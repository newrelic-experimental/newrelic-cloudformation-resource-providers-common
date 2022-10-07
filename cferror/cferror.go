package cferror

import (
   "encoding/json"
   "fmt"
   "github.com/aws/aws-sdk-go/service/cloudformation"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/logging"
   log "github.com/sirupsen/logrus"
   "strings"
)

type InvalidRequest struct {
   Err error
}

func (e *InvalidRequest) Error() string {
   return cloudformation.HandlerErrorCodeInvalidRequest
}

func (e *InvalidRequest) Unwrap() error {
   return e.Err
}

type UnknownError struct {
   Err error
}

func (e *UnknownError) Error() string {
   return cloudformation.HandlerErrorCodeGeneralServiceException
}

func (e *UnknownError) Unwrap() error {
   return e.Err
}

type NotFound struct {
   Err error
}

func (e *NotFound) Error() string {
   return cloudformation.HandlerErrorCodeNotFound
}

func (e *NotFound) Unwrap() error {
   return e.Err
}

type AlreadyExists struct {
   Err error
}

func (e *AlreadyExists) Error() string {
   return cloudformation.HandlerErrorCodeAlreadyExists
}

func (e *AlreadyExists) Unwrap() error {
   return e.Err
}

type ServiceInternalError struct {
   Err error
}

func (e *ServiceInternalError) Error() string {
   return cloudformation.HandlerErrorCodeServiceInternalError
}

func (e *ServiceInternalError) Unwrap() error {
   return e.Err
}

// TODO abstract this as it's type specific
type workloadError struct {
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

type root struct {
   Errors []map[string]interface{} `json:"errors,omitempty"`
}

func HasErrors(data *[]byte) (err error) {
   // Empty
   if data == nil {
      return
   }

   // No error keyword
   s := string(*data)
   if !strings.Contains(strings.ToLower(s), "error") {
      return
   }

   r := root{}
   err = json.Unmarshal(*data, &r)
   if err != nil {
      log.Errorf("hasError: unmarshal %v", err)
      return fmt.Errorf("%w %s", &UnknownError{}, s)
   }

   // Nothing in the error array
   if r.Errors == nil || len(r.Errors) <= 0 {
      return
   }

   // At this point we actually have something
   if len(r.Errors) > 1 {
      log.Warnf("%d errors returned from NerdGraph, the first is used the remainder logged", len(r.Errors))
   }
   for i, e := range r.Errors {
      // Don't log the first error, we'll return it as the error value
      if i == 0 {
         continue
      }
      logging.Dump(log.ErrorLevel, e, "NerdGraph error")
   }
   var errorCode = 0
   var errorMessage = ""
   _splunkError(r.Errors[0], &errorCode, &errorMessage)
   log.Infof("code: %d message: %s", errorCode, errorMessage)

   // In-case we can't find a specific error
   if errorCode == 0 {
      log.Errorf("hasError: non-specific error %s", s)
      return fmt.Errorf("%w %s", &UnknownError{}, s)
   }

   switch errorCode {
   case 404:
      err = fmt.Errorf("%w Not found", &NotFound{})
   default:
      err = fmt.Errorf("%w code: %d message: %s", &UnknownError{}, errorCode, errorMessage)
   }
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
