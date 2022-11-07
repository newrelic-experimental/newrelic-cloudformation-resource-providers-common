package model

// ErrorHandler decouples error handling from a possible implementation.
type ErrorHandler interface {
   TypeSpecificError(data *[]byte, s string) (err error)
   GetErrorMap(v interface{}, result map[string]interface{})
   ServerError(data *[]byte, s string) (err error)
}
