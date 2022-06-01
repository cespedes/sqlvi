package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lib/pq"
)

func main() {
	err := run(os.Args[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

type change struct {
	query string
	args  []string
}

type app struct {
	Debug      bool
	ConfigFile string
	Editor     string
	Format     string
	Connect    string
	Select     string
	Insert     string
	Update     string
	Delete     string

	InsertChanges []change
	UpdateChanges []change
	DeleteChanges []change
}

func run(args []string) error {
	app := app{}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.BoolVar(&app.Debug, "debug", false, "Debugging")
	flags.StringVar(&app.ConfigFile, "config", filepath.Join(os.Getenv("HOME"), ".sqlvirc"), "Config file")
	flags.StringVar(&app.Connect, "connect", "", "Connect string to database")
	flags.StringVar(&app.Format, "format", "", "Output format to use (default \"org\")")
	flags.StringVar(&app.Editor, "editor", "", "Editor to use")
	flags.StringVar(&app.Select, "select", "", "Database query to display results")
	flags.StringVar(&app.Insert, "insert", "", "Database query to add a new row")
	flags.StringVar(&app.Update, "update", "", "Database query to update a row")
	flags.StringVar(&app.Delete, "delete", "", "Database query to delete a row")
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: sqlvi [options] [<mode from config file>]")
		fmt.Fprintln(os.Stderr, "Options:")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if len(flags.Args()) > 1 {
		return fmt.Errorf("too many arguments")
	}
	modeName := ""
	if len(flags.Args()) == 1 {
		modeName = flags.Args()[0]
	}

	err := app.readConfig(modeName)
	if err != nil {
		return err
	}
	if app.Connect == "" {
		return fmt.Errorf("no connection to database specified")
	}
	if app.Select == "" {
		return fmt.Errorf("no query specified")
	}

	db, err := sqlConnect(app.Connect)
	if err != nil {
		return err
	}
	orig, err := sqlGenericQuery(db, app.Select)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "sqlvi.*.org")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	defer os.Remove(tmpName)

	switch app.Format {
	case "org":
		writeOrgTable(f, orig.Columns, orig.Strings)
	case "ini":
		writeINI(f, orig.Columns, orig.Strings)
	case "ldap":
		fallthrough
	case "yaml":
		writeYAML(f, orig.Columns, orig.Strings)
	default:
		return fmt.Errorf("unknown output format: %q", app.Format)
	}
	if err = f.Close(); err != nil {
		return err
	}

editorLoop:
	for {
		var tx *sql.Tx
		var c rune

		err = app.callEditor(f.Name())
		if err != nil {
			return err
		}
		f, err := os.Open(tmpName)
		if err != nil {
			return err
		}
		var data [][]string
		switch app.Format {
		case "org":
			data, err = readOrgTable(f, orig.Columns)
		case "ini":
			data, err = readINI(f, orig.Columns)
		case "ldap":
			fallthrough
		case "yaml":
			data, err = readYAML(f, orig.Columns)
		default:
			panic("unknown output format")
		}
		f.Close()
		if err != nil {
			goto errorFound
		}
		if app.Debug {
			log.Printf("data = %v\n", data)
		}

		// Wrong lines (primary key different from all the original ones):
	findWrongLines:
		for _, row := range data {
			for _, o := range orig.Strings {
				if row[0] == "" || row[0] == o[0] {
					continue findWrongLines
				}
			}
			err = fmt.Errorf("unexpected entry with id %q in table", row[0])
			goto errorFound
		}

		// Lines with duplicated primary keys:
		for i, row := range data {
			for j, row2 := range data {
				if i != j && row[0] != "" && row[0] == row2[0] {
					err = fmt.Errorf("duplicate entry %s", row[0])
					goto errorFound
				}
			}
		}

		app.InsertChanges = nil
		app.UpdateChanges = nil
		app.DeleteChanges = nil

		// New rows:
		for _, row := range data {
			if row[0] == "" {
				query, args := sqlBind(db, app.Insert, row)
				if app.Debug {
					log.Printf("New row: %s %v", query, args)
				}
				app.InsertChanges = append(app.InsertChanges, change{query, args})
			}
		}

		// Removed rows (primary key exists in original but not in new):
	findRemovedLines:
		for _, o := range orig.Strings {
			for _, row := range data {
				if o[0] == row[0] {
					continue findRemovedLines
				}
			}
			query, args := sqlBind(db, app.Delete, o)
			if app.Debug {
				log.Printf("Updated row: %s %v", query, args)
			}
			app.DeleteChanges = append(app.DeleteChanges, change{query, args})
		}

		// Modified rows (different fields in orig and new):
		for _, o := range orig.Strings {
			for _, row := range data {
				if o[0] == row[0] {
					for i := range row {
						if o[i] != row[i] {
							query, args := sqlBind(db, app.Update, row)
							if app.Debug {
								log.Printf("Updated row: %s %v", query, args)
							}
							app.UpdateChanges = append(app.UpdateChanges, change{query, args})
							break
						}
					}
					break
				}
			}
		}
		if len(app.InsertChanges) == 0 && len(app.UpdateChanges) == 0 && len(app.DeleteChanges) == 0 {
			fmt.Println("No changes.")
			return nil
		}
		if len(app.InsertChanges) == 0 {
			fmt.Print("add: 0")
		} else {
			fmt.Printf("\033[32;1madd: %d\033[m", len(app.InsertChanges))
		}
		fmt.Printf(", ")
		if len(app.UpdateChanges) == 0 {
			fmt.Print("modify: 0")
		} else {
			fmt.Printf("\033[33;1mmodify: %d\033[m", len(app.UpdateChanges))
		}
		fmt.Printf(", ")
		if len(app.DeleteChanges) == 0 {
			fmt.Print("delete: 0")
		} else {
			fmt.Printf("\033[31;1mdelete: %d\033[m", len(app.DeleteChanges))
		}
		fmt.Println()
		c = ask("Action?", []askStruct{
			{'y', "commit changes"},
			{'e', "open editor again"},
			{'Q', "discard changes and quit"},
		})
		switch c {
		case 'e':
			continue editorLoop
		case 'Q':
			return nil
		}
		if app.Debug {
			log.Println("Commiting changes...")
		}
		tx, err = db.Begin()
		if err != nil {
			return err
		}
		for _, c := range append(app.InsertChanges, append(app.UpdateChanges, app.DeleteChanges...)...) {
			if app.Debug {
				log.Printf(">> %v\n", c)
			}
			args := make([]interface{}, len(c.args))
			for i, v := range c.args {
				args[i] = v
			}
			_, err = tx.Exec(c.query, args...)
			if err != nil {
				goto errorFound
			}
		}
		err = tx.Commit()
		if err != nil {
			goto errorFound
		}
		return nil
	errorFound:
		fmt.Printf("Error: %s\n", err.Error())
		if err, ok := err.(*pq.Error); ok && err.Detail != "" {
			fmt.Printf("Details: %s\n", err.Detail)
		}

		c = ask(`What now?`, []askStruct{
			{'e', "open editor again"},
			{'Q', "discard changes and quit"},
		})
		if c == 'e' {
			continue editorLoop
		}
		return nil
	}
}
