package supervisor

import "context"

type Process func(ctx context.Context) error
