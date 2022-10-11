package nerdgraph

import (
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "time"
)

func (i *nerdgraph) Delete(m model.Model) (err error) {
   variables := m.GetVariables()
   i.config.InjectIntoMap(&variables)
   mutation := m.GetDeleteMutation()

   // Render the mutation
   mutation, err = model.Render(mutation, variables)
   if err != nil {
      log.Errorf("Delete: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }
   log.Debugln("Delete: rendered mutation: ", mutation)
   log.Debugln("")

   // Validate mutation
   err = model.Validate(&mutation)
   if err != nil {
      log.Errorf("Delete: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   start := time.Now()
   body, err := i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   _ = body
   if err != nil {
      return
   }

   // v, err := findKeyValue(body, m.GetGuidKey())
   // if err != nil {
   //    log.Errorf("error finding guid: %s in response: %s", m.GetGuidKey(), string(body))
   //    return
   // }

   // if v == nil {
   //    log.Errorf("Delete: guid not returned by NerdGraph operation")
   //    err = fmt.Errorf("%w Delete: guid not returned by NerdGraph operation", &cferror.InvalidRequest{})
   //    return
   // }

   // Allow for the NRDB propagation delay by doing a spin Read
   // FUTURE add some sort of timeout interrupt (channel?)
   // Call Read until it returns Not Found
   err = nil
   for err == nil {
      err = i.Read(m)
   }
   delta := time.Now().Sub(start)
   log.Debugf("DeleteMutation: propagation delay: %v", delta)
   return nil
}
