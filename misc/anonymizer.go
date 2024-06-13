package misc

import (
	"context"
	"io"
)

type Anonymizer interface {
	Run(context.Context, io.Writer) error
}
