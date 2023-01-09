package repository

import (
	"github.com/couchbase/gocb/v2"
	"time"
)

const (
	BucketName   = "Bucket"
	QueryTimeout = 2 * time.Second
)

type Repository interface {
	Exist(id string) (bool, error)
	Insert(doc Document) error
}

type repository struct {
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

func NewRepository(cluster *gocb.Cluster) Repository {
	return &repository{
		cluster:    cluster,
		collection: cluster.Bucket(BucketName).DefaultCollection(),
	}
}

func (h *repository) Exist(id string) (bool, error) {
	exists, err := h.collection.Exists(id, nil)
	if err != nil {
		return false, err
	}
	return exists.Exists(), nil
}

func (h *repository) Insert(entity Document) error {
	_, err := h.collection.Insert(
		entity.Id,
		entity,
		&gocb.InsertOptions{
			Timeout: QueryTimeout,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

type Document struct {
	Id    string   `json:"id"`
	Field string   `json:"field"`
	Value string   `json:"value"`
	Cas   gocb.Cas `json:"cas"`
}
