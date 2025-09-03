package supervisor

type StopToken struct{}

type Process func(stop <-chan StopToken) error
