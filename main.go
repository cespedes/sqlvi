package main

import (
	"log"
	"os"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := readConfig()
	if err != nil {
		return err
	}

	mode := config.Modes[config.Default]
	if len(os.Args) == 2 {
		var ok bool
		mode, ok = config.Modes[os.Args[1]]
		if !ok {
			log.Fatalf("mode %q not found in config file", os.Args[1])
		}
	}
	log.Printf("Using mode = %q\n", mode)

	db, err := sqlConnect(mode.Connect)
	if err != nil {
		return err
	}
	result, err := db.sqlGenericQuery(mode.Select)
	if err != nil {
		return err
	}
	log.Printf("result = %v\n", result)

	f, err := os.CreateTemp("", "sqlvi.*.org")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	defer os.Remove(tmpName)

	writeOrgTable(f, result.Columns, result.Strings)
	if err = f.Close(); err != nil {
		return err
	}

	err = callEditor(f.Name())
	if err != nil {
		return err
	}

	f, err = os.Open(tmpName)
	if err != nil {
		return err
	}
	cols, data, err := readOrgTable(f)
	if err != nil {
		return err
	}
	log.Printf("cols = %v\n", cols)
	log.Printf("data = %v\n", data)
	return nil
}
