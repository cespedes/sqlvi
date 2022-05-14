package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	result, err := db.sqlGenericQuery(mode.Select)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("result = %v\n", result)
	//return

	/*
		columns := []string{"one", "two", "three"}
		data := [][]string{
			{"1", "2", "3"},
			{"4", "5", "6"},
		}
	*/
	f, err := os.CreateTemp("", "sqlvi.*.org")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())

	writeOrgtable(f, result.Columns, result.Strings)

	err = callEditor(f.Name())
	if err != nil {
		log.Fatal(err)
	}
}
