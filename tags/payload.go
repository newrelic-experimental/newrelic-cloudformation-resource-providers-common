package tags

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "strings"
)

//
// Generic, should be able to leave these as-is
//

// Muy importante Read only cares about error/non-error!

type Payload struct {
   model  *Model
   models []interface{}
   // preReadTags []TagObject
}

func (p *Payload) SetIdentifier(g *string) {
   // TODO implement me
   panic("implement me")
}

func (p *Payload) GetIdentifier() *string {
   return p.model.Guid
}

func (p *Payload) GetTagIdentifier() *string {
   return p.model.Guid
}

func (p *Payload) GetIdentifierKey(a model.Action) string {
   // We don't get a guid/id back from tag api's
   return ""
}

func (p *Payload) HasTags() bool {
   return p.model.Tags != nil
}

func (p *Payload) GetTags() map[string]string {
   return make(map[string]string)
}
func NewPayload(m *Model) (p *Payload) {
   // // Should only be true on Create
   // if m.Guid == nil {
   //    m.Guid = m.EntityGuid
   // }

   p = &Payload{
      model:  m,
      models: make([]interface{}, 0),
   }
   // p.processTagString()

   return
}

// func (p *Payload) processTagString() {
//    defer func() {
//       log.Debugf("processTagString: tagString: %s", *p.model.TagString)
//    }()
//
//    // Create
//    if p.model.TagString == nil {
//       log.Debugf("processTagString: Create")
//       p.model.TagString = tagsToString(p.model.Tags)
//       return
//    }
//
//    // Delete or Read
//    if len(p.model.Tags) <= 0 {
//       log.Debugf("processTagString: Delete/Read")
//       copy(p.model.Tags, tagsFromString(p.model.TagString))
//       // FIXME- Delete should return an empty TagString
//       return
//    }
//
//    // Must be an Update
//    log.Debugf("processTagString: Update")
//    // For pre-read, required by an actual stack
//    // TODO how to use this?
//    p.preReadTags = tagsFromString(p.model.TagString)
//    if p.model.Variables == nil {
//       p.model.Variables = make(map[string]string)
//    }
//    for _, v := range p.preReadTags {
//       p.model.Variables["ORIGINALTAG"] = *v.Key
//       p.model.Variables["ORIGINALVALUE"] = v.Values[0]
//    }
//    // For result
//    p.model.TagString = tagsToString(p.model.Tags)
// }

// func tagsFromString(tagString *string) (t []TagObject) {
//    if tagString == nil {
//       log.Errorf("tagsFromString: nil tagString")
//       return
//    }
//    err := json.Unmarshal([]byte(*tagString), &t)
//    if err != nil {
//       log.Errorf("tagsFromString: json.Unmarshal: %v", err)
//    }
//    return
// }
//
// func tagsToString(tags []TagObject) (s *string) {
//    ba, err := json.Marshal(tags)
//    if err != nil {
//       log.Errorf("tagsToString: json.Marshal: %v", err)
//       return
//    }
//
//    ts := string(ba)
//    return &ts
// }

func (p *Payload) GetResourceModel() interface{} {
   return p.model
}

func (p *Payload) GetResourceModels() []interface{} {
   log.Debugf("GetResourceModels: returning %+v", p.models)
   return p.models
}

// Tagging doesn't implement List
func (p *Payload) AppendToResourceModels(m model.Model) {
   p.models = append(p.models, m.GetResourceModel())
}

//
// These are API specific, must be configured per API
//

var TypeName = "NewRelic::Observability::Tagging"

func (p *Payload) NewModelFromGuid(g interface{}) (m model.Model) {
   s := fmt.Sprintf("%s", g)
   return NewPayload(&Model{Guid: &s})
}

var emptyFragment = ""

func (p *Payload) GetGraphQLFragment() *string {
   return &emptyFragment
}

func (p *Payload) SetGuid(g *string) {
   p.model.Guid = g
   log.Debugf("SetIdentifier: %s", *p.model.Guid)
}

