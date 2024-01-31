package deployments

import (
	"encoding/json"

	"github.com/HeeSeoung/admission-controller"

	v1 "k8s.io/api/apps/v1"
)

// NewMutationHook creates a new instance of pods mutation hook
func NewMutationHook() admissioncontroller.Hook {
	return admissioncontroller.Hook{
		Create: mutateCreate(),
	}
}

func parseDeployment(object []byte) (*v1.Deployment, error) {
	var dp v1.Deployment
	if err := json.Unmarshal(object, &dp); err != nil {
		return nil, err
	}

	return &dp, nil
}