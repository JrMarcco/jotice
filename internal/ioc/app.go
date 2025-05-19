package ioc

import "context"

type Task interface {
	Start(ctx context.Context)
}

type App struct {
	Tasks []Task
}
