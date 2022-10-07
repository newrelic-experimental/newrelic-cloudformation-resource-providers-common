package nerdgraph

import (
   log "github.com/sirupsen/logrus"
   "strings"
)

// TODO
func hasErrors(data *[]byte) (err error) {
   // Empty
   if data == nil {
      return
   }

   // No error keyword
   s := string(*data)
   if !strings.Contains(strings.ToLower(s), "error") {
      return
   }

   // r := root{}
   // err = json.Unmarshal(*data, &r)
   // if err != nil {
   //    log.Errorf("hasError: unmarshal %v", err)
   //    return fmt.Errorf("%w %s", &cferror.UnknownError{}, s)
   // }
   //
   // // Nothing in the error array
   // if r.Errors == nil || len(r.Errors) <= 0 {
   //    return
   // }
   //
   // // At this point we actually have something
   // if len(r.Errors) > 1 {
   //    log.Warnf("%d errors returned from NerdGraph, the first is used the remainder logged", len(r.Errors))
   // }
   // for i, e := range r.Errors {
   //    // Don't log the first error, we'll return it as the error value
   //    if i == 0 {
   //       continue
   //    }
   //    utils.DumpModel(log.ErrorLevel, e, "NerdGraph error")
   // }
   // var errorCode = 0
   // var errorMessage = ""
   // _splunkError(r.Errors[0], &errorCode, &errorMessage)
   // log.Infof("code: %d message: %s", errorCode, errorMessage)
   //
   // // In-case we can't find a specific error
   // if errorCode == 0 {
   //    log.Errorf("hasError: non-specific error %s", s)
   //    return fmt.Errorf("%w %s", &cferror.UnknownError{}, s)
   // }
   //
   // switch errorCode {
   // case 404:
   //    err = fmt.Errorf("%w Not found", &cferror.NotFound{})
   // default:
   //    err = fmt.Errorf("%w code: %d message: %s", &cferror.UnknownError{}, errorCode, errorMessage)
   // }
   return
}

func _splunkError(m map[string]interface{}, i *int, s *string) {
   for k, v := range m {
      k = strings.ToLower(k)
      log.Debugf("_splunkError: key: %s value: %v type: %T", k, v, v)
      switch v.(type) {
      case float64:
         if strings.Contains(k, "code") {
            log.Debugln("_splunkError: setting error code")
            *i = int(v.(float64))
         }
      case int:
         if strings.Contains(k, "code") {
            log.Debugln("_splunkError: setting error code")
            *i = v.(int)
         }
      case string:
         if strings.Contains(k, "message") {
            log.Debugln("_splunkError: setting error message")
            *s = v.(string)
         }
      case map[string]interface{}:
         log.Debugln("_splunkError: recursion on map[string]interface{}")
         _splunkError(v.(map[string]interface{}), i, s)
      case interface{}:
         log.Debugln("_splunkError: interface{}")
      default:
         log.Debugln("_splunkError: default")
      }
   }
   return
}
