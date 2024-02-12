package daemonsets

import (
	"github.com/HeeSeoung/admission-controller"

	admission "k8s.io/api/admission/v1"

	v1 "k8s.io/api/core/v1"
	log "k8s.io/klog/v2"
)

func mutateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		var operations []admissioncontroller.PatchOperation
		ds, err := parseDaemonset(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}
		log.Infof("%+v", ds)
		// Very simple logic to inject a new "sidecar" container.
		if ds.Namespace == "production" {
			ds.Spec.Template.Spec.InitContainers = append(ds.Spec.Template.Spec.InitContainers, v1.Container{
				Name:  "test",
				Image: "busybox",
				Command: []string{
					"/bin/sh",
					"-c",
					"sleep 20",
				},
			})
			operations = append(operations, admissioncontroller.ReplacePatchOperation("/spec/template/spec/", ds.Spec.Template.Spec.InitContainers))
		}

		// Add a simple annotation using `AddPatchOperation`
		metadata := map[string]string{"origin": "fromMutation"}
		operations = append(operations, admissioncontroller.AddPatchOperation("/metadata/annotations", metadata))
		return &admissioncontroller.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}
