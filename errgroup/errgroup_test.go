package errgroup

import (
	"context"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestGroup(t *testing.T) {
	t.Parallel()

	g := WithContext(t.Context())
	ctx := g.Context()
	ctxCh := make(chan context.Context, 2)
	g.Go(func(ctx context.Context) error {
		ctxCh <- ctx
		return nil
	})
	g.TryGo(func(ctx context.Context) error {
		ctxCh <- ctx
		return nil
	})
	err := g.Wait()
	require.NoError(t, err)
	for range 2 {
		groupCtx := <-ctxCh
		require.Equal(t, ctx, groupCtx)
	}
}
