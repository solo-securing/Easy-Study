package main

import (
	"fmt"
	"log"
	"packages/configloader"
	"packages/configloader/parsers/json"
	"packages/configloader/parsers/yaml"
	"packages/configloader/providers/file"
)

var cfLoader = configloader.New()

func main() {
	// Load JSON config.
	fileProvider := file.Provider("../../mocks/mock.json")
	if err := cfLoader.Load(fileProvider, json.Parser()); err != nil {
		log.Fatal(err)
	}

	// Load YAML config and merge into the previously loaded config (because we can).
	cfLoader.Load(file.Provider("../../mocks/mock.yml"), yaml.Parser())

	fmt.Println("parent's name is = ", cfLoader.String("parent1.name"))
	fmt.Println("parent's ID is = ", cfLoader.Int("parent1.id"))
}
