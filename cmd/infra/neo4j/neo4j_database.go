package neo4j

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	existQuery  = "MATCH(node:Person) WHERE node.personId = $personIdValue RETURN node"
	mergeQuery  = "MERGE (node:Person {personId: $personIdValue})"
	personIdKey = "personIdValue"
)

type databaseImp struct {
	Driver       neo4j.Driver
	DatabaseName string
}

type Database interface {
	InsertPerson(personId int) error
	Exist(personId int) (bool, error)
}

func NewService(neo4jDriver neo4j.Driver, databaseName string) Database {
	return &databaseImp{
		Driver:       neo4jDriver,
		DatabaseName: databaseName,
	}
}

func (service *databaseImp) InsertPerson(personId int) error {
	session, err := service.createSession()
	if err != nil {
		return err
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

	if err != nil {
		return err
	}

	return nil
}

func (service *databaseImp) Exist(personId int) (bool, error) {
	session, err := service.createSession()
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

func (service *databaseImp) createSession() (neo4j.Session, error) {
	return service.Driver.NewSession(neo4j.SessionConfig{DatabaseName: service.DatabaseName})
}
