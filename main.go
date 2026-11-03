package main

import (
	"log"
	"os"

	lib "obs-lan-controller/Lib"
)

func main() {
	if err := lib.Run(os.Args[1:]); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