func (p *Payload) GetGuid() *string {
   return p.model.Guid
}

/*
func (p *Payload) GetGuidKey() string {
   // FIXME Only List returns a guid, this causes all other calls to fail
   return "guid"
}
*/

func (p *Payload) GetResultKey(a model.Action) string {
   switch a {
   case model.List:
      return "guid"
   case model.Read:
      return "guid"
   }
   return ""
}

func (p *Payload) GetVariables() map[string]string {
   // ACCOUNTID comes from the configuration
   // NEXTCURSOR is a _convention_

   log.Debugf("GetVariables: p.model.Guid: %s", *p.model.Guid)
   log.Debugf("GetVariables: p.model.Variables: %v", p.model.Variables)
   log.Debugf("GetVariables: p.model.Tags: %+v", p.model.Tags)

   if p.model.Variables == nil {
      p.model.Variables = make(map[string]string)
   }

   // if p.model.EntityGuid != nil {
   //    p.model.Variables["GUID"] = *p.model.EntityGuid
   // }
   if p.model.Guid != nil {
      p.model.Variables["GUID"] = *p.model.Guid
   }

   if p.model.Tags != nil {
      // JSON Stringify the tags
      ba, err := json.Marshal(p.model.Tags)
      if err != nil {
         panic(err)
      }
      log.Debugf("GetVariables: ba: %+v", ba)

      // Fix the case + GraphQL/JSON snafu
      tagString := string(ba)
      log.Debugf("GetVariables: tagString pre : %s", tagString)
      // FIXME remove qoutes from all instances of "key": and "values":   It's a GraphQL'ism
      tagString = strings.ReplaceAll(tagString, `"Key":`, `key:`)
      tagString = strings.ReplaceAll(tagString, `"Values":`, `values:`)
      log.Debugf("GetVariables: tagString post: %s", tagString)
      // tagString = strings.ReplaceAll(tagString, `"Key"`, "key")
      // tagString = strings.ReplaceAll(tagString, `"Values"`, "values")
      p.model.Variables["TAGS"] = tagString

      // Build an array of just the keys.
      // This blocked needed by Read & Delete
      // FIXME dash-case keys need a backtick, document that for the end-user
      keys := make([]string, 0)
      for _, t := range p.model.Tags {
         keys = append(keys, t.Key)
         if len(t.Values) > 0 {
            p.model.Variables["TAG"] = t.Key
            p.model.Variables["VALUE"] = t.Values[0]
         }
      }
      // JSON Stringify the key array
      ba, err = json.Marshal(keys)
      if err != nil {
         panic(err)
      }
      p.model.Variables["KEYS"] = string(ba)
   }

   lqf := ""
   // if p.model.ListQueryFilter != nil {
   //    lqf = *p.model.ListQueryFilter
   // }
   p.model.Variables["LISTQUERYFILTER"] = lqf

   return p.model.Variables
}

func (p *Payload) GetErrorKey() string {
   return "type"
}

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
   Domain     string `json:"domain"`
   EntityType string `json:"entityType"`
   Guid       string `json:"guid"`
   Name       string `json:"name"`
   Tags       []tag  `json:"tags,omitempty"`
   Type       string `json:"type"`
}
type tag struct {
   Key    string   `json:"key"`
   Values []string `json:"values"`
}

