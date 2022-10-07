package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "time"
)

type deleteResponse struct {
   Data   deleteData      `json:"data"`
   Errors []workloadError `json:"errors,omitempty"`
}
type deleteData struct {
   Payload payload `json:"workloadDelete"`
}

var delay = time.Second * 35

func (i *nerdgraph) Delete(m model.Model) error {
   if m == nil {
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, "nil model")
   }
   if m.GetGuid() == nil {
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, "nil guid")
   }

   mutation, err := model.Render(m.GetDeleteMutation(), map[string]string{"GUID": *m.GetGuid()})
   if err != nil {
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   // There seems to be a nrdb propagation delay on delete
   start := time.Now()
   log.Debugf("Delete: guid: %s model: %+v", *m.GetGuid(), m)
   body, err := i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return err
   }

   response := deleteResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   if response.Data.Payload.Guid == "" {
      log.Errorf("Delete: guid not returned by NerdGraph operation")
      err = fmt.Errorf("%w Delete: guid not returned by NerdGraph operation", &cferror.InvalidRequest{})
      return err
   }
   // Sleep to allow the guid deletion to propagate through nrdb
   // delta := delay - (time.Now().Sub(start))
   // log.Debugf("DeleteMutation: sleeping: %v", delta)
   // time.Sleep(delta)

   // Call Read until it returns Not Found
   err = nil
   for err == nil {
      err = i.Read(m)
   }
   delta := time.Now().Sub(start)
   log.Debugf("DeleteMutation: propagation delay: %v", delta)
   return nil
}

const deleteMutation = `
mutation {
  workloadDelete(guid: "{{{GUID}}}") {
    guid
  }
}
`
