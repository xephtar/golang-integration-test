package couchbase

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"runtime"
	"time"
)

type Container struct {
	container        testcontainers.Container
	containerRequest testcontainers.ContainerRequest
	address          string
	ip               string
	port             int
}

const Image = "docker.io/trendyoltech/couchbase-testcontainer:6.5.1"

func NewContainer(image string) *Container {

	exposedPort := []string{"8091:8091/tcp", "8093:8093/tcp", "11210:11210/tcp"}

	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: exposedPort,
		WaitingFor:   wait.ForLog("couchbase-dev started").WithStartupTimeout(45 * time.Second),
		Env:          map[string]string{"USERNAME": "Administrator", "PASSWORD": "password", "BUCKET_NAME": "Bucket"},
	}

	return &Container{
		containerRequest: req,
	}
}

func (c *Container) Run() (err error) {

	c.container, err = testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: c.containerRequest,
		Started:          true,
	})
	if err != nil {
		return err
	}

	c.ip, err = c.container.Host(context.Background())
	if err != nil {
		return err
	}

	if isRunningOnOSX() {
		c.ip = "127.0.0.1"
	}

	return nil
}

func (c *Container) Ip() string {
	return c.ip
}

func (c *Container) ForceRemoveAndPrune() (err error) {
	if c.container != nil {
		return c.container.Terminate(context.Background())
	}

	return nil
}

func isRunningOnOSX() bool {
	return runtime.GOOS == "darwin"
}
