package tasks

import "context"

type TaskAction func(ctx context.Context) error
