package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//CreateBackup finds ETCD volumes, then creates snapshots
func CreateBackup() ([]string, error) {
	ec2Client := ec2.New(session.New())
	err := RotateSnapshots(ec2Client, 4000)
	if err != nil {
		fmt.Println(err)
	}
	volumes, err := FindETCDVolumes(ec2Client)
	if err != nil {
		return nil, err
	}
	return CreateSnapshots(ec2Client, volumes)
}

// FindETCDVolumes returns a list of etcd volumes to backup
func FindETCDVolumes(client *ec2.EC2) ([]ec2.Volume, error) {
	etcdVolumes := []ec2.Volume{}
	volumes, err := client.DescribeVolumes(&ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, err
	}
	for _, v := range volumes.Volumes {
		for _, t := range v.Tags {
			if *t.Key == "k8s.io/etcd/events" || *t.Key == "k8s.io/etcd/main" {
				etcdVolumes = append(etcdVolumes, *v)
			}
		}
	}
	return etcdVolumes, nil
}

// CreateSnapshots creates a snapshot for a list of volumes
func CreateSnapshots(client *ec2.EC2, volumes []ec2.Volume) ([]string, error) {
	errors := ""
	snapshots := []string{}
	for _, v := range volumes {
		snapshot, err := client.CreateSnapshot(&ec2.CreateSnapshotInput{
			VolumeId: v.VolumeId,
			TagSpecifications: []*ec2.TagSpecification{
				&ec2.TagSpecification{
					Tags:         v.Tags,
					ResourceType: aws.String("snapshot"),
				},
			},
		})
		if err != nil {
			errors += "\n" + err.Error()
		} else {
			snapshots = append(snapshots, *snapshot.SnapshotId)
		}
	}
	if errors == "" {
		return snapshots, nil
	}
	return snapshots, fmt.Errorf(errors)
}

// RotateSnapshots ensures that no more than depth amount of snapshots exist
func RotateSnapshots(client *ec2.EC2, depth int) error {
	snapshots, err := client.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag-key"),
				Values: []*string{
					aws.String("k8s.io/etcd/events"),
					aws.String("k8s.io/etcd/main"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	if len(snapshots.Snapshots) <= depth {
		return nil
	}
	for i := len(snapshots.Snapshots) - 1; i >= depth; i-- {
		_, err := client.DeleteSnapshot(&ec2.DeleteSnapshotInput{
			SnapshotId: snapshots.Snapshots[i].SnapshotId,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
