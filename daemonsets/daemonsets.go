package daemonsets

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

func parseDaemonset(object []byte) (*v1.DaemonSet, error) {
	var ds v1.DaemonSet
	if err := json.Unmarshal(object, &ds); err != nil {
		return nil, err
	}

	return &ds, nil
}