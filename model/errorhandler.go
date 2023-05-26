package model

// ErrorHandler decouples error handling from a possible implementation.
// ErrorHandler magic is due to Go's embedding:  https://go.dev/doc/effective_go#embedding   https://genekuo.medium.com/composition-with-structures-and-method-forwarding-in-go-9a9045f2c1fd
type ErrorHandler interface {
   TypeSpecificError(data *[]byte, s string) (err error)
   GetErrorMap(v interface{}, result map[string]interface{})
   ServerError(data *[]byte, s string) (err error)
}
