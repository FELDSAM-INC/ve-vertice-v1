/*
** Copyright [2013-2017] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package one

/**
import(
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/virtengine/libgo/hc"
	"gopkg.in/check.v1"
)
*/

/*

func (s *S) TestHealthCheckDocker(c *check.C) {
	var request *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("OK"))
	}))
	defer server.Close()
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(nil, &cluster.MapStorage{}, cluster.Node{Address: server.URL})
	c.Assert(err, check.IsNil)
	err = healthCheckDocker()
	c.Assert(err, check.IsNil)
	c.Assert(request.Method, check.Equals, "GET")
	c.Assert(request.URL.Path, check.Equals, "/_ping")
}

func (s *S) TestHealthCheckDockerMultipleNodes(c *check.C) {
	var request *http.Request
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("OK"))
	}))
	defer server1.Close()
	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("OK"))
	}))
	defer server2.Close()
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(nil, &cluster.MapStorage{},
		cluster.Node{Address: server1.URL}, cluster.Node{Address: server2.URL})
	c.Assert(err, check.IsNil)
	err = healthCheckDocker()
	c.Assert(err, check.Equals, hc.ErrDisabledComponent)
	c.Assert(request, check.IsNil)
}

func (s *S) TestHealthCheckDockerNoNodes(c *check.C) {
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(nil, &cluster.MapStorage{})
	c.Assert(err, check.IsNil)
	err = healthCheckDocker()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "error - no nodes available for running containers")
}

func (s *S) TestHealthCheckDockerFailure(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
	}))
	defer server.Close()
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(nil, &cluster.MapStorage{}, cluster.Node{Address: server.URL})
	c.Assert(err, check.IsNil)
	err = healthCheckDocker()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "ping failed - API error (500): something went wrong")
}
*/
