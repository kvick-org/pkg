package errgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	*errgroup.Group
	ctx context.Context
}

func WithContext(ctx context.Context) *Group {
	g, gCtx := errgroup.WithContext(ctx)
	return &Group{
		Group: g,
		ctx:   gCtx,
	}
}

func (g *Group) Go(f func(ctx context.Context) error) {
	g.Group.Go(func() error {
		return f(g.ctx)
	})
}

func (g *Group) TryGo(f func(ctx context.Context) error) {
	g.Group.TryGo(func() error {
		return f(g.ctx)
	})
}

func (g *Group) Context() context.Context {
	return g.ctx
}
