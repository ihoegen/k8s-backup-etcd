package main

import (
	"log"

	"github.com/ihoegen/backup-etcd/pkg/providers"
	"github.com/integrii/flaggy"
)

var cloudProvider = ""
var snapshotsSaved = 500

func main() {
	flaggy.AddPositionalValue(&cloudProvider, "cloudProvider", 1, true, "The cloud provider that is being used")
	flaggy.Int(&snapshotsSaved, "", "snapshots-saved", "The number of snapshots to keep on rotation")
	flaggy.Parse()
	var p providers.Provider
	switch cloudProvider {
	case "aws":
		p = providers.NewAWSProvider(snapshotsSaved)
		break
	default:
		log.Panic("Invalid cloud provider provided")
	}
	err := p.RotateSnapshots()
	if err != nil {
		log.Println(err)
	}
	volumes, err := p.GetVolumes()
	if err != nil {
		log.Fatal(err)
	}
	snapshots, err := p.CreateSnapshots(volumes)
	if err != nil {
		log.Println(err)
	}
	if snapshots != nil {
		log.Printf("Following snapshots taken: %v", snapshots)
	}

}
