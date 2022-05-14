package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
       drop: DROP countries WHERE id=$1
     products:
       connect: postgres://user:pass@host/database
       select: SELECT id,name,price FROM products
       insert: INSERT INTO products (name,price) VALUES ($2,$3)
       update: UPDATE products SET name=$2,price=$3 WHERE id=$1
       drop: DROP product WHERE id=$1
*/

type configMode struct {
	Connect string
	Select  string
	Insert  string
	Update  string
	Drop    string
}

type config struct {
	Editor  string
	Default string
	Modes   map[string]configMode
}

func readConfig() (config, error) {
	configFileName := filepath.Join(os.Getenv("HOME"), ".sqlvi.yaml")

	log.Printf("Reading config file %s", configFileName)
	c := config{}
	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return c, fmt.Errorf("reading %s: %w", configFileName, err)
	}
	if err := yaml.UnmarshalStrict(data, &c); err != nil {
		return c, fmt.Errorf("parsing %s: %w", configFileName, err)
	}
	return c, nil

}
