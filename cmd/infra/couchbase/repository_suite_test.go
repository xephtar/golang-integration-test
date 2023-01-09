package repository_test

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	testifyAssert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	repository "golang-project-with-intregtation-test/cmd/infra/couchbase"
	"golang-project-with-intregtation-test/containers/couchbase"
	"os"
	"testing"
	"time"
)

type testSuite struct {
	suite.Suite
	CouchbaseContainer *couchbase.Container
	CouchbaseCluster   *gocb.Cluster
	collection         *gocb.Collection
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) SetupSuite() {
	err := s.setupCouchbase()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "repository couchbase setup failed: %s\n", err)
		s.T().FailNow()
		return
	}
}

func (s *testSuite) TearDownSuite() {

	err := s.CouchbaseContainer.ForceRemoveAndPrune()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "repository couchbase terminate failed: %s\n", err)
	}

}

func (s *testSuite) setupCouchbase() (err error) {

	s.CouchbaseContainer = couchbase.NewContainer(couchbase.Image)
	err = s.CouchbaseContainer.Run()
	if err != nil {
		return err
	}

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		},
	}

	s.CouchbaseCluster, err = gocb.Connect("couchbase://"+s.CouchbaseContainer.Ip(), opts)
	if err != nil {
		return err
	}

	err = s.CouchbaseCluster.WaitUntilReady(60*time.Second, nil)
	if err != nil {
		return err
	}

	s.collection = s.CouchbaseCluster.Bucket(repository.BucketName).DefaultCollection()

	return nil
}

func (s *testSuite) insert(doc repository.Document) {
	_, err := s.collection.Insert(doc.Id, doc, nil)
	testifyAssert.Nil(s.T(), err)
}

func (s *testSuite) exist(id string) bool {
	received, err := s.collection.Exists(id, nil)
	testifyAssert.Nil(s.T(), err)

	return received.Exists()
}

func (s *testSuite) delete(docs ...repository.Document) {
	for _, doc := range docs {
		s.collection.Remove(doc.Id, nil)
	}
}

func (s *testSuite) get(id string) (doc repository.Document) {
	received, err := s.collection.Get(id, nil)
	testifyAssert.Nil(s.T(), err)

	err = received.Content(&doc)
	testifyAssert.Nil(s.T(), err)
	doc.Cas = received.Cas()
	return doc
}
