package cferror

import (
   "github.com/aws/aws-sdk-go/service/cloudformation"
)

type InvalidRequest struct {
   Err error
}

func (e *InvalidRequest) Error() string {
   return cloudformation.HandlerErrorCodeInvalidRequest
}

func (e *InvalidRequest) Unwrap() error {
   return e.Err
}

type UnknownError struct {
   Err error
}

func (e *UnknownError) Error() string {
   return cloudformation.HandlerErrorCodeGeneralServiceException
}

func (e *UnknownError) Unwrap() error {
   return e.Err
}

type NotFound struct {
   Err error
}

func (e *NotFound) Error() string {
   return cloudformation.HandlerErrorCodeNotFound
}

func (e *NotFound) Unwrap() error {
   return e.Err
}

type AlreadyExists struct {
   Err error
}

func (e *AlreadyExists) Error() string {
   return cloudformation.HandlerErrorCodeAlreadyExists
}

func (e *AlreadyExists) Unwrap() error {
   return e.Err
}

type ServiceInternalError struct {
   Err error
}

func (e *ServiceInternalError) Error() string {
   return cloudformation.HandlerErrorCodeServiceInternalError
}

func (e *ServiceInternalError) Unwrap() error {
   return e.Err
}
