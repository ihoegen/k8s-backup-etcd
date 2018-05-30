package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/iStreamplanet/backup-etcd/pkg/cloudproviders/aws"
)

func main() {
	lambda.Start(handleLambda)
}

func handleLambda(ctx context.Context) (string, error) {
	snapshots, err := aws.CreateBackup()
	return fmt.Sprintf("Snapshots created: %v", snapshots), err
}
