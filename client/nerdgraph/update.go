package nerdgraph

import (
   "errors"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "time"
)

// TODO Update should *only* update

func (i *nerdgraph) Update(m model.Model) (err error) {
   log.Debugf("nerdgraph.Update: enter")
   // TODONE move this up one level, it's part of the contract
   // if err = i.Read(m); err != nil {
   //    return
   // }

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

   start := time.Now()
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
   // TODO move the spin-lock up one level, don't wait for it, return IN-PROGRESS
   // Allow for the NRDB propagation delay by doing a spin Read
   err = i.Read(m)
   var nf *cferror.NotFound
   for err != nil && errors.As(err, &nf) {
      err = i.Read(m)

      var timeout *cferror.Timeout
      if errors.As(err, &timeout) {
         log.Warnf("Update: retrying due to timeout %v", err)
         err = nil
      }

      log.Debugf("common.Create: spin lock: %+v", err)
      time.Sleep(1 * time.Second)
      // FUTURE add some sort of timeout interrupt
   }
   // Delete _wants_ to wait for NotFound, therefore return nil to indicate OK
   if err != nil && errors.As(err, &nf) {
      err = nil
   }
   delta := time.Now().Sub(start)
   log.Debugf("CreateMutation: exit: err: %+v propagation delay: %v", err, delta)
   return
}
