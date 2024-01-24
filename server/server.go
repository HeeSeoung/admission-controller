package server

import (
	"fmt"
	"net/http"

	"github.com/douglasmakey/admissioncontroller/deployments"
)

// NewServer creates and return a http.Server
func NewServer(port string) *http.Server {
	// Instances hooks
	deploymentMutation := deployments.NewMutationHook()

	// Routers
	ah := newAdmissionHandler()
	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz())
	mux.Handle("/mutate/deployments", ah.Serve(deploymentMutation))

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}