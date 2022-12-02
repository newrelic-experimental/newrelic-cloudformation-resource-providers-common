package nerdgraph

import (
   "errors"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "time"
)

func (i *nerdgraph) Delete(m model.Model) (err error) {
   if err = i.Read(m); err != nil {
      return
   }
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

   key := m.GetResultKey(model.Delete)
   if key != "" {
      var v interface{}
      v, err = FindKeyValue(body, key)
      if err != nil {
         log.Errorf("error finding result key: %s in response: %s", key, string(body))
         return
      }

      if v == nil {
         log.Errorf("Delete: result not returned by NerdGraph operation")
         err = fmt.Errorf("%w Delete: result not returned by NerdGraph operation", &cferror.InvalidRequest{})
         return
      }
   }
   // Allow for the NRDB propagation delay by doing a spin Read
   // FUTURE add some sort of timeout interrupt (channel?)
   // Call Read until it returns an arror, hopefully NotFound
   // Allow for the NRDB propagation delay by doing a spin Read
   err = i.Read(m)
   for err == nil {
      err = i.Read(m)
      log.Debugf("common.Delete: spin lock: %+v", err)
      time.Sleep(1 * time.Second)
      // FUTURE add some sort of timeout interrupt
   }
   // Delete _wants_ to wait for NotFound, therefore return nil to indicate OK

   var nf *cferror.NotFound
   if err != nil && errors.As(err, &nf) {
      err = nil
   }
   delta := time.Now().Sub(start)
   log.Debugf("DeleteMutation: exit: err: %+v propagation delay: %v", err, delta)
   return
}
