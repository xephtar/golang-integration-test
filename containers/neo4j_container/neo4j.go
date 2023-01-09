package neo4j_container

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

const (
	Image    = "neo4j:4.3-community"
	Username = "neo4j"
	Password = "test"
)

func NewContainer(image string) *Container {

	exposedPort := []string{
		"7687/tcp",
		"7474:7474",
		"7473:7473",
		"7687:7687",
	}

	env := map[string]string{
		"NEO4J_AUTH": "neo4j/test",
	}

	req := testcontainers.ContainerRequest{
		Image:        image,
		Env:          env,
		ExposedPorts: exposedPort,
		WaitingFor:   wait.ForLog("Bolt enabled").WithStartupTimeout(time.Minute * 2),
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
