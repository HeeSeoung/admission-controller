package deployments

import (
	"github.com/HeeSeoung/admission-controller"

	admission "k8s.io/api/admission/v1"

	v1 "k8s.io/api/core/v1"
	log "k8s.io/klog/v2"
)

func validateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		dp, err := parseDeployment(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}

		if dp.Namespace == "special" {
			return &admissioncontroller.Result{Msg: "You cannot create a deployment in `special` namespace."}, nil
		}

		return &admissioncontroller.Result{Allowed: true}, nil
	}
}

func mutateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		var operations []admissioncontroller.PatchOperation
		dp, err := parseDeployment(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}
		log.Infof("%+v", dp)
		// Very simple logic to inject a new "sidecar" container.
		if dp.Namespace == "production" {
			var containers []v1.Container
			containers = append(containers, dp.Spec.Template.Spec.Containers...)
			sideC := v1.Container{
				Name:    "test-sidecar",
				Image:   "busybox:stable",
				Command: []string{"sh", "-c", "while true; do echo 'I am a container injected by mutating webhook'; sleep 2; done"},
			}
			containers = append(containers, sideC)
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