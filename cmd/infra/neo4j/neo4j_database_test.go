package neo4j_test

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"golang-project-with-intregtation-test/cmd/infra/neo4j"
)

func (s *testSuite) Test_it_should_return_true_if_person_exist() {
	//given
	var (
		assert = testifyAssert.New(s.T())
	)
	personId := 123
	s.insert(personId)

	db := neo4j.NewService(s.Neo4jGraph, s.GraphName)

	// when
	exist, err := db.Exist(personId)

	// then
	assert.Nil(err)
	assert.True(exist)
}

func (s *testSuite) Test_it_should_insert_person() {
	//given
	var (
		assert = testifyAssert.New(s.T())
	)
	personId := 123

	db := neo4j.NewService(s.Neo4jGraph, s.GraphName)

	// when
	err := db.InsertPerson(personId)

	// then
	assert.Nil(err)
	exist, err := s.exist(personId)
	assert.Nil(err)
	assert.True(exist)
}
