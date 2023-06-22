package logging

import (
   "encoding/json"
   log "github.com/sirupsen/logrus"
)

var DefaultLogLevel = "info"

func SetLogLevel(lvl string) {
   log.Printf("logging.SetLogLevel: enter: LogLevel: %s", lvl)
   // parse string, this is built-in feature of logrus
   ll, err := log.ParseLevel(lvl)
   if err != nil {
      ll = log.DebugLevel
   }
   // set global log level
   log.SetLevel(ll)
   log.SetFormatter(&log.TextFormatter{ForceQuote: false, DisableQuote: true})
}

func Dump(level log.Level, m interface{}, s string) {
   l := log.StandardLogger()
   l.SetLevel(log.GetLevel())
   b, err := json.MarshalIndent(m, "", "   ")
   if err == nil {
      l.Logln(level, s, " ", string(b))
   } else {
      l.Logf(level, "Marshal error: %s %v", s, err)
   }
}
