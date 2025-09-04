package supervisor

import "context"

type StopToken struct{}

type Process func(ctx context.Context) error
