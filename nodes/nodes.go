package nodes

import (
	"encoding/json"

	"github.com/HeeSeoung/admission-controller"

	"k8s.io/api/core/v1"
)

// NewValidationHook updates a new instance of deployment validation hook
func NewValidationHook() admissioncontroller.Hook {
	return admissioncontroller.Hook{
	}
}

// NewMutationHook updates a new instance of pods mutation hook
func NewMutationHook() admissioncontroller.Hook {
	return admissioncontroller.Hook{
		Update: mutateUpdate(),
	}
}

func parseNode(object []byte) (*v1.Node, error) {
	var node v1.Node
	if err := json.Unmarshal(object, &node); err != nil {
		return nil, err
	}

	return &node, nil
}
