package threading

import "context"

type CancellationToken interface {
	Context() context.Context
	Cancel()
}

type cancellationToken struct {
	context context.Context
	cancel  context.CancelFunc
}

func (token *cancellationToken) Context() context.Context {
	return token.context
}

func (token *cancellationToken) Cancel() {
	token.cancel()
}

func NewToken() CancellationToken {
	return NewTokenFrom(context.TODO())
}

func NewTokenFrom(ctx context.Context) CancellationToken {
	ctx, cancel := context.WithCancel(ctx)

	return &cancellationToken{
		context: ctx,
		cancel:  cancel,
	}
}
