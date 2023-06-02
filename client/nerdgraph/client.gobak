package nerdgraph

import (
   "encoding/json"
   "errors"
   "fmt"
   "github.com/go-resty/resty/v2"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/configuration"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/logging"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "time"
)

type nerdgraph struct {
   client        *resty.Client
   config        *configuration.Config
   errorHandler  model.ErrorHandler
   resultHandler model.ResultHandler
}

func NewClient(config *configuration.Config, errorHandler model.ErrorHandler, resultHandler model.ResultHandler) *nerdgraph {
   log.Debugf("client.NewClient: errorHandler: %p", errorHandler)
   // FIXME this can be a singleton based on typeName as the same type will use the same errorHandler
   return &nerdgraph{client: resty.New(), config: config, errorHandler: errorHandler, resultHandler: resultHandler}
}

func (i *nerdgraph) emit(body string, apiKey string, apiEndpoint string) (respBody []byte, err error) {
   defer func() {
      log.Debugf("client.emit exit: errorHandler: %p %+v", i.errorHandler, i.errorHandler)
   }()
   log.Debugf("client.emit enter: errorHandler: %p %+v", i.errorHandler, i.errorHandler)
   log.Debugln("emit: body: ", body)
   log.Debugln("")

   bodyJson, err := json.Marshal(map[string]string{"query": body})
   if err != nil {
      return
   }

   headers := map[string]string{
      "Content-Type": "application/json",
      "Api-Key":      apiKey,
      "deep-trace":   "true",
      "User-Agent":   i.config.GetUserAgent(),
   }
   log.Debugf("emit: headers: %+v", headers)
   type PostResult interface {
   }
   type PostError interface {
   }
   var postResult PostResult
   var postError PostError

   var timeout *cferror.Timeout
   retry := true
   for retry {
      var resp *resty.Response
      resp, err = i.client.R().
         SetBody(bodyJson).
         SetHeaders(headers).
         SetResult(&postResult).
         SetError(&postError).
         Post(apiEndpoint)

      if err != nil {
         log.Errorf("Error POSTing %v", err)
         return
      }
      if resp.StatusCode() >= 300 {
         log.Errorf("Bad status code POSTing %s error: %s ", resp.Status(), bodyJson)
         err = fmt.Errorf("%s", resp.Status())
         return
      }

      respBody = resp.Body()
      logging.Dump(log.DebugLevel, string(respBody), "emit: response: ")

      err = HasErrors(i.errorHandler, &respBody)
      // NOTE: This spin lock must be sync due to create having to return the guid first try
      if errors.As(err, &timeout) {
         log.Warnf("emit: retrying due to timeout %v", err)
         time.Sleep(1 * time.Second)
      } else {
         retry = false
      }
   }
   return
}

// func captureResult(m model.Model, action model.Action, body []byte) (err error) {
//    for _, key := range m.GetCaptureKeys(action) {
//       var v interface{}
//       v, err = FindKeyValue(body, key)
//       if err != nil {
//          log.Errorf("Create: error finding result key: %s in response: %s", key, string(body))
//          return
//       }
//       s := fmt.Sprintf("%v", v)
//       if action == model.Create {
//          setModel(m, key, s)
//       }
//    }
//    return
// }

// func setModel(m model.Model, n string, v string) {
//    // pointer to struct - addressable
//    ps := reflect.ValueOf(m)
//    pv := reflect.ValueOf(v)
//    // struct
//    s := ps.Elem()
//    if s.Kind() == reflect.Struct {
//       // exported field
//       f := s.FieldByName(n)
//       if f.IsValid() {
//          // A Value can be changed only if it is
//          // addressable and was not obtained by
//          // the use of unexported struct fields.
//          if f.CanSet() {
//             // change value of N
//             if f.Kind() == pv.Kind() {
//                f.SetString(v)
//             }
//          }
//       }
//    }
// }
