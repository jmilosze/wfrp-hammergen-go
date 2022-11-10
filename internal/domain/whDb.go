package domain

import "context"

type WhDbService interface {
	Create(ctx context.Context, wh Warhammer) (Warhammer, *DbError)
	Retrieve(ctx context.Context, whId string, users []string, sharedUsers []string) (Warhammer, *DbError)
}
