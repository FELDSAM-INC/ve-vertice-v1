package dockertest
/*
import (
	"errors"
	"testing"

	"github.com/fsouza/go-dockerclient"
	dtesting "github.com/fsouza/go-dockerclient/testing"
	"gopkg.in/check.v1"

)

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&S{})

type S struct{}

func (s *S) SetUpTest(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	defer conn.Close()
	dbtest.ClearAllCollections(conn.Apps().Database)
}

func (s *S) TestNewFakeDockerProvisioner(c *check.C) {
	server, err := dtesting.NewServer("127.0.0.1:0", nil, nil)
	c.Assert(err, check.IsNil)
	defer server.Stop()
	p, err := NewFakeDockerProvisioner(server.URL())
	c.Assert(err, check.IsNil)
	_, err = p.storage.RetrieveNode(server.URL())
	c.Assert(err, check.IsNil)
	opts := docker.PullImageOptions{Repository: "tsuru/bs"}
	err = p.Cluster().PullImage(opts, p.RegistryAuthConfig())
	c.Assert(err, check.IsNil)
	client, err := docker.NewClient(server.URL())
	c.Assert(err, check.IsNil)
	_, err = client.InspectImage("tsuru/bs")
	c.Assert(err, check.IsNil)
}

func (s *S) TestStartMultipleServersCluster(c *check.C) {
	p, err := StartMultipleServersCluster()
	c.Assert(err, check.IsNil)
	err = p.Cluster().PullImage(docker.PullImageOptions{Repository: "tsuru/bs"}, p.RegistryAuthConfig())
	c.Assert(err, check.IsNil)
	nodes, err := p.Cluster().Nodes()
	c.Assert(err, check.IsNil)
	c.Assert(nodes, check.HasLen, 2)
}

func (s *S) TestDestroy(c *check.C) {
	p, err := StartMultipleServersCluster()
	c.Assert(err, check.IsNil)
	p.Destroy()
	c.Assert(p.servers, check.IsNil)
	err = p.Cluster().PullImage(docker.PullImageOptions{Repository: "tsuru/bs"}, p.RegistryAuthConfig())
	c.Assert(err, check.NotNil)
	e, ok := err.(cluster.DockerNodeError)
	c.Assert(ok, check.Equals, true)
	c.Assert(e.BaseError(), check.ErrorMatches, "cannot connect to Docker endpoint")
}

func (s *S) TestServers(c *check.C) {
	server, err := dtesting.NewServer("127.0.0.1:0", nil, nil)
	c.Assert(err, check.IsNil)
	defer server.Stop()
	var p FakeDockerProvisioner
	p.servers = append(p.servers, server)
	c.Assert(p.Servers(), check.DeepEquals, p.servers)
}

func (s *S) TestCluster(c *check.C) {
	var p FakeDockerProvisioner
	cluster, err := cluster.New(nil, &cluster.MapStorage{})
	c.Assert(err, check.IsNil)
	p.cluster = cluster
	c.Assert(p.Cluster(), check.Equals, cluster)
}

func (s *S) TestPushImage(c *check.C) {

	var p FakeDockerProvisioner
	err := p.PushImage("tsuru/bs", "v1")
	c.Assert(err, check.IsNil)
	expected := []Push{{Name: "tsuru/bs", Tag: "v1"}}
	c.Assert(p.Pushes(), check.DeepEquals, expected)
}

func (s *S) TestPushImageFailure(c *check.C) {
	p := FakeDockerProvisioner{pushErrors: make(chan error, 1)}
	prepErr := errors.New("fail to push")
	p.FailPush(prepErr)
	err := p.PushImage("tsuru/bs", "v1")
	c.Assert(err, check.Equals, prepErr)
	expected := []Push{{Name: "tsuru/bs", Tag: "v1"}}
	c.Assert(p.Pushes(), check.DeepEquals, expected)
}

func (s *S) TestRegistryAuthConfig(c *check.C) {
	var p FakeDockerProvisioner
	config := p.RegistryAuthConfig()
	c.Assert(config, check.Equals, p.authConfig)
}

func (s *S) TestAllContainers(c *check.C) {
	p, err := NewFakeDockerProvisioner()
	c.Assert(err, check.IsNil)
	defer p.Destroy()
	cont1 := container.Container{ID: "cont1"}
	cont2 := container.Container{ID: "cont2"}
	p.SetContainers("localhost", []container.Container{cont1})
	p.SetContainers("remotehost", []container.Container{cont2})
	cont1.HostAddr = "localhost"
	cont2.HostAddr = "remotehost"
	containers := p.AllContainers()
	expected := []container.Container{cont1, cont2}
	if expected[0].HostAddr != containers[0].HostAddr {
		expected = []container.Container{cont2, cont1}
	}
	c.Assert(containers, check.DeepEquals, expected)
}

func (s *S) TestStartContainers(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	p, err := StartMultipleServersCluster()
	c.Assert(err, check.IsNil)
	defer p.Destroy()
	containers, err := p.StartContainers(StartContainersArgs{
		Amount:    map[string]int{"web": 2, "worker": 1},
		Image:     "tsuru/python",
		PullImage: true,
		Endpoint:  p.Servers()[0].URL(),
		App:       app,
	})
	c.Assert(err, check.IsNil)
	c.Assert(containers, check.HasLen, 3)
	c.Assert(p.Containers(urlToHost(p.Servers()[0].URL())), check.DeepEquals, containers)
}
*/
