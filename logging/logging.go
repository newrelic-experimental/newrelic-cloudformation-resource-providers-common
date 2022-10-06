package logging

import (
   "encoding/json"
   log "github.com/sirupsen/logrus"
)

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
