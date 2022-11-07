package nerdgraph

import (
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

func (i *nerdgraph) Update(m model.Model) (err error) {
   if err = i.Read(m); err != nil {
      return
   }
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

   key := m.GetResultKey(model.Update)
   if key != "" {
      var v interface{}
      v, err = FindKeyValue(body, key)
      if err != nil {
         log.Errorf("Update: error finding result key: %s in response: %s", key, string(body))
         return
      }

      if v == nil {
         log.Errorf("Update: result not returned by NerdGraph operation")
         err = fmt.Errorf("%w Update: result not returned by NerdGraph operation", &cferror.InvalidRequest{})
         return
      }
   }
   return
}
