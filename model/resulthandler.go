package model

// ResultHandler magic is due to Go's embedding:  https://go.dev/doc/effective_go#embedding  and https://genekuo.medium.com/composition-with-structures-and-method-forwarding-in-go-9a9045f2c1fd
type ResultHandler interface {
   Create(m Model, b []byte) (err error)
   Delete(m Model, b []byte) (err error)
   List(m Model, b []byte) (err error)
   Read(m Model, b []byte) (err error)
   Update(m Model, b []byte) (err error)
}
