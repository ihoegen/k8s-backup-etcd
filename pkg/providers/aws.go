package providers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//AWSProvider is an implementation of the Provider class for AWS
type AWSProvider struct {
	RotationDepth int
	EC2Client     *ec2.EC2
}

// NewAWSProvider returns an AWS Provider Object
func NewAWSProvider(rotationDepth int) Provider {
	return AWSProvider{
		RotationDepth: rotationDepth,
		EC2Client:     ec2.New(session.New()),
	}
}

// GetVolumes returns a list of etcd volumes to backup
func (a AWSProvider) GetVolumes() ([]string, error) {
	etcdVolumes := []string{}
	volumes, err := a.EC2Client.DescribeVolumes(&ec2.DescribeVolumesInput{
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
		return nil, err
	}
	for _, v := range volumes.Volumes {
		etcdVolumes = append(etcdVolumes, *v.VolumeId)
	}
	return etcdVolumes, nil
}

// CreateSnapshots creates a snapshot from a list of volumes
func (a AWSProvider) CreateSnapshots(volumes []string) ([]string, error) {
	errors := ""
	snapshots := []string{}
	for _, v := range volumes {
		tags, err := a.EC2Client.DescribeTags(&ec2.DescribeTagsInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name: aws.String("resource-id"),
					Values: []*string{
						&v,
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
		var tagList []*ec2.Tag
		for _, t := range tags.Tags {
			tagList = append(tagList, &ec2.Tag{
				Key:   t.Key,
				Value: t.Value,
			})
		}
		snapshot, err := a.EC2Client.CreateSnapshot(&ec2.CreateSnapshotInput{
			VolumeId: &v,
			TagSpecifications: []*ec2.TagSpecification{
				&ec2.TagSpecification{
					Tags:         tagList,
					ResourceType: aws.String("snapshot"),
				},
			},
		})
		if err != nil {
			errors += err.Error() + "\n"
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
func (a AWSProvider) RotateSnapshots() error {
	snapshots, err := a.EC2Client.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
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
	if len(snapshots.Snapshots) <= a.RotationDepth {
		return nil
	}
	errorString := ""
	for i := len(snapshots.Snapshots) - 1; i >= a.RotationDepth; i-- {
		_, err := a.EC2Client.DeleteSnapshot(&ec2.DeleteSnapshotInput{
			SnapshotId: snapshots.Snapshots[i].SnapshotId,
		})
		if err != nil {
			errorString += err.Error() + "\n"
		}
	}
	if errorString != "" {
		return fmt.Errorf(errorString)
	}
	return nil
}
