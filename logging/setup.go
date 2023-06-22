//go:build logging

package logging

import (
   log "github.com/sirupsen/logrus"
   goLog "log"
   "sync"
)

var doOnce sync.Once

// Setup cannot be init as it is dependent on the generated main running first
func Setup() {
   goLog.Printf("logging.Setup: build tag logging")
   // doOnce.Do(func() {
   log.SetOutput(goLog.Writer())
   //    // fmt.Println("")
   //    // fmt.Println(os.Environ())
   //    // fmt.Println("")
   //    // if !(strings.ToLower(os.Getenv("AWS_SAM_LOCAL")) == "true") {
   //    //    fmt.Println("Inside if statement")
   //    //    fmt.Println("")
   //    //    debug.SetPanicOnFault(true)
   //    //    w := ol.Writer()
   //    //    _ = w
   //    //    // log.Debugf("Log writer: %T %#v", w, w)
   //    //    log.SetOutput(w)
   //    //    debug.SetPanicOnFault(false)
   //    // } else {
   //    //    log.SetOutput(os.Stderr)
   //    //    fmt.Println("Inside else statement")
   //    //    fmt.Println("")
   //    // }
   // })
}
