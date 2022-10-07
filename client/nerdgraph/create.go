package nerdgraph

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/logging"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "time"
)

type createResponse struct {
   Data   createData      `json:"data"`
   Errors []workloadError `json:"errors,omitempty"`
}
type createData struct {
   Payload payload `json:"workloadCreate"`
}
type payload struct {
   Guid string `json:"guid"`
}

func (i *nerdgraph) Create(m model.Model) error {
   log.Debugf("nerdgraph/client.Create model: %+v", m)
   // TODO abstract the rendering out of the "framework"

   variables := map[string]string{"ACCOUNTID": *i.config.AccountID}
   mutation := m.GetCreateMutation()
   variables["WORKLOAD"] = *m.GetGraphQL()

   // Render the mutation
   // TODO abstract the substitution map
   mutation, err := model.Render(mutation, variables)
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

   // There seems to be a nrdb propagation delay on create
   start := time.Now()
   body, err := i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return err
   }

   response := createResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }
   logging.Dump(log.DebugLevel, response, "Create: response: ")
   if response.Data.Payload.Guid == "" {
      log.Errorf("Create: guid not returned by NerdGraph operation")
      err = fmt.Errorf("%w Create: guid not returned by NerdGraph operation", &cferror.InvalidRequest{})
      return err
   }
   m.SetGuid(&response.Data.Payload.Guid)

   err = fmt.Errorf("placeholder")
   for err != nil {
      err = i.Read(m)
   }
   delta := time.Now().Sub(start)
   log.Debugf("CreateMutation: propagation delay: %v", delta)
   return nil
}

const createMutation = `
mutation {
  workloadCreate(accountId: {{{ACCOUNTID}}}, {{{WORKLOAD}}}) {
    guid
  }
}
`
