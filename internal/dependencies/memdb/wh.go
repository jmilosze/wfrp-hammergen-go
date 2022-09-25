package memdb

import (
	"github.com/hashicorp/go-memdb"
)

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
			"mutation": {
				Name: "mutation",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
			"spell": {
				Name: "spell",
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

//func getOneWh[W domain.WhType](db *memdb.MemDB, table string, fieldName string, fieldValue string) (*W, *domain.DbError) {
//	txn := db.Txn(false)
//	whRaw, err := txn.First(table, fieldName, fieldValue)
//	if err != nil {
//		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
//	}
//
//	if whRaw == nil {
//		return nil, &domain.DbError{Type: domain.DbNotFoundError, Err: errors.New("user not found")}
//	}
//	user := whRaw.(*domain.UserDb)
//
//	return copyWh(user), nil
//}

//func copyWh[W domain.WhType](from *W) *W {
//	switch v := any(from).(type) {
//	case domain.Mutation:
//
//	}
//}
