package main

import (
	"log"
	"os"

	"github.com/Liar233/throttles-tank/internal/client"
)

func main() {

	app := client.NewClient()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
