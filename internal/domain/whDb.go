package domain

import (
	"context"
	"fmt"
)

type WhDbService[W WhType] interface {
	Create(ctx context.Context, wh *W) (*W, *DbError)
	Retrieve(ctx context.Context, fieldName string, fieldValue string) (*W, *DbError)
	RetrieveAll(ctx context.Context) ([]*W, *DbError)
	Update(ctx context.Context, user *W) (*W, *DbError)
	Delete(ctx context.Context, id string) *DbError
}

func GetTableName[W WhType](x W) (string, error) {
	switch v := any(x).(type) {
	case Mutation:
		return "mutation", nil
	case Spell:
		return "spell", nil
	default:
		return "", fmt.Errorf("could get table name of %T", v)
	}
}
