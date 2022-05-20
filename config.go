package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

/*
   editor: vim  # optional
   default: countries
   modes:
     countries:
       connect: postgres://user:pass@host/database
       select: SELECT id,country,capital,population FROM countries
       insert: INSERT INTO countries (country,capital,population) VALUES ($2,$3,$4)
       update: UPDATE countries SET country=$2,capital=$3,population=$4 WHERE id=$1
       delete: DELETE FROM countries WHERE id=$1
     products:
       connect: postgres://user:pass@host/database
       select: SELECT id,name,price FROM products
       insert: INSERT INTO products (name,price) VALUES ($2,$3)
       update: UPDATE products SET name=$2,price=$3 WHERE id=$1
       delete: DELETE FROM product WHERE id=$1
*/

type configMode struct {
	Connect string
	Select  string
	Insert  string
	Update  string
	Delete  string
}

type config struct {
	Editor  string
	Default string
	Modes   map[string]configMode
}

func (app *app) readConfig(modeName string) error {
	var config config

	if app.Debug {
		log.Printf("Reading config file %s", app.ConfigFile)
	}
	data, err := ioutil.ReadFile(app.ConfigFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", app.ConfigFile, err)
	}
	if err := yaml.UnmarshalStrict(data, &config); err != nil {
		return fmt.Errorf("parsing %s: %w", app.ConfigFile, err)
	}

	if modeName == "" {
		modeName = config.Default
	}

	mode, ok := config.Modes[modeName]
	if !ok {
		return fmt.Errorf("mode %q not found in config file", modeName)
	}

	if app.Debug {
		log.Printf("Using mode = %q\n", modeName)
	}

	if app.Connect == "" {
		app.Connect = mode.Connect
	}
	if app.Select == "" {
		app.Select = mode.Select
	}
	if app.Insert == "" {
		app.Insert = mode.Insert
	}
	if app.Update == "" {
		app.Update = mode.Update
	}
	if app.Delete == "" {
		app.Delete = mode.Delete
	}

	return nil
}
