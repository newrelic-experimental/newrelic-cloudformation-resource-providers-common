package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
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

func (i *nerdgraph) Read(m model.Model) error {
   if m.GetGuid() == nil {
      return fmt.Errorf("%w %s", &cferror.NotFound{}, "missing guid")
   }
   if i.config.APIKey == nil {
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, "missing APIKey")
   }

   mutation, err := model.Render(readQuery, map[string]string{"GUID": *m.GetGuid()})
   if err != nil {
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   body, err := i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return err
   }

   response := readResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   // An entity query returns no errors so nothing to check
   if response.Data.Actor.Entity == nil {
      return fmt.Errorf("%w %s", &cferror.NotFound{}, "guid not found")
   }

   return nil
}

const readQuery = `
{
  actor {
    entity(guid: "{{{GUID}}}") {
        guid
        name
    }
  }
}
`
