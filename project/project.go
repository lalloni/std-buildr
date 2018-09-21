package project

import "context"

type Publisher interface {
	Publish(ctx *context.Context) error
}

type Cleaner interface {
	Clean(ctx *context.Context) error
}
