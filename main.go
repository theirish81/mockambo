package main

import (
	"log"
	"mockambo/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
