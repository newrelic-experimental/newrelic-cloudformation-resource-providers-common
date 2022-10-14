package client

import (
   "errors"
   "fmt"
   "github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
   "github.com/aws/aws-sdk-go/aws/session"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/client/nerdgraph"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/configuration"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

/*
   Contract adherence here

   Do not modify the model!
*/

type IClient interface {
   Create(model model.Model) error
   Delete(model model.Model) error
   Update(m model.Model) error
   Read(m model.Model) error
   List(m model.Model) error
}

type GraphqlClient struct {
   client IClient
}

var graphqlClient *GraphqlClient

func NewGraphqlClient(sess *session.Session, typeName *string, model model.Model) *GraphqlClient {
   if graphqlClient == nil {
      graphqlClient = &GraphqlClient{
         client: nerdgraph.NewClient(configuration.NewConfiguration(sess, typeName), model),
      }
      log.Debugln("NewGraphqlClient: returning NerdGraph client")
   }
   return graphqlClient
}

func (i *GraphqlClient) CreateMutation(model model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.Create(model)

   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "Create complete",
         ResourceModel:   model.GetResourceModel(),
      }
   } else {
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          fmt.Sprintf("Create error: %s", err.Error()),
      }
   }
   return event, nil
}

func (i *GraphqlClient) DeleteMutation(model model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.Delete(model)
   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "Delete complete",
      }
   } else {
      fmt.Printf("DeleteMutation: error: %+v", err)
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          fmt.Sprintf("Delete error: %s", err.Error()),
      }
   }
   return event, nil
}

func (i *GraphqlClient) UpdateMutation(model model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.Update(model)
   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "Update complete",
         ResourceModel:   model.GetResourceModel(),
      }
   } else {
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          fmt.Sprintf("Update error: %s", err.Error()),
         ResourceModel:    model.GetResourceModel(),
      }
   }
   return event, nil
}

func (i *GraphqlClient) ReadQuery(model model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.Read(model)
   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "Read complete",
         ResourceModel:   model.GetResourceModel(),
      }
   } else {
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          fmt.Sprintf("Read error: %s", err.Error()),
      }
   }
   return event, nil
}

func (i *GraphqlClient) ListQuery(model model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.List(model)
   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "List complete",
         ResourceModels:  model.GetResourceModels(),
      }
   } else {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         fmt.Sprintf("List error: %s", err.Error()),
         ResourceModels:  []interface{}{},
      }
   }
   return event, nil
}
