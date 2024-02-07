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
			var containers []v1.Container
			containers = append(containers, ds.Spec.Template.Spec.InitContainers...)
			initC := v1.Container{
				Name:    "test-sidecar",
				Image:   "busybox:stable",
				Command: []string{"sh", "-c", "while true; do echo 'I am a container injected by mutating webhook'; sleep 2; done"},
			}
			containers = append(containers, initC)
			log.Infof("%+v", containers)
			operations = append(operations, admissioncontroller.ReplacePatchOperation("/spec/template/spec/containers", containers))
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
