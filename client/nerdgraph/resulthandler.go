package nerdgraph

//
// Provide a default results handler that covers the common endpoint guid/id patterns
//

import (
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

type ResultHandler struct {
}

func NewResultHandler() (h model.ResultHandler) {
   // Hmm- should we init with the model?
   h = &ResultHandler{}
   return
}

func (n ResultHandler) Create(m model.Model, body []byte) (err error) {
   key := m.GetIdentifierKey(model.Create)
   if key != "" {
      var v interface{}
      v, err = FindKeyValue(body, key)
      if err != nil {
         log.Errorf("Create: error finding result key: %s in response: %s", key, string(body))
         return err
      }
      s := fmt.Sprintf("%v", v)
      m.SetIdentifier(&s)
   }
   return
}

func (n ResultHandler) Delete(m model.Model, body []byte) (err error) {
   key := m.GetIdentifierKey(model.Delete)
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
   return
}

func (n ResultHandler) List(m model.Model, body []byte) (err error) {
   key := m.GetIdentifierKey(model.List)
   if key != "" {
      var guids []interface{}
      guids, err = findAllKeyValues(body, key)
      if err != nil {
         return
      }

      log.Debugf("List: guids: %+v", guids)
      for _, g := range guids {
         m.AppendToResourceModels(m.NewModelFromGuid(g))
      }
   } else {
      log.Errorf("No result expected for List, this is probably an error")
   }
   return
}

func (n ResultHandler) Read(m model.Model, body []byte) (err error) {
   key := m.GetIdentifierKey(model.Read)
   if key != "" {
      var v interface{}
      v, err = FindKeyValue(body, key)
      if err != nil {
         log.Errorf("error finding result key: %s in response: %s", key, string(body))
         err = fmt.Errorf("%w Not found key: %s", &cferror.NotFound{}, key)
         return
      }

      if v == nil {
         log.Errorf("Read: result not returned by NerdGraph operation")
         // err = fmt.Errorf("%w Read: result not returned by NerdGraph operation", &cferror.InvalidRequest{})
         err = fmt.Errorf("%w Not found key: %s value: %v", &cferror.NotFound{}, key, v)
         return
      }
   }
   return
}

func (n ResultHandler) Update(m model.Model, body []byte) (err error) {
   key := m.GetIdentifierKey(model.Update)
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
