package nerdgraph

import (
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

// List only gets 30 seconds to do its work, IN_PROGRESS is not allowed
// NOTE: entitySearch requires several seconds to index a newly created entity. Read the guid in the model and append it to the list result.
func (i *nerdgraph) List(m model.Model) (err error) {
   // Because of the indexing delay on the guid after create do an entity query
   err = i.Read(m)
   if err != nil {
      return
   }
   // Add current model to result list
   m.AppendToResourceModels(m)

   variables := m.GetVariables()
   i.config.InjectIntoMap(&variables)
   mutation := m.GetListQuery()

   // Render the mutation
   mutation, err = model.Render(mutation, variables)
   if err != nil {
      log.Errorf("List: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }
   log.Debugln("List: rendered mutation: ", mutation)
   log.Debugln("")

   // Validate mutation
   err = model.Validate(&mutation)
   if err != nil {
      log.Errorf("List: %v", err)
      return fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   body, err := i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return
   }

   err = i.resultHandler.List(m, body)
   if err != nil {
      return
   }
   // TODO process cursor
   // DOC By convention NEXTCURSOR is the field to substitute in the template
   return
}
