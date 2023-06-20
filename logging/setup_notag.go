//go:build !logging

package logging

import (
   log "github.com/sirupsen/logrus"
)

// Setup cannot be init as it is dependent on the generated main running first
func Setup() {
   log.Debugf("logging.Setup: no build tags")
}
