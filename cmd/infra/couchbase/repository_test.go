package repository_test

import (
	repository "golang-project-with-intregtation-test/cmd/infra/couchbase"

	testifyAssert "github.com/stretchr/testify/assert"
)

func (s *testSuite) Test_it_should_insert_document() {
	//given
	var (
		assert = testifyAssert.New(s.T())
	)

	doc := repository.Document{
		Id:    "123",
		Field: "field",
		Value: "value",
		Cas:   0,
	}
	defer s.delete(doc)

	repo := repository.NewRepository(s.CouchbaseCluster)

	// when
	err := repo.Insert(doc)

	// then
	assert.Nil(err)
	retrievedDoc := s.get("123")
	assert.Equal(doc.Id, retrievedDoc.Id)
}

func (s *testSuite) Test_it_should_return_true_if_document_exists() {
	//given
	var (
		assert = testifyAssert.New(s.T())
	)

	doc := repository.Document{
		Id:    "123",
		Field: "field",
		Value: "value",
		Cas:   0,
	}
	defer s.delete(doc)
	s.insert(doc)

	repo := repository.NewRepository(s.CouchbaseCluster)

	// when
	exist, err := repo.Exist("123")

	// then
	assert.Nil(err)
	assert.True(exist)
}

func (s *testSuite) Test_it_should_return_false_if_document_does_not_exist() {
	//given
	var (
		assert = testifyAssert.New(s.T())
	)

	repo := repository.NewRepository(s.CouchbaseCluster)

	// when
	exist, err := repo.Exist("not_exist")

	// then
	assert.Nil(err)
	assert.False(exist)
}
