package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/go-resty/resty/v2"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/configuration"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/logging"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

type nerdgraph struct {
   client *resty.Client
   config *configuration.Config
   model  model.Model
}

func NewClient(config *configuration.Config, model model.Model) *nerdgraph {
   return &nerdgraph{client: resty.New(), config: config, model: model}
}

func (i *nerdgraph) emit(body string, apiKey string, apiEndpoint string) (respBody []byte, err error) {
   log.Debugln("emit: body: ", body)
   log.Debugln("")

   bodyJson, err := json.Marshal(map[string]string{"query": body})
   if err != nil {
      return
   }

   headers := map[string]string{"Content-Type": "application/json", "Api-Key": apiKey, "deep-trace": "true"}
   log.Debugf("emit: headers: %+v", headers)
   type PostResult interface {
   }
   type PostError interface {
   }
   var postResult PostResult
   var postError PostError

   resp, err := i.client.R().
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

   err = i.hasErrors(&respBody)
   return
}
