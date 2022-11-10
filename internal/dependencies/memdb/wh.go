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

func whFromRaw(raw any, whType int) (domain.Warhammer, error) {
	var wh domain.Warhammer
	var ok bool

	switch whType {
	case domain.WhTypeMutation:
		wh, ok = raw.(*domain.Mutation)
	case domain.WhTypeSpell:
		wh, ok = raw.(*domain.Spell)
	default:
		ok = false
	}
	if !ok {
		return nil, fmt.Errorf("could not populate wh from raw %v", raw)
	} else {
		return wh, nil
	}

}

type WhDbService struct {
	Db     *memdb.MemDB
	WhType int
}

func NewWhDbService(whType int) *WhDbService {
	table := whTypeToTable(whType)
	db, err := createNewWhMemDb(table)
	if err != nil {
		panic(err)
	}

	return &WhDbService{Db: db, WhType: whType}
}

func createNewWhMemDb(table string) (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			table: {
				Name: table,
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

func (s *WhDbService) Retrieve(ctx context.Context, whId string, users []string, sharedUsers []string) (domain.Warhammer, *domain.DbError) {
	txn := s.Db.Txn(false)
	table := whTypeToTable(s.WhType)
	whRaw, err1 := txn.First(table, "id", whId)
	if err1 != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err1}
	}

	if whRaw == nil {
		return nil, &domain.DbError{Type: domain.DbNotFoundError, Err: errors.New("wh not found")}
	}

	wh, err2 := whFromRaw(whRaw, s.WhType)
	if err2 != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err2}
	}

	whFields := wh.GetCommonFields()

	if slices.Contains(users, whFields.OwnerId) {
		return wh, nil
	}

	if slices.Contains(sharedUsers, whFields.OwnerId) && whFields.Shared {
		return wh, nil
	}

	return nil, &domain.DbError{Type: domain.DbNotFoundError, Err: errors.New("wh not found")}

}

func (s *WhDbService) Create(ctx context.Context, w domain.Warhammer) (domain.Warhammer, *domain.DbError) {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	table := whTypeToTable(s.WhType)
	if err2 := txn.Insert(table, w); err2 != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err2}
	}
	txn.Commit()

	return w.Copy(), nil
}
