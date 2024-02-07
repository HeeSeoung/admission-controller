package server

import (
	"fmt"
	"net/http"

	"github.com/HeeSeoung/admission-controller/daemonsets"
	"github.com/HeeSeoung/admission-controller/deployments"
	"github.com/HeeSeoung/admission-controller/nodes"
)

// NewServer creates and return a http.Server
func NewServer(port string) *http.Server {
	// Instances hooks
	deploymentMutation := deployments.NewMutationHook()
	nodeMutation := nodes.NewMutationHook()
	daemonsetMutation := daemonsets.NewMutationHook()

	// Routers
	ah := newAdmissionHandler()
	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz())
	mux.Handle("/mutate/deployments", ah.Serve(deploymentMutation))
	mux.Handle("/mutate/daemonsets", ah.Serve(daemonsetMutation))
	mux.Handle("/mutate/nodes", ah.Serve(nodeMutation))
	

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
