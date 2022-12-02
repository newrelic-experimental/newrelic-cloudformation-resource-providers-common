package nerdgraph

import (
   "errors"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "time"
)

func (i *nerdgraph) Create(m model.Model) (err error) {
   variables := m.GetVariables()
   i.config.InjectIntoMap(&variables)
   mutation := m.GetCreateMutation()

   // Render the mutation
   mutation, err = model.Render(mutation, variables)
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }
   log.Debugln("Create: rendered mutation: ", mutation)
   log.Debugln("")

   // Validate mutation
   err = model.Validate(&mutation)
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   start := time.Now()
   body, err := i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return err
   }

   key := m.GetResultKey(model.Create)
   if key != "" {
      var v interface{}
      v, err = FindKeyValue(body, key)
      if err != nil {
         log.Errorf("Create: error finding result key: %s in response: %s", key, string(body))
         return err
      }
      s := fmt.Sprintf("%v", v)
      m.SetGuid(&s)
   }

   // Allow for the NRDB propagation delay by doing a spin Read
   err = i.Read(m)
   var nf *cferror.NotFound
   for err != nil && errors.As(err, &nf) {
      err = i.Read(m)
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
