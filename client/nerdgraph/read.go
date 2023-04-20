package nerdgraph

import (
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

type readResponse struct {
   Data readData `json:"data"`
}
type readData struct {
   Actor readActor `json:"actor"`
}
type readActor struct {
   Entity *readEntity `json:"entity,omitempty"`
}
type readEntity struct {
   Guid string `json:"guid"`
   Name string `json:"name"`
}

func (i *nerdgraph) Read(m model.Model) (err error) {
   defer func() {
      log.Debugf("Read: returning value: %+v type: %T", err, err)
   }()

   if m.GetGuid() == nil {
      log.Errorf("Read: missing guid")
      err = fmt.Errorf("%w Missing guid", &cferror.NotFound{})
      return
   }
   query := m.GetReadQuery()
   variables := m.GetVariables()
   i.config.InjectIntoMap(&variables)

   // Render the query
   query, err = model.Render(query, variables)
   if err != nil {
      log.Errorf("Read: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }
   log.Debugln("Read: rendered query: ", query)
   log.Debugln("")

   // Validate query
   err = model.Validate(&query)
   if err != nil {
      log.Errorf("Read: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   body, err := i.emit(query, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return
   }

   // Some NerdGraph APIs do not return an error on NOT FOUND, rather they return an empty result
   key := m.GetResultKey(model.Read)
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
