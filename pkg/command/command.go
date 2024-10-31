package command

import (
	"context"
)

type Name string

func (c Name) String() string { return string(c) }

type Command struct {
	processor ProcessorFn
}

type ProcessorFn func(context.Context, *Message) (*Message, []Name)

type MiddlewareFn func(fn ProcessorFn) ProcessorFn

func New(processor ProcessorFn, middleware ...MiddlewareFn) *Command {
	for _, m := range middleware {
		processor = m(processor)
	}

	return &Command{
		processor: processor,
	}
}

func (c *Command) Process(
	ctx context.Context, msg *Message,
) (*Message, []Name) {
	return c.processor(ctx, msg)
}
