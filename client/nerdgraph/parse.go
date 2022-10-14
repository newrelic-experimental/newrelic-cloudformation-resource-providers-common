package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   log "github.com/sirupsen/logrus"
)

func findAllKeyValues(data []byte, key string) (values []interface{}, err error) {
   defer func() {
      log.Debugf("findKeyValues: returning: %v %T err: %v", values, values, err)
   }()

   m := map[string]interface{}{}
   err = json.Unmarshal(data, &m)
   if err != nil {
      log.Errorf("findKeyValue: unmarshal %v", err)
      err = fmt.Errorf("%w %v", &cferror.UnknownError{}, err)
      return
   }
   values = make([]interface{}, 0)
   _findAllKeyValues(m, key, &values)
   if len(values) <= 0 {
      err = fmt.Errorf("%w key not found: %s", &cferror.NotFound{}, key)
   }

   return
}

func findKeyValue(data []byte, key string) (value interface{}, err error) {
   defer func() {
      log.Debugf("findKeyValue: returning: %v %T err: %v", value, value, err)
   }()

   v, err := findAllKeyValues(data, key)
   if err != nil {
      return
   }
   if len(v) <= 0 {
      err = fmt.Errorf("%w key not found: %s", &cferror.NotFound{}, key)
      return
   }
   value = v[0]
   return
}

func _findAllKeyValues(m map[string]interface{}, key string, values *[]interface{}) {
   for k, v := range m {
      log.Tracef("_findAllKeyValues: k(ey): %s v(alue): %v type: %T", k, v, v)
      if k == key {
         *values = append(*values, v)
         return
      }
      switch v.(type) {
      case map[string]interface{}:
         _findAllKeyValues(v.(map[string]interface{}), key, values)
      case []interface{}:
         for _, e := range v.([]interface{}) {
            if m, ok := e.(map[string]interface{}); ok {
               _findAllKeyValues(m, key, values)
            } else {
               log.Warnf("_findAllKeyValues: skipping [] of unknown type: %T at key: %s", e, k)
            }
         }
      default:
      }
   }
}
