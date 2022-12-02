package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/logging"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
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
   Type       string `json:"type"`
}

func NewCommonErrorHandler(m model.Model) (eh model.ErrorHandler) {
   log.Debugf("errors.NewCommonErrorHandler: enter: model %p", m)
   defer func() {
      log.Debugf("errors.NewCommonErrorHandler: exit : %p", eh)
   }()
   eh = &CommonErrorHandler{M: m}
   return eh
}

// CommonErrorHandler implements ErrorHandler and provides a common, default, implementation of dealing with errors coming out of API calls.
type CommonErrorHandler struct {
   M model.Model
}

// HasErrors can't be a method on ErrorHandler due to the way Go handles (or doesn't) dispatch- no v-table :-(
func HasErrors(e model.ErrorHandler, data *[]byte) (err error) {
   log.Debugf("HasErrors: %p enter", e)
   defer func() {
      log.Debugf("HasErrors: %p returning %v", e, err)
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

   if err = e.ServerError(data, s); err != nil {
      return
   }

   if err = e.TypeSpecificError(data, s); err != nil {
      return
   }
   return
}

// TypeSpecificError is a bit complex, we don't know the shape so we have to travel a map[string]interface{}
func (e *CommonErrorHandler) TypeSpecificError(data *[]byte, s string) (err error) {
   // TODO drop 's' on next refactor
   log.Debugf("TypeSpecificError: %p enter", e)
   defer func() {
      log.Debugf("TypeSpecificError: %p returning %v", e, err)
   }()
   v, err := FindKeyValue(*data, "errors")
   log.Debugf("TypeSpecificError: found: %v %T", v, v)
   if err != nil {
      return
   }
   if v == nil {
      return
   }

   errorMap := make(map[string]interface{})
   e.GetErrorMap(v, errorMap)

   if errorMap == nil {
      log.Warnf("Empty errors array: %v+ %T", e, e)
      return
   }
   _type := fmt.Sprintf("%v", errorMap[e.M.GetErrorKey()])
   if strings.Contains(strings.ToLower(_type), "not_found") || strings.Contains(strings.ToLower(_type), "not found") {
      err = fmt.Errorf("%w Not found", &cferror.NotFound{})
      return
   }
   return
}

func (e *CommonErrorHandler) GetErrorMap(v interface{}, result map[string]interface{}) {
   log.Debugf("GetErrorMap: enter %p", e)
   defer func() {
      log.Debugf("GetErrorMap: exit %p", e)
   }()
   switch k := v.(type) {
   case []interface{}:
      for _, j := range k {
         e.GetErrorMap(j, result)
      }
   case map[string]interface{}:
      for key, value := range k {
         result[key] = value
      }
   default:
      log.Warnf("GetErrorMap: unknown value/type: %+v %T", k, k)
   }
   return
}

// ServerError is relatively simple, we know its shape
func (e *CommonErrorHandler) ServerError(data *[]byte, s string) (err error) {
   // TODO drop 's' on next refactor
   log.Debugf("ServerError: %p exit", e)
   defer func() {
      log.Debugf("ServerError: %p returning %v", e, err)
   }()

   r := genericRoot{}
   err = json.Unmarshal(*data, &r)
   if err != nil {
      log.Errorf("ServerError: unmarshal %v", err)
      err = fmt.Errorf("%w %s", &cferror.UnknownError{}, s)
      return
   }

   // Nothing in the error array
   if r.Errors == nil || len(r.Errors) <= 0 {
      return
   }

   // At this point we actually have something
   if len(r.Errors) > 1 {
      log.Warnf("ServerError: %d errors returned from NerdGraph, the first is used the remainder logged", len(r.Errors))
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
   var _type = r.Errors[0].Extensions.Type
   var errorType = r.Errors[0].Extensions.ErrorType
   msg := strings.ToLower(errorMessage + _type + errorType)
   log.Infof("ServerError: code: %d message: %s errorType: %s type: %s", errorCode, errorMessage, errorType, _type)

   if strings.Contains(msg, "not_found") || strings.Contains(msg, "not found") {
      err = fmt.Errorf("%w Not found", &cferror.NotFound{})
      return
   }

   // In-case we can't find a specific error
   if errorCode == 0 {
      log.Errorf("ServerError: non-specific error %s", s)
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
