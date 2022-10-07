package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

type workloadError struct {
}
type updateResponse struct {
   Data   updateData      `json:"data"`
   Errors []workloadError `json:"errors,omitempty"`
}
type updateData struct {
   Payload payload `json:"workloadUpdate"`
}

func (i *nerdgraph) Update(m model.Model) error {
   log.Debugf("nerdgraph/model.Update model: %+v", m)
   if m == nil {
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, "nil model")
   }
   if m.GetGuid() == nil {
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, "nil guid")
   }
   // Render the mutation
   // TODO abstract the substitution map
   mutation, err := model.Render(updateMutation, map[string]string{"GUID": *m.GetGuid(), "WORKLOAD": *m.GetGraphQL()})
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }
   log.Debugf("Update- rendered mutation: %s", mutation)

   // Validate mutation
   err = model.Validate(&mutation)
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   log.Debugf("Update: mutation: %s\n model: %+v", mutation, m)
   body, err := i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return err
   }

   response := updateResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   // if err = hasError(response.Errors); err != nil {
   //    return err
   // }
   if response.Data.Payload.Guid == "" {
      log.Errorf("Update: guid not returned by NerdGraph operation")
      err = fmt.Errorf("%w Update: guid not returned by NerdGraph operation", &cferror.InvalidRequest{})
      return err
   }
   return nil
}

const updateMutation = `
mutation {
  workloadUpdate(guid: "{{{GUID}}}", {{{WORKLOAD}}}) {
    guid
  }
}
`
