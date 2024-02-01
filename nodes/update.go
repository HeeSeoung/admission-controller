package nodes

import (
	"github.com/HeeSeoung/admission-controller"

	admission "k8s.io/api/admission/v1"

	log "k8s.io/klog/v2"
)

func mutateUpdate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		var operations []admissioncontroller.PatchOperation
		node, err := parseNode(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}

		log.Infof("%+v", node)

		newNodeLabel := node.Labels
		if newNodeLabel == nil {
			newNodeLabel = make(map[string]string)
		}
		newNodeLabel["node-mutating"] = "true"
		log.Infof("%+v", newNodeLabel)
		operations = append(operations, admissioncontroller.ReplacePatchOperation("metadata/labels", newNodeLabel))

		// Add a simple annotation using `AddPatchOperation`
		metadata := map[string]string{"origin": "fromMutation"}
		operations = append(operations, admissioncontroller.AddPatchOperation("/metadata/annotations", metadata))
		return &admissioncontroller.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}