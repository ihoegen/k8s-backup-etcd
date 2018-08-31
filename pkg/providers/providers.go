package providers

type Provider interface {
	GetVolumes() ([]string, error)
	CreateSnapshots([]string) ([]string, error)
	RotateSnapshots() error
}
