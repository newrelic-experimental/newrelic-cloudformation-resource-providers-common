package nerdgraph

import (
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

func (i *nerdgraph) Update(m model.Model) (err error) {
   variables := m.GetVariables()
   i.config.InjectIntoMap(&variables)
   mutation := m.GetUpdateMutation()

   // Render the mutation
   mutation, err = model.Render(mutation, variables)
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }
   log.Debugln("Update: rendered mutation: ", mutation)
   log.Debugln("")

   // Validate mutation
   err = model.Validate(&mutation)
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   body, err := i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return
   }

   v, err := findKeyValue(body, m.GetGuidKey())
   if err != nil {
      log.Errorf("Update: error finding guid: %s in response: %s", m.GetGuidKey(), string(body))
      return
   }

   if v == nil {
      log.Errorf("Update: guid not returned by NerdGraph operation")
      err = fmt.Errorf("%w Update: guid not returned by NerdGraph operation", &cferror.InvalidRequest{})
      return
   }
   return
}
