package main

import (
	"log"

	"github.com/lightclient/bazooka/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if err.Error() == "^C" {
			return
		}
		log.Fatalln(err)
	}
}
