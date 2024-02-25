package hopper

import (
	"fmt"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

const (
	defaultDBName = "default"
)

type M map[string]string

type Collection struct {
	*bbolt.Bucket
}

type Hopper struct {
	db *bbolt.DB
}

func New() (*Hopper, error) {
	dbname := fmt.Sprintf("%s.hopper", defaultDBName)
	db, err := bbolt.Open(dbname, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &Hopper{
		db: db,
	}, nil
}

func (h *Hopper) CreateCollection(name string) (*Collection, error) {
	tx, err := h.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// bucket := tx.Bucket([]byte(name))
	// if bucket != nil {
	// 	return &Collection{Bucket: bucket}, nil
	// }
	bucket, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	return &Collection{Bucket: bucket}, nil
}

func (h *Hopper) Insert(collName string, data M) (uuid.UUID, error) {
	id := uuid.New()
	tx, err := h.db.Begin(true)
	if err != nil {
		return id, err
	}
	defer tx.Rollback()

	bucket, err := tx.CreateBucketIfNotExists([]byte(collName))
	if err != nil {
		return id, err
	}
	for k, v := range data {
		if err := bucket.Put([]byte(k), []byte(v)); err != nil {
			return id, err
		}
	}
	if err := bucket.Put([]byte("id"), []byte(id.String())); err != nil {
		return id, err
	}
	return id, tx.Commit()
}

func (h *Hopper) Select(coll string, query M) (M, error) {
	tx, err := h.db.Begin(false)
	if err != nil {
		return nil, err
	}
	bucket := tx.Bucket([]byte(coll))
	if bucket == nil {
		return nil, fmt.Errorf("collection (%s) not found", coll)
	}
	return nil, nil
}
