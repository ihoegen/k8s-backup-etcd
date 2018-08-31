package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ihoegen/backup-etcd/pkg/providers"
)

func main() {
	lambda.Start(handleLambda)
}

func handleLambda(ctx context.Context) (string, error) {
	var p providers.Provider
	err := p.RotateSnapshots()
	if err != nil {
		log.Info(err)
	}
	volumes, err := p.GetVolumes()
	if err != nil {
		return "", err
	}
	snapshots, err := p.CreateSnapshots(volumes)
	if err != nil {
		return "", err
	}
	if snapshots != nil {
		return fmt.Sprintf("Snapshots taken: %v", snapshots), nil
	}
}
