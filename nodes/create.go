package nodes

import (
	"github.com/HeeSeoung/admission-controller"

	admission "k8s.io/api/admission/v1"

	// log "k8s.io/klog/v2"
)

func mutateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		var operations []admissioncontroller.PatchOperation
		node, err := parseNode(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}

		// Add a new node label
		newNodeLabel := node.Labels
		if newNodeLabel == nil {
			newNodeLabel = make(map[string]string)
		}
		newNodeLabel["test-label"] = "true"
		operations = append(operations, admissioncontroller.ReplacePatchOperation("/metadata/labels", newNodeLabel))

		// Add a simple annotation using `AddPatchOperation`
		metadata := map[string]string{"origin": "fromMutation"}
		operations = append(operations, admissioncontroller.AddPatchOperation("/metadata/annotations", metadata))
		return &admissioncontroller.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}
