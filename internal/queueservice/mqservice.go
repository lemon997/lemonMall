package queueservice

import (
	"context"
)

type MQService struct {
	ctx context.Context
}

func MQNew(ctx context.Context) MQService {
	svc := MQService{ctx: ctx}
	return svc
}
