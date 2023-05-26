package tags

import (
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/client/nerdgraph"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

// Resulthandler at a minimum provides access to the default error processing.
// If required we can provide custom processing here via composition overrides
type Resulthandler struct {
   // Use Go composition to access the default implementation
   model.ResultHandler
}

// NewResultHandler This is all pretty magical. We return the interface so common is insulated from an implementation. Payload implements model.Model so all is good
func NewResultHandler() (h model.ResultHandler) {
   defer func() {
      log.Debugf("(tagging) errorHandler.NewErrorHandler: exit %p", h)
   }()
   // Initialize ourself with the common core
   h = &Resulthandler{ResultHandler: nerdgraph.NewResultHandler()}
   return
}
