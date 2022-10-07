package nerdgraph

import (
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

type listResults struct {
   Entities   []listEntities `json:"entities"`
   NextCursor string         `json:"nextCursor"`
}
type listEntities struct {
   Guid string `json:"guid"`
   Name string `json:"name"`
}

// List only gets 30 seconds to do its work, IN_PROGRESS is not allowed
// NOTE: entitySearch requires several seconds to index a newly created entity. Read the guid in the model and append it to the list result.
func (i *nerdgraph) List(m model.Model) ([]interface{}, error) {
   log.Debugf("List: enter: guid: %s", *m.GetGuid())
   result := make([]interface{}, 0)

   // Because of the indexing delay on the guid after create do an entity query
   err := i.Read(m)
   if err != nil {
      return result, err
   }
   result = append(result, m)

   filter := ""
   if m.GetListQueryFilter() != nil {
      filter = *m.GetListQueryFilter()
   }
   mutation, err := model.Render(listQuery, map[string]string{"LISTQUERYFILTER": filter})
   if err != nil {
      return nil, fmt.Errorf("%w %s", &cferror.InvalidRequest{}, err.Error())
   }

   _, err = i.emit(mutation, *i.config.APIKey, i.config.GetEndpoint())
   if err != nil {
      return nil, err
   }

   // TODO
   // response := listResponse{}
   // err = json.Unmarshal(body, &response)
   // if err != nil {
   //    return result, err
   // }
   //
   // // NOTE: entitySearch does not return errors
   // for _, e := range response.Data.Actor.EntitySearch.Results.Entities {
   //    result = append(result, &model.Model{
   //       Guid:   &e.Guid,
   //       APIKey: m.APIKey,
   //    })
   // }
   // TODO process cursor
   return result, nil
}

const listQuery = `
{
  actor {
    entitySearch(queryBuilder: {type: WORKLOAD}) {
      count
      results {
        nextCursor
        entities {
            guid
            name
        }
      }
    }
  }
}
`
const listQueryNextCursor = `
{
  actor {
    entitySearch(queryBuilder: {type: WORKLOAD}) {
      count
      results(cursor: "{{{NEXTCURSOR}}}") {
        nextCursor
        entities {
            guid
            name
        }
      }
    }
  }
}
`
