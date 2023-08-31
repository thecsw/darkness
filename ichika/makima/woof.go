package makima

import "io"

type Woof interface {
	Parse() Woof
	Export() Woof
	Result() (string, io.Reader)
}
