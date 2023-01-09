package neo4j_test

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/suite"
	"golang-project-with-intregtation-test/containers/neo4j_container"
	"os"
	"testing"
)

type testSuite struct {
	suite.Suite
	Neo4jContainer *neo4j_container.Container
	Neo4jGraph     neo4j.Driver
	GraphName      string `default:"neo4j"`
}

const (
	existQuery  = "MATCH(node:Person) WHERE node.personId = $personIdValue RETURN node"
	mergeQuery  = "MERGE (node:Person {personId: $personIdValue})"
	personIdKey = "personIdValue"
)

func TestNeo4j(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) SetupSuite() {
	err := s.setupNeo4j()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "repository neo4j setup failed: %s\n", err)
		s.T().FailNow()
		return
	}
}

func (s *testSuite) TearDownSuite() {

	err := s.Neo4jContainer.ForceRemoveAndPrune()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "repository neo4j terminate failed: %s\n", err)
	}

}

func (s *testSuite) setupNeo4j() (err error) {
	s.Neo4jContainer = neo4j_container.NewContainer(neo4j_container.Image)

	err = s.Neo4jContainer.Run()
	if err != nil {
		return
	}
	s.Neo4jGraph, err = s.createNeo4jDriver()
	if err != nil {
		return err
	}

	err = s.Neo4jGraph.VerifyConnectivity()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return
}

func (s *testSuite) createNeo4jDriver() (neo4j.Driver, error) {
	return neo4j.NewDriver("bolt://"+s.Neo4jContainer.Ip(), neo4j.BasicAuth(neo4j_container.Username, neo4j_container.Password, ""), func(c *neo4j.Config) { c.Encrypted = false })
}

func (s *testSuite) insert(personId int) (err error) {
	session, err := s.createSession()
	if err != nil {
		return
	}
	defer session.Close()

	params := map[string]interface{}{
		personIdKey: personId,
	}

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(mergeQuery, params)

		if err != nil {
			return nil, err
		}

		return result.Consume()
	})

	return
}

func (s *testSuite) exist(personId int) (bool, error) {
	session, err := s.createSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	params := map[string]interface{}{
		personIdKey: personId,
	}

	result, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		records, err := transaction.Run(existQuery, params)

		if err != nil {
			return nil, err
		}

		return records.Next(), nil
	})

	return result.(bool), err
}

func (s *testSuite) createSession() (neo4j.Session, error) {
	return s.Neo4jGraph.NewSession(neo4j.SessionConfig{DatabaseName: s.GraphName})
}
