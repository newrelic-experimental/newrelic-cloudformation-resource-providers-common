package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/go-resty/resty/v2"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/configuration"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/logging"
   log "github.com/sirupsen/logrus"
)

type nerdgraph struct {
   client *resty.Client
   config *configuration.Config
}

// TODO abstract errors to insulate from different calls using different shapes
type errors struct {
   Description string `json:"description"`
   Type        string `json:"type"`
   Message     string `json:"message"`
}

type topLevelErrors struct {
   Errors []interface{} `json:"errors"`
}

// TODO Refactor Create & Update as they're almost the same

func NewClient(config *configuration.Config) *nerdgraph {
   return &nerdgraph{client: resty.New(), config: config}
}

// TODO refactor this as it is per-NerdGraph-type
const idField = "guid"

func (i *nerdgraph) emit(body string, apiKey string, apiEndpoint string) ([]byte, error) {
   log.Debugln("emit: body: ", body)
   log.Debugln("")

   bodyJson, err := json.Marshal(map[string]string{"query": body})
   if err != nil {
      return nil, err
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
      return nil, err
   }
   if resp.StatusCode() >= 300 {
      log.Errorf("Bad status code POSTing %s error: %s ", resp.Status(), bodyJson)
      err = fmt.Errorf("%s", resp.Status())
      return nil, err
   }

   respBody := resp.Body()
   logging.Dump(log.DebugLevel, string(respBody), "emit: response: ")

   err = hasErrors(&respBody)
   if err != nil {
      return nil, err
   }
   return respBody, nil
}
