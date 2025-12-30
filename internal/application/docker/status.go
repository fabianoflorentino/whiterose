package docker

type DockerStatus string

const (
	StatusUnKnown  DockerStatus = "unknown"
	StatusMissing  DockerStatus = "missing"
	StatusBuilding DockerStatus = "building"
	StatusReady    DockerStatus = "ready"
	StatusFailed   DockerStatus = "failed"
)

type ImageStatus struct {
	Status DockerStatus
	Error  string
}

func NewMissingStatus() ImageStatus {
	return ImageStatus{Status: StatusMissing}
}

func NewBuildingStatus() ImageStatus {
	return ImageStatus{Status: StatusBuilding}
}

func NewReadyStatus() ImageStatus {
	return ImageStatus{Status: StatusReady}
}

func NewFailedStatus(err error) ImageStatus {
	return ImageStatus{
		Status: StatusFailed,
		Error:  err.Error(),
	}
}
