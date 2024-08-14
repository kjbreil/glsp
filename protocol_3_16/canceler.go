package protocol

import (
	"context"
)

type cancelReturn struct {
	rtn any
	err error
}

func cancelErr(ctx context.Context, fn func() error) error {
	wfn := func() (any, error) {
		err := fn()
		return nil, err
	}
	_, err := cancelRtnErr(ctx, wfn)
	return err
}

func cancelRtnErr(ctx context.Context, fn func() (any, error)) (any, error) {
	ch := make(chan cancelReturn, 1)
	defer close(ch)
	go func() {
		cr := cancelReturn{}
		cr.rtn, cr.err = fn()
		ctxErr := ctx.Err()
		if ctxErr == nil {
			ch <- cr
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case rtn := <-ch:
		return rtn.rtn, rtn.err
	}
}
