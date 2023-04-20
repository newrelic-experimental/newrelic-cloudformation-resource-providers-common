package tags

import (
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/client/nerdgraph"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

// ErrorHandler at a minimum provides access to the default error processing.
// If required we can provide custom processing here via composition overrides
type ErrorHandler struct {
   // Use Go composition to access the default implementation
   model.ErrorHandler
   M model.Model
}

// NewErrorHandler This is all pretty magical. We return the interface so common is insulated from an implementation. Payload implements model.Model so all is good
func NewErrorHandler(p *Payload) (h model.ErrorHandler) {
   defer func() {
      log.Debugf("(tagging) errorHandler.NewErrorHandler: exit %p", h)
   }()
   // Initialize ourself with the common core
   h = &ErrorHandler{ErrorHandler: nerdgraph.NewCommonErrorHandler(p), M: p}
   return
}
