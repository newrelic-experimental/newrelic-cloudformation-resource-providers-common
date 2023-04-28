package client

import (
   "errors"
   "fmt"
   "github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/client/nerdgraph"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/configuration"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/tags"
   log "github.com/sirupsen/logrus"
   "sync"
)

/*
   Contract adherence here

   1. Do not modify the model!
   2. TODO If we're returning a non-nil handler.ProgressEvent then RETURN a nil error!
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
   req    handler.Request
}

// FIXME this should not be a singleton!
// var graphqlClient *GraphqlClient

// func NewGraphqlClient(sess *session.Session, typeName *string, model model.Model, errorHandler model.ErrorHandler) *GraphqlClient {
func NewGraphqlClient(req handler.Request, typeName *string, model model.Model, errorHandler model.ErrorHandler) *GraphqlClient {
   log.Debugf("client.NewGraphqlClient: enter: errorHandler: %p", errorHandler)
   // if graphqlClient == nil {
   //    graphqlClient = &GraphqlClient{
   //       client: nerdgraph.NewClient(configuration.NewConfiguration(sess, typeName), model, errorHandler),
   //    }
   //    log.Debugln("NewGraphqlClient: returning NerdGraph client")
   // }
   // return graphqlClient
   return &GraphqlClient{
      client: nerdgraph.NewClient(configuration.NewConfiguration(req.Session, typeName), model, errorHandler),
      req:    req}
}

func (i *GraphqlClient) CreateMutation(model model.Model) (event handler.ProgressEvent, err error) {
   // If we're returning a valid ProgressEvent DO NOT return an error
   defer func() {
      if event.OperationStatus != "" {
         err = nil
      }
   }()

   // Test for and return a queued event/err pair if any
   evt, err := getEvent(i.req)
   if evt != nil {
      return *evt, err
   }

   // Create the entity
   err = i.client.Create(model)

   if err == nil {
      if model.HasTags() {
         event = handler.ProgressEvent{
            OperationStatus: handler.InProgress,
            Message:         "Create waiting on tagging",
            ResourceModel:   model.GetResourceModel(),
         }
         setEvent(i.req, event, err)
         // Run tagging in a goroutine. Because this is running inside a Lambda wrapper provided by the go plugin there's no worry of the "main" exiting- it's async all the way down.
         go func() {
            tagModel := tags.NewTagModel(model.GetGuid(), model.GetTags(), model.GetVariables())
            sm := tags.NewPayload(tagModel)
            c := NewGraphqlClient(i.req, &tags.TypeName, sm, tags.NewErrorHandler(sm))

            var evt2 handler.ProgressEvent
            er := c.client.Create(sm)
            if er != nil {
               evt2 = handler.ProgressEvent{
                  OperationStatus:  handler.Failed,
                  HandlerErrorCode: errors.Unwrap(er).Error(),
                  Message:          fmt.Sprintf("Create error: %s", er.Error()),
               }
            } else {
               evt2 = handler.ProgressEvent{
                  OperationStatus: handler.Success,
                  Message:         "Create complete",
                  ResourceModel:   model.GetResourceModel(),
               }
            }
            setEvent(i.req, evt2, err)
         }()
      } else {
         // No tags, no error, we're done. Return Success
         event = handler.ProgressEvent{
            OperationStatus: handler.Success,
            Message:         "Create complete",
            ResourceModel:   model.GetResourceModel(),
         }
      }
   } else {
      // Error creating the entity. Return Failed
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          fmt.Sprintf("Create error: %s", err.Error()),
      }
   }
   return event, nil
}

var mu sync.Mutex
var events = make(map[string]*handler.ProgressEvent)
var errs = make(map[string]error)

func setEvent(req handler.Request, event handler.ProgressEvent, err error) {
   defer mu.Unlock()
   mu.Lock()
   events[req.LogicalResourceID] = &event
   errs[req.LogicalResourceID] = err
}

func getEvent(req handler.Request) (*handler.ProgressEvent, error) {
   defer mu.Unlock()
   mu.Lock()
   evt := events[req.LogicalResourceID]
   if evt == nil {
      return nil, nil
   }

   err := errs[req.LogicalResourceID]

   if evt.OperationStatus == handler.Success {
      delete(events, req.LogicalResourceID)
      delete(errs, req.LogicalResourceID)
   }

   return evt, err
}

func (i *GraphqlClient) DeleteMutation(model model.Model) (event handler.ProgressEvent, err error) {
   // If we're returning a valid ProgressEvent DO NOT return an error
   defer func() {
      if event.OperationStatus != "" {
         err = nil
      }
   }()

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
   // If we're returning a valid ProgressEvent DO NOT return an error
   defer func() {
      if event.OperationStatus != "" {
         err = nil
      }
   }()

   log.Debugf("client.UpdateMutation: enter")

   // Test for and return a queued event/err pair if any
   evt, err := getEvent(i.req)
   if evt != nil {
      return *evt, nil
   }

   // Verify the entity exists as this is an update
   if err = i.client.Read(model); err != nil {
      log.Debugf("UpdateMutation: client.Read: HandlerErrorCode: %v", errors.Unwrap(err))
      log.Debugf("Update mutation: Failed 1")
      return handler.ProgressEvent{
         OperationStatus: handler.Failed,
         //         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message: fmt.Sprintf("Update error: %s", err.Error()),
         //         ResourceModel:    model.GetResourceModel(),
      }, nil
   }

   err = i.client.Update(model)

   if err == nil {
      if model.HasTags() {
         event = handler.ProgressEvent{
            OperationStatus: handler.InProgress,
            Message:         "Update waiting on tagging",
            ResourceModel:   model.GetResourceModel(),
         }
         setEvent(i.req, event, err)
         // Run tagging in a goroutine. Because this is running inside a Lambda wrapper provided by the go plugin there's no worry of the "main" exiting- it's async all the way down.
         go func() {
            tagModel := tags.NewTagModel(model.GetGuid(), model.GetTags(), model.GetVariables())
            sm := tags.NewPayload(tagModel)
            c := NewGraphqlClient(i.req, &tags.TypeName, sm, tags.NewErrorHandler(sm))

            var evt2 handler.ProgressEvent
            er := c.client.Update(sm)
            if er != nil {
               log.Debugf("Update mutation: Failed 2")
               evt2 = handler.ProgressEvent{
                  OperationStatus:  handler.Failed,
                  HandlerErrorCode: errors.Unwrap(er).Error(),
                  //                  ResourceModel:    model.GetResourceModel(),
                  Message: fmt.Sprintf("Update error: %s", er.Error()),
               }
            } else {
               evt2 = handler.ProgressEvent{
                  OperationStatus: handler.Success,
                  Message:         "Update complete",
                  ResourceModel:   model.GetResourceModel(),
               }
            }
            setEvent(i.req, evt2, err)
         }()
      } else {
         // No tags, no error, we're done. Return Success
         event = handler.ProgressEvent{
            OperationStatus: handler.Success,
            Message:         "Update complete",
            ResourceModel:   model.GetResourceModel(),
         }
      }
   } else {
      log.Debugf("Update mutation: Failed 3")
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          fmt.Sprintf("Update error: %s", err.Error()),
         //         ResourceModel:    model.GetResourceModel(),
      }
   }
   return event, nil
}

func (i *GraphqlClient) ReadQuery(model model.Model) (event handler.ProgressEvent, err error) {
   // If we're returning a valid ProgressEvent DO NOT return an error
   defer func() {
      if event.OperationStatus != "" {
         err = nil
      }
   }()

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
   // If we're returning a valid ProgressEvent DO NOT return an error
   defer func() {
      if event.OperationStatus != "" {
         err = nil
      }
   }()

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
