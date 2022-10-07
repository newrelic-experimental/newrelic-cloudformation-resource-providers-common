package model

type Model interface {
   GetResourceModel() interface{}
   GetResourceModels() []interface{}
   GetGraphQL() *string
   SetGuid(g *string)
   GetGuid() *string
   GetCreateMutation() string
   GetDeleteMutation() string
   GetUpdateMutation() string
   GetReadQuery() string
   GetListQuery() string
   GetCreateResponse() interface{}
   GetDeleteResponse() interface{}
   GetUpdateResponse() interface{}
   GetReadResponse() interface{}
   GetListResponse() interface{}
   GetListQueryFilter() *string
}
