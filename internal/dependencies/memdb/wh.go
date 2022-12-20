package memdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"golang.org/x/exp/slices"
)

func whTypeToTable(whType int) string {
	switch whType {
	case domain.WhTypeMutation:
		return "mutation"
	case domain.WhTypeSpell:
		return "spell"
	default:
		panic("invalid whType")
	}
}

type WhDbService struct {
	Db *memdb.MemDB
}

func NewWhDbService() *WhDbService {
	db, err := createNewWhMemDb()
	if err != nil {
		panic(err)
	}

	return &WhDbService{Db: db}
}

func createNewWhMemDb() (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"wh": {
				Name: "wh",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
		},
	}
	return memdb.NewMemDB(schema)
}

func (s *WhDbService) Retrieve(ctx context.Context, whType int, whId string, users []string, sharedUsers []string) (*domain.Wh, *domain.DbError) {
	txn := s.Db.Txn(false)
	whRaw, err1 := txn.First("wh", "id", whId)
	if err1 != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err1}
	}

	if whRaw == nil {
		return nil, &domain.DbError{Type: domain.DbNotFoundError, Err: errors.New("wh not found")}
	}

	wh, ok := whRaw.(*domain.Wh)
	if !ok {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: fmt.Errorf("could not populate wh from raw %v", whRaw)}
	}

	if slices.Contains(users, wh.OwnerId) {
		return wh, nil
	}

	if slices.Contains(sharedUsers, wh.OwnerId) && wh.Shared {
		return wh, nil
	}

	return nil, &domain.DbError{Type: domain.DbNotFoundError, Err: errors.New("wh not found")}

}

func (s *WhDbService) Create(ctx context.Context, whType int, w *domain.Wh) (*domain.Wh, *domain.DbError) {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err2 := txn.Insert("wh", w); err2 != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err2}
	}
	txn.Commit()

	return w.Copy(), nil
}
