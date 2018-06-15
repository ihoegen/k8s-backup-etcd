package main

import (
	"log"

	"github.com/integrii/flaggy"
)

var cloudProvider = ""

func main() {
	flaggy.AddPositionalValue(&cloudProvider, "cloudProvider", 1, true, "The cloud provider that is being used")
	flaggy.Parse()
	switch cloudProvider {
	case "aws":
		backupAWS()
		break
	default:
		log.Panic("Invalid cloud provider provided")
	}
}
