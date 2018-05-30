package main

import (
	"log"

	"github.com/iStreamplanet/backup-etcd/pkg/cloudproviders/aws"
)

func backupAWS() {
	snapshots, err := aws.CreateBackup()
	if snapshots != nil {
		log.Println("Snapshots created:", snapshots)
	}
	if err != nil {
		log.Println("The following errors occurred:", err)
	}
}
