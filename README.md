# Backup ETCD
This project is designed to backup ETCD instances, and can be run on a Kubernetes cluster, via the command line, or as a lambda.

## Usage

`./backup-etcd [cloudProvider]`

## Deploy

### AWS
To deploy to AWS, a Makefile command exists:

`make deploy-aws`

### Kubernetes


In the `deploy/kubernetes` directory, there is a helm chart named backup-etcd.

To use this, simply run `helm install deploy/kubernetes/backup-etcd`

This assumes that you have Kube2IAM installed, and that `Values.iam` is mapped to a role with the necessary permissions. 