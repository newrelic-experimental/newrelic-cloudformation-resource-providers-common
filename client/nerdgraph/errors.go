package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/logging"
   log "github.com/sirupsen/logrus"
   "strings"
)

type genericRoot struct {
   Errors []genericError `json:"errors,omitempty"`
}

type genericError struct {
   Extensions extensions    `json:"extensions"`
   Stacktrace []string      `json:"stacktrace"`
   Locations  []interface{} `json:"locations"`
   Message    string        `json:"message"`
   Path       []string      `json:"path"`
}

type extensions struct {
   ErrorClass string `json:"errorClass"`
   ErrorCode  int    `json:"errorCode"`
   ErrorType  string `json:"errorType"`
}

func hasErrors(data *[]byte) (err error) {
   defer func() {
      log.Debugf("hassErrors: returning %v", err)
   }()
   // Empty
   if data == nil {
      return
   }

   // No error keyword
   s := string(*data)
   if !(strings.Contains(strings.ToLower(s), `"error":`) || strings.Contains(strings.ToLower(s), `"errors":`)) {
      return
   }

   if err = serverError(data, s); err != nil {
      return
   }

   if err = typeSpecificError(data, s); err != nil {
      return
   }
   return
}

// typeSpecific error is a bit complex, we don't know the shape so we have to travel a map[string]interface{}
func typeSpecificError(data *[]byte, s string) (err error) {
   // TODO
   defer func() {
      log.Debugf("typeSpecificError: returning %v", err)
   }()

   return
}

// serverError is relatively simple, we know its shape
func serverError(data *[]byte, s string) (err error) {
   defer func() {
      log.Debugf("serverError: returning %v", err)
   }()

   r := genericRoot{}
   err = json.Unmarshal(*data, &r)
   if err != nil {
      log.Errorf("serverError: unmarshal %v", err)
      err = fmt.Errorf("%w %s", &cferror.UnknownError{}, s)
      return
   }

   // Nothing in the error array
   if r.Errors == nil || len(r.Errors) <= 0 {
      return
   }

   // At this point we actually have something
   if len(r.Errors) > 1 {
      log.Warnf("serverError: %d errors returned from NerdGraph, the first is used the remainder logged", len(r.Errors))
   }
   for i, e := range r.Errors {
      // Don't log the first error, we'll return it as the error value
      if i == 0 {
         continue
      }
      logging.Dump(log.ErrorLevel, e, "serverError: NerdGraph error")
   }
   var errorCode = r.Errors[0].Extensions.ErrorCode
   var errorMessage = r.Errors[0].Message
   log.Infof("serverError: code: %d message: %s", errorCode, errorMessage)

   // In-case we can't find a specific error
   if errorCode == 0 {
      log.Errorf("serverError: non-specific error %s", s)
      err = fmt.Errorf("%w %s", &cferror.UnknownError{}, s)
      return
   }

   switch errorCode {
   case 404:
      err = fmt.Errorf("%w Not found", &cferror.NotFound{})
   default:
      err = fmt.Errorf("%w code: %d message: %s", &cferror.UnknownError{}, errorCode, errorMessage)
   }
   return
}
