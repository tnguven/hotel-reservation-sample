package store

import "context"

type Dropper interface {
	Drop(context.Context) error
}
