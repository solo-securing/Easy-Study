package main

import (
	"fmt"
	"log"
	"os"
	"packages/configloader"
	"packages/configloader/parsers/json"
	"packages/configloader/providers/env"
	"packages/configloader/providers/file"
)

var cfLoader = configloader.New()

func main() {
	// Load JSON config.
	if err := cfLoader.Load(file.Provider("../../mocks/mock.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	key := "MYVAR_KEY_1"
	val := "value_1"
	os.Setenv(key, val)
	defer os.Unsetenv(key)

	cfLoader.Load(env.Provider(env.Opt{
		Prefix: "MYVAR_",
	}), nil)

	fmt.Println(cfLoader.String("parent1.child1.name"))
	fmt.Print(cfLoader.String("MYVAR_KEY_1"))
}
