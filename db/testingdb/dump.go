package testingdb

import (
	"encoding/json"
	"fmt"

	"github.com/script-development/RT-CV/db"
)

// Dump prints the full database contents in the console
// This can be used in tests to dump the contents of the database might something fail or to debug
//
// shouldPanic controls if the output is only printed or also should panic
func (c *TestConnection) Dump(shouldPanic bool) {
	data := map[string][]db.Entry{}
	for _, collection := range c.collections {
		data[collection.name] = collection.data
	}

	jsonBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)
	if shouldPanic {
		panic(jsonString)
	} else {
		fmt.Println(jsonString)
	}
}

// DumpCollection prints a full database collection it's contents in the console
// This can be used in tests to dump the contents of the database might something fail or to debug
//
// shouldPanic controls if the output is only printed or also should panic
func (c *TestConnection) DumpCollection(entry db.Entry, shouldPanic bool) {
	collection := c.getCollectionFromEntry(entry)

	jsonBytes, err := json.MarshalIndent(collection.data, "", "    ")
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)
	if shouldPanic {
		panic(jsonString)
	} else {
		fmt.Println(jsonString)
	}
}
