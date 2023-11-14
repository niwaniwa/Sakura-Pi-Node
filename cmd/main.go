package main

import (
	config "Sakura-Hardware/pkg"
	"log"
)

func main() {
	config := config.LoadEnvironments()
	log.Print(config)
}
