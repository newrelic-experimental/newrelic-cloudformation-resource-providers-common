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
   log.Debugf("Read: Enter: model: %+v", m.GetResourceModel())
   log.Debugf("Read: Enter: variables: %+v", m.GetVariables())

   if m.GetIdentifier() == nil {
      log.Errorf("Read: missing identifier")
      err = fmt.Errorf("%w Missing identifier", &cferror.NotFound{})
      return
   }
   query := m.GetReadQuery()
   variables := m.GetVariables()
   i.config.InjectIntoMap(&variables)
   log.Debugf("Read: variables: %+v", variables)

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

   return i.resultHandler.Read(m, body)
}
