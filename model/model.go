package model

type Model interface {
   // GetResourceModel Return the generated resource.Model
   GetResourceModel() interface{}
   // GetResourceModels Return the generated resource.Model(s) as an array for List
   GetResourceModels() []interface{}
   // GetGraphQL get the GraphQL fragment from resource.Model
   GetGraphQLFragment() *string
   // SetGuid set the guid in the generated resource.Model
   SetGuid(g *string)
   // GetGuid get the guid in the generated resource.Model
   GetGuid() *string
   // GetCreateMutation get the GraphQL mutation for Create
   GetCreateMutation() string
   // GetDeleteMutation get the GraphQL mutation for Delete
   GetDeleteMutation() string
   // GetUpdateMutation get the GraphQL mutation for Update
   GetUpdateMutation() string
   // GetReadQuery get the GraphQL query for Read
   GetReadQuery() string
   // GetListQuery get the GraphQL query for List
   GetListQuery() string
   // GetListQueryNextCursor get the GraphQL query for pagination, NEXTCURSOR is the template param
   GetListQueryNextCursor() string
   // GetGuidKey the response key that has the guid/id/pk
   GetGuidKey() string
   // GetListQueryFilter GraphQL entitySearch query string (optional) from resource.Model
   // GetListQueryFilter() *string
   // GetVariables return moustache substitution variables from resource.Model
   GetVariables() map[string]string
   // AppendToResourceModels adds a resource.Model to resource.Model.ResourceModels
   AppendToResourceModels(m Model)
   // NewModelFromGuid creates a new Model with the passed guid
   NewModelFromGuid(g interface{}) Model
}