// func (p *Payload) TestReadResponse(data []byte) (err error) {
//    r := readResponse{}
//    if err = json.Unmarshal(data, &r); err != nil {
//       return
//    }
//    if p.model.Tags == nil && r.Data.Actor.Entity.Tags == nil {
//       return
//    }
//    if p.model.Tags != nil && r.Data.Actor.Entity.Tags == nil {
//       err = fmt.Errorf("%w model Tags nil and read tags not nil", &cferror.NotFound{})
//       return
//    }
//    if p.model.Tags == nil && r.Data.Actor.Entity.Tags != nil {
//       err = fmt.Errorf("%w model Tags not nil and read tags nil", &cferror.NotFound{})
//       return
//    }
//    tagsEqual := true
//    // Compare the model to the read result. Everything in the model must be present in the read- the read may be bigger
//    for _, modelTag := range p.model.Tags {
//       // If the model tag key is in the read tag array
//       if readValue, ok := containsTagKey(*modelTag.Key, r.Data.Actor.Entity.Tags); ok {
//          // Test the model key's values and ensure each is in the read tag's value array
//          for _, modelValue := range modelTag.Values {
//             if containsTagValue(modelValue, readValue) {
//                continue
//             } else {
//                tagsEqual = false
//                break
//             }
//          }
//       } else {
//          tagsEqual = false
//       }
//    }
//    if tagsEqual {
//       return
//    }
//    err = fmt.Errorf("%w model Tags not nil and read tags nil", &cferror.NotFound{})
//    return
// }

// func containsTagValue(rv string, mvs []string) bool {
//    for _, mv := range mvs {
//       if rv == mv {
//          return true
//       }
//    }
//    return false
// }
//
// func containsTagKey(key string, tags []tag) ([]string, bool) {
//    for _, t := range tags {
//       if t.Key == key {
//          return t.Values, true
//       }
//    }
//    return nil, false
// }

func (p *Payload) GetCreateMutation() string {
   return `
mutation {
  taggingAddTagsToEntity(guid: "{{{GUID}}}" tags: {{{TAGS}}} ) {
    errors {
      message
      type
    }
  }
}
`
}

func (p *Payload) GetDeleteMutation() string {
   return `
mutation {
  taggingDeleteTagFromEntity(guid: "{{{GUID}}}", tagKeys: {{{KEYS}}}) {
    errors {
      message
      type
    }
  }
}
`
}

func (p *Payload) GetUpdateMutation() string {
   return `
mutation {
  taggingReplaceTagsOnEntity(guid: "{{{GUID}}}", tags: {{{TAGS}}}) {
    errors {
      message
      type
    }
  }
}
`
}

// func (p *Payload) GetReadQueryOriginal() string {
//    return `
// {
//   actor {
//     entity(guid: "{{{GUID}}}") {
//       tags {
//          key
//          values
//       }
//       domain
//       entityType
//       guid
//       name
//       type
//     }
//   }
// }
// `
// }

func (p *Payload) GetReadQuery() string {
   if len(p.model.Tags) <= 0 {
      return `
{
  actor {
    entitySearch(query: "id = '{{{GUID}}}'") {
      results {
        entities {
          guid
        }
      }
    }
  }
}
`

   } else {
      return `
   {
     actor {
       entitySearch(query: "id = '{{{GUID}}}' and tags.{{{TAG}}} = '{{{VALUE}}}'") {
         results {
           entities {
             guid
           }
         }
       }
     }
   }
   `
   }
}

// func (p *Payload) GetPreUpdateReadQuery() string {
//    if len(p.model.Tags) <= 0 {
//       return `
// {
//   actor {
//     entitySearch(query: "id = '{{{GUID}}}'") {
//       results {
//         entities {
//           guid
//         }
//       }
//     }
//   }
// }
// `
//
//    } else {
//       return `
// {
//   actor {
//     entitySearch(query: "id = '{{{GUID}}}' and tags.{{{ORIGINALTAG}}} = '{{{ORIGINALVALUE}}}'") {
//       results {
//         entities {
//           guid
//         }
//       }
//     }
//   }
// }
// `
//    }
// }

func (p *Payload) GetListQuery() string {
   return `
{
  actor {
    entitySearch(query: "tags is null") {
      results {
        entities {
          guid
          entityType
          tags {
            key
            values
          }
        }
        nextCursor
      }
    }
  }
}
`
}

func (p *Payload) GetListQueryNextCursor() string {
   return `
{
  actor {
    entitySearch(query: "tags is null") {
      results(cursor: "{{{NEXTCURSOR}}}") {
        entities {
          guid
          entityType
          tags {
            key
            values
          }
        }
        nextCursor
      }
    }
  }
}
`
}
