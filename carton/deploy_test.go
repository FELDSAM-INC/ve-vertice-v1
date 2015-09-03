/*
** copyright [2013-2015] [Megam Systems]
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
package carton

import (
	"bytes"
	"errors"
	"io/ioutil"
	"time"

	"github.com/megamsys/megamd/provision/provisiontest"
	"github.com/megamsys/megamd/repository"
	"github.com/megamsys/megamd/repository/repositorytest"
	"gopkg.in/check.v1"
)

func (s *S) TestListDeployByNonAdminUsers(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer nativeScheme.Remove(user)
	team := &auth.Team{Name: "someteam", Users: []string{user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, check.IsNil)
	defer s.conn.Teams().RemoveId("someteam")
	s.conn.Deploys().RemoveAll(nil)
	a := App{Name: "g1", Teams: []string{team.Name}}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	a2 := App{Name: "ge"}
	err = s.conn.Apps().Insert(a2)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	defer s.conn.Apps().Remove(bson.M{"name": a2.Name})
	deploys := []DeployData{
		{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second)},
		{App: "ge", Timestamp: time.Now()},
	}
	for _, deploy := range deploys {
		s.conn.Deploys().Insert(deploy)
	}
	defer s.conn.Deploys().RemoveAll(bson.M{"app": a.Name})
	result, err := ListDeploys(nil, nil, user, 0, 0)
	c.Assert(err, check.IsNil)
	c.Assert(result, check.HasLen, 1)
	c.Assert(result[0].App, check.Equals, "g1")
}

func (s *S) TestListDeployByAdminUsers(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer nativeScheme.Remove(user)
	team := &auth.Team{Name: "adminteam", Users: []string{user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, check.IsNil)
	defer s.conn.Teams().RemoveId("adminteam")
	s.conn.Deploys().RemoveAll(nil)
	adminTeamName, err := config.GetString("admin-team")
	c.Assert(err, check.IsNil)
	config.Set("admin-team", "adminteam")
	defer config.Set("admin-team", adminTeamName)
	a := App{Name: "g1", Teams: []string{team.Name}}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	a2 := App{Name: "ge"}
	err = s.conn.Apps().Insert(a2)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	defer s.conn.Apps().Remove(bson.M{"name": a2.Name})
	deploys := []DeployData{
		{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second)},
		{App: "ge", Timestamp: time.Now()},
	}
	for _, deploy := range deploys {
		s.conn.Deploys().Insert(deploy)
	}
	defer s.conn.Deploys().RemoveAll(bson.M{"app": a.Name})
	result, err := ListDeploys(nil, nil, user, 0, 0)
	c.Assert(err, check.IsNil)
	c.Assert(result, check.HasLen, 2)
	c.Assert(result[0].App, check.Equals, "ge")
	c.Assert(result[1].App, check.Equals, "g1")
}

func (s *S) TestListDeployByAppAndService(c *check.C) {
	s.conn.Deploys().RemoveAll(nil)
	srv := service.Service{Name: "mysql"}
	instance := service.ServiceInstance{
		Name:        "myinstance",
		ServiceName: "mysql",
		Apps:        []string{"g1"},
	}
	err := s.conn.ServiceInstances().Insert(instance)
	err = s.conn.Services().Insert(srv)
	c.Assert(err, check.IsNil)
	defer s.conn.ServiceInstances().Remove(bson.M{"apps": instance.Apps})
	defer s.conn.Services().Remove(bson.M{"_id": srv.Name})
	a := App{Name: "g1"}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	a2 := App{Name: "ge"}
	err = s.conn.Apps().Insert(a2)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	defer s.conn.Apps().Remove(bson.M{"name": a2.Name})
	deploys := []DeployData{
		{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second)},
		{App: "ge", Timestamp: time.Now()},
	}
	for _, deploy := range deploys {
		s.conn.Deploys().Insert(deploy)
	}
	defer s.conn.Deploys().RemoveAll(bson.M{"app": a.Name})
	result, err := ListDeploys(&a2, &srv, nil, 0, 0)
	c.Assert(err, check.IsNil)
	c.Assert(result, check.IsNil)
}

func (s *S) TestListAppDeploys(c *check.C) {
	s.conn.Deploys().RemoveAll(nil)
	a := App{Name: "g1"}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	insert := []interface{}{
		DeployData{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second)},
		DeployData{App: "g1", Timestamp: time.Now()},
	}
	s.conn.Deploys().Insert(insert...)
	defer s.conn.Deploys().RemoveAll(bson.M{"app": a.Name})
	expected := []DeployData{insert[1].(DeployData), insert[0].(DeployData)}
	deploys, err := a.ListDeploys(nil)
	c.Assert(err, check.IsNil)
	for i := 0; i < 2; i++ {
		ts := expected[i].Timestamp
		expected[i].Timestamp = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
		ts = deploys[i].Timestamp
		deploys[i].Timestamp = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
		expected[i].ID = deploys[i].ID
	}
	c.Assert(deploys, check.DeepEquals, expected)
}

func (s *S) TestListServiceDeploys(c *check.C) {
	s.conn.Deploys().RemoveAll(nil)
	srv := service.Service{Name: "mysql"}
	instance := service.ServiceInstance{
		Name:        "myinstance",
		ServiceName: "mysql",
		Apps:        []string{"g1"},
	}
	err := s.conn.ServiceInstances().Insert(instance)
	err = s.conn.Services().Insert(srv)
	c.Assert(err, check.IsNil)
	defer s.conn.ServiceInstances().Remove(bson.M{"apps": instance.Apps})
	defer s.conn.Services().Remove(bson.M{"_id": srv.Name})
	insert := []interface{}{
		DeployData{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second)},
		DeployData{App: "g1", Timestamp: time.Now()},
	}
	s.conn.Deploys().Insert(insert...)
	defer s.conn.Deploys().RemoveAll(bson.M{"apps": instance.Apps})
	expected := []DeployData{insert[1].(DeployData), insert[0].(DeployData)}
	deploys, err := ListDeploys(nil, &srv, nil, 0, 0)
	c.Assert(err, check.IsNil)
	for i := 0; i < 2; i++ {
		ts := expected[i].Timestamp
		expected[i].Timestamp = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
		ts = deploys[i].Timestamp
		deploys[i].Timestamp = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
		expected[i].ID = deploys[i].ID
	}
	c.Assert(deploys, check.DeepEquals, expected)
}

func (s *S) TestListAllDeploys(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer user.Delete()
	team := &auth.Team{Name: "team", Users: []string{user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, check.IsNil)
	defer s.conn.Teams().Remove(team)
	a := App{
		Name:     "g1",
		Platform: "zend",
		Teams:    []string{team.Name},
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	a = App{
		Name:     "ge",
		Platform: "zend",
		Teams:    []string{team.Name},
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.conn.Deploys().RemoveAll(nil)
	insert := []interface{}{
		DeployData{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second)},
		DeployData{App: "ge", Timestamp: time.Now(), Image: "app-image"},
	}
	s.conn.Deploys().Insert(insert...)
	defer s.conn.Deploys().RemoveAll(nil)
	expected := []DeployData{insert[1].(DeployData), insert[0].(DeployData)}
	expected[0].CanRollback = true
	deploys, err := ListDeploys(nil, nil, user, 0, 0)
	c.Assert(err, check.IsNil)
	for i := 0; i < 2; i++ {
		ts := expected[i].Timestamp
		expected[i].Timestamp = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
		ts = deploys[i].Timestamp
		deploys[i].Timestamp = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
		expected[i].ID = deploys[i].ID
	}
	c.Assert(deploys, check.DeepEquals, expected)
}

func (s *S) TestListAllDeploysSkipAndLimit(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer user.Delete()
	team := &auth.Team{Name: "team", Users: []string{user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, check.IsNil)
	defer s.conn.Teams().Remove(team)
	a := App{
		Name:     "app1",
		Platform: "zend",
		Teams:    []string{team.Name},
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.conn.Deploys().RemoveAll(nil)
	insert := []interface{}{
		DeployData{App: "app1", Commit: "v1", Timestamp: time.Now().Add(-30 * time.Second)},
		DeployData{App: "app1", Commit: "v2", Timestamp: time.Now().Add(-20 * time.Second)},
		DeployData{App: "app1", Commit: "v3", Timestamp: time.Now().Add(-10 * time.Second)},
		DeployData{App: "app1", Commit: "v4", Timestamp: time.Now()},
	}
	s.conn.Deploys().Insert(insert...)
	defer s.conn.Deploys().RemoveAll(nil)
	expected := []DeployData{insert[2].(DeployData), insert[1].(DeployData)}
	deploys, err := ListDeploys(nil, nil, user, 1, 2)
	c.Assert(err, check.IsNil)
	c.Assert(deploys, check.HasLen, 2)
	for i := 0; i < len(deploys); i++ {
		ts := expected[i].Timestamp
		newTs := time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
		expected[i].Timestamp = newTs
		ts = deploys[i].Timestamp
		deploys[i].Timestamp = newTs
		expected[i].ID = deploys[i].ID
	}
	c.Assert(deploys, check.DeepEquals, expected)
}

func (s *S) TestListDeployByAppAndUser(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer user.Delete()
	team := &auth.Team{Name: "team", Users: []string{user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, check.IsNil)
	defer s.conn.Teams().Remove(team)
	a := App{
		Name:     "g1",
		Platform: "zend",
		Teams:    []string{team.Name},
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	a = App{
		Name:     "ge",
		Platform: "zend",
		Teams:    []string{team.Name},
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.conn.Deploys().RemoveAll(nil)
	insert := []interface{}{
		DeployData{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second)},
		DeployData{App: "ge", Timestamp: time.Now()},
	}
	s.conn.Deploys().Insert(insert...)
	defer s.conn.Deploys().RemoveAll(nil)
	expected := []DeployData{insert[1].(DeployData)}
	deploys, err := ListDeploys(&a, nil, user, 0, 0)
	c.Assert(err, check.IsNil)
	c.Assert(expected[0].App, check.DeepEquals, deploys[0].App)
	c.Assert(len(expected), check.Equals, len(deploys))
}

func (s *S) TestGetDeploy(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer user.Delete()
	team := &auth.Team{Name: "team", Users: []string{user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, check.IsNil)
	defer s.conn.Teams().Remove(team)
	a := App{
		Name:     "g1",
		Platform: "zend",
		Teams:    []string{team.Name},
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.conn.Deploys().RemoveAll(nil)
	newDeploy := DeployData{ID: bson.NewObjectId(), App: "g1", Timestamp: time.Now()}
	err = s.conn.Deploys().Insert(&newDeploy)
	c.Assert(err, check.IsNil)
	defer s.conn.Deploys().Remove(bson.M{"name": newDeploy.App})
	lastDeploy, err := GetDeploy(newDeploy.ID.Hex(), user)
	c.Assert(err, check.IsNil)
	ts := lastDeploy.Timestamp
	lastDeploy.Timestamp = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
	ts = newDeploy.Timestamp
	newDeploy.Timestamp = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC)
	c.Assert(lastDeploy.ID, check.Equals, newDeploy.ID)
	c.Assert(lastDeploy.App, check.Equals, newDeploy.App)
	c.Assert(lastDeploy.Timestamp, check.Equals, newDeploy.Timestamp)
}

func (s *S) TestGetDeployWithoutAccess(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer user.Delete()
	a := App{
		Name:     "g1",
		Platform: "zend",
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.conn.Deploys().RemoveAll(nil)
	newDeploy := DeployData{ID: bson.NewObjectId(), App: "g1", Timestamp: time.Now()}
	err = s.conn.Deploys().Insert(&newDeploy)
	c.Assert(err, check.IsNil)
	defer s.conn.Deploys().Remove(bson.M{"name": newDeploy.App})
	result, err := GetDeploy(newDeploy.ID.Hex(), user)
	c.Assert(err.Error(), check.Equals, "Deploy not found.")
	c.Assert(result, check.IsNil)
}

func (s *S) TestGetDeployNotFound(c *check.C) {
	idTest := bson.NewObjectId()
	deploy, err := GetDeploy(idTest.Hex(), nil)
	c.Assert(err.Error(), check.Equals, "not found")
	c.Assert(deploy, check.IsNil)
}

func (s *S) TestGetDiffInDeploys(c *check.C) {
	s.conn.Deploys().RemoveAll(nil)
	myDeploy := DeployData{
		App:       "g1",
		Timestamp: time.Now().Add(-3600 * time.Second),
		Commit:    "545b1904af34458704e2aa06ff1aaffad5289f8g",
		Origin:    "git",
	}
	deploys := []DeployData{
		{App: "ge", Timestamp: time.Now(), Commit: "hwed834hf8y34h8fhn8rnr823nr238runh23x", Origin: "git"},
		{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second * 2), Commit: "545b1904af34458704e2aa06ff1aaffad5289f8f", Origin: "git"},
		myDeploy,
		{App: "g1", Timestamp: time.Now(), Commit: "1b970b076bbb30d708e262b402d4e31910e1dc10", Origin: "git"},
	}
	for _, d := range deploys {
		s.conn.Deploys().Insert(d)
	}
	defer s.conn.Deploys().RemoveAll(nil)
	err := s.conn.Deploys().Find(bson.M{"commit": myDeploy.Commit}).One(&myDeploy)
	c.Assert(err, check.IsNil)
	repository.Manager().CreateRepository("g1", nil)
	diffOutput, err := GetDiffInDeploys(&myDeploy)
	c.Assert(err, check.IsNil)
	c.Assert(diffOutput, check.Equals, repositorytest.Diff)
}

func (s *S) TestGetDiffInDeploysWithOneCommit(c *check.C) {
	s.conn.Deploys().RemoveAll(nil)
	lastDeploy := DeployData{App: "g1", Timestamp: time.Now(), Commit: "1b970b076bbb30d708e262b402d4e31910e1dc10"}
	s.conn.Deploys().Insert(lastDeploy)
	defer s.conn.Deploys().RemoveAll(nil)
	err := s.conn.Deploys().Find(bson.M{"commit": lastDeploy.Commit}).One(&lastDeploy)
	c.Assert(err, check.IsNil)
	diffOutput, err := GetDiffInDeploys(&lastDeploy)
	c.Assert(err, check.IsNil)
	c.Assert(diffOutput, check.Equals, "The deployment must have at least two commits for the diff.")
}

func (s *S) TestGetDiffInDeploysNoGit(c *check.C) {
	s.conn.Deploys().RemoveAll(nil)
	myDeploy := DeployData{
		App:       "g1",
		Timestamp: time.Now().Add(-3600 * time.Second),
		Commit:    "545b1904af34458704e2aa06ff1aaffad5289f8g",
		Origin:    "app-deploy",
	}
	deploys := []DeployData{
		{App: "ge", Timestamp: time.Now(), Commit: "hwed834hf8y34h8fhn8rnr823nr238runh23x", Origin: "git"},
		{App: "g1", Timestamp: time.Now().Add(-3600 * time.Second * 2), Commit: "545b1904af34458704e2aa06ff1aaffad5289f8f", Origin: "git"},
		myDeploy,
		{App: "g1", Timestamp: time.Now(), Commit: "1b970b076bbb30d708e262b402d4e31910e1dc10", Origin: "git"},
	}
	for _, d := range deploys {
		s.conn.Deploys().Insert(d)
	}
	defer s.conn.Deploys().RemoveAll(nil)
	err := s.conn.Deploys().Find(bson.M{"commit": myDeploy.Commit}).One(&myDeploy)
	c.Assert(err, check.IsNil)
	repository.Manager().CreateRepository("g1", nil)
	diffOutput, err := GetDiffInDeploys(&myDeploy)
	c.Assert(err, check.IsNil)
	c.Assert(diffOutput, check.Equals, "Cannot have diffs between git based and app-deploy based deployments")
}

func (s *S) TestDeployApp(c *check.C) {
	a := App{
		Name:     "someApp",
		Platform: "django",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	err = Deploy("github.com/megamsys{
		App:          &a,
		Version:      "version",
		Commit:       "1ee1f1084927b3a5db59c9033bc5c4abefb7b93c",
		OutputStream: writer,
	})
	c.Assert(err, check.IsNil)
	logs := writer.String()
	c.Assert(logs, check.Equals, "Git deploy called")
}

func (s *S) TestDeployAppWithUpdatePlatform(c *check.C) {
	a := App{
		Name:           "someApp",
		Platform:       "django",
		Teams:          []string{s.team.Name},
		UpdatePlatform: true,
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	err = Deploy("github.com/megamsys{
		App:          &a,
		Version:      "version",
		Commit:       "1ee1f1084927b3a5db59c9033bc5c4abefb7b93c",
		OutputStream: writer,
	})
	c.Assert(err, check.IsNil)
	logs := writer.String()
	c.Assert(logs, check.Equals, "Git deploy called")
	var updatedApp App
	s.conn.Apps().Find(bson.M{"name": "someApp"}).One(&updatedApp)
	c.Assert(updatedApp.UpdatePlatform, check.Equals, false)
}

func (s *S) TestDeployAppIncrementDeployNumber(c *check.C) {
	a := App{
		Name:     "otherapp",
		Platform: "zend",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	err = Deploy("github.com/megamsys{
		App:          &a,
		Version:      "version",
		Commit:       "1ee1f1084927b3a5db59c9033bc5c4abefb7b93c",
		OutputStream: writer,
	})
	c.Assert(err, check.IsNil)
	s.conn.Apps().Find(bson.M{"name": a.Name}).One(&a)
	c.Assert(a.Deploys, check.Equals, uint(1))
}

func (s *S) TestDeployAppSaveDeployData(c *check.C) {
	a := App{
		Name:     "otherapp",
		Platform: "zend",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	commit := "1ee1f1084927b3a5db59c9033bc5c4abefb7b93c"
	err = Deploy("github.com/megamsys{
		App:          &a,
		Version:      "version",
		Commit:       commit,
		OutputStream: writer,
		User:         "someone@themoon",
	})
	c.Assert(err, check.IsNil)
	s.conn.Apps().Find(bson.M{"name": a.Name}).One(&a)
	c.Assert(a.Deploys, check.Equals, uint(1))
	var result map[string]interface{}
	s.conn.Deploys().Find(bson.M{"app": a.Name}).One(&result)
	c.Assert(result["app"], check.Equals, a.Name)
	now := time.Now()
	diff := now.Sub(result["timestamp"].(time.Time))
	c.Assert(diff < 60*time.Second, check.Equals, true)
	c.Assert(result["duration"], check.Not(check.Equals), 0)
	c.Assert(result["commit"], check.Equals, commit)
	c.Assert(result["image"], check.Equals, "app-image")
	c.Assert(result["log"], check.Equals, "Git deploy called")
	c.Assert(result["user"], check.Equals, "someone@themoon")
	c.Assert(result["origin"], check.Equals, "git")
}

func (s *S) TestDeployAppSaveDeployDataOriginRollback(c *check.C) {
	a := App{
		Name:     "otherapp",
		Platform: "zend",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	err = Deploy("github.com/megamsys{
		App:          &a,
		OutputStream: writer,
		Image:        "some-image",
	})
	c.Assert(err, check.IsNil)
	s.conn.Apps().Find(bson.M{"name": a.Name}).One(&a)
	c.Assert(a.Deploys, check.Equals, uint(1))
	var result map[string]interface{}
	s.conn.Deploys().Find(bson.M{"app": a.Name}).One(&result)
	c.Assert(result["app"], check.Equals, a.Name)
	now := time.Now()
	diff := now.Sub(result["timestamp"].(time.Time))
	c.Assert(diff < 60*time.Second, check.Equals, true)
	c.Assert(result["duration"], check.Not(check.Equals), 0)
	c.Assert(result["image"], check.Equals, "some-image")
	c.Assert(result["log"], check.Equals, "Image deploy called")
	c.Assert(result["origin"], check.Equals, "rollback")
}

func (s *S) TestDeployAppSaveDeployDataOriginAppDeploy(c *check.C) {
	a := App{
		Name:     "otherapp",
		Platform: "zend",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	err = Deploy("github.com/megamsys{
		App:          &a,
		OutputStream: writer,
		File:         ioutil.NopCloser(bytes.NewBuffer([]byte("my file"))),
	})
	c.Assert(err, check.IsNil)
	s.conn.Apps().Find(bson.M{"name": a.Name}).One(&a)
	c.Assert(a.Deploys, check.Equals, uint(1))
	var result map[string]interface{}
	s.conn.Deploys().Find(bson.M{"app": a.Name}).One(&result)
	c.Assert(result["app"], check.Equals, a.Name)
	now := time.Now()
	diff := now.Sub(result["timestamp"].(time.Time))
	c.Assert(diff < 60*time.Second, check.Equals, true)
	c.Assert(result["duration"], check.Not(check.Equals), 0)
	c.Assert(result["image"], check.Equals, "app-image")
	c.Assert(result["log"], check.Equals, "Upload deploy called")
	c.Assert(result["origin"], check.Equals, "app-deploy")
}

func (s *S) TestDeployAppSaveDeployErrorData(c *check.C) {
	provisioner := provisiontest.NewFakeProvisioner()
	provisioner.PrepareFailure("GitDeploy", errors.New("deploy error"))
	Provisioner = provisioner
	defer func() {
		Provisioner = s.provisioner
	}()
	a := App{
		Name:     "testErrorApp",
		Platform: "zend",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	provisioner.Provision(&a)
	defer provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	err = Deploy("github.com/megamsys{
		App:          &a,
		Version:      "version",
		Commit:       "1ee1f1084927b3a5db59c9033bc5c4abefb7b93c",
		OutputStream: writer,
	})
	c.Assert(err, check.NotNil)
	var result map[string]interface{}
	s.conn.Deploys().Find(bson.M{"app": a.Name}).One(&result)
	c.Assert(result["app"], check.Equals, a.Name)
	c.Assert(result["error"], check.NotNil)
}

func (s *S) TestUserHasPermission(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer user.Delete()
	team := &auth.Team{Name: "team", Users: []string{user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, check.IsNil)
	defer s.conn.Teams().Remove(team)
	a := App{
		Name:     "g1",
		Platform: "zend",
		Teams:    []string{team.Name},
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	hasPermission := userHasPermission(user, a.Name)
	c.Assert(hasPermission, check.Equals, true)
}

func (s *S) TestUserHasNoPermission(c *check.C) {
	user := &auth.User{Email: "user@user.com", Password: "123456"}
	nativeScheme := auth.ManagedScheme(native.NativeScheme{})
	AuthScheme = nativeScheme
	_, err := nativeScheme.Create(user)
	c.Assert(err, check.IsNil)
	defer user.Delete()
	team := &auth.Team{Name: "team", Users: []string{user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, check.IsNil)
	defer s.conn.Teams().Remove(team)
	a := App{
		Name:     "g1",
		Platform: "zend",
	}
	err = s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	hasPermission := userHasPermission(user, a.Name)
	c.Assert(hasPermission, check.Equals, false)
}

func (s *S) TestIncrementDeploy(c *check.C) {
	a := App{
		Name:     "otherapp",
		Platform: "zend",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	incrementDeploy(&a)
	s.conn.Apps().Find(bson.M{"name": a.Name}).One(&a)
	c.Assert(a.Deploys, check.Equals, uint(1))
}

func (s *S) TestDeployToProvisioner(c *check.C) {
	a := App{
		Name:     "someApp",
		Platform: "django",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	opts := "github.com/megamsys{App: &a, Version: "version"}
	_, err = deployToProvisioner(&opts, writer)
	c.Assert(err, check.IsNil)
	logs := writer.String()
	c.Assert(logs, check.Equals, "Git deploy called")
}

func (s *S) TestDeployToProvisionerArchive(c *check.C) {
	a := App{
		Name:     "someApp",
		Platform: "django",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	opts := "github.com/megamsys{App: &a, ArchiveURL: "https://s3.amazonaws.com/smt/archive.tar.gz"}
	_, err = deployToProvisioner(&opts, writer)
	c.Assert(err, check.IsNil)
	logs := writer.String()
	c.Assert(logs, check.Equals, "Archive deploy called")
}

func (s *S) TestDeployToProvisionerUpload(c *check.C) {
	a := App{
		Name:     "someApp",
		Platform: "django",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	opts := "github.com/megamsys{App: &a, File: ioutil.NopCloser(bytes.NewBuffer([]byte("my file")))}
	_, err = deployToProvisioner(&opts, writer)
	c.Assert(err, check.IsNil)
	logs := writer.String()
	c.Assert(logs, check.Equals, "Upload deploy called")
}

func (s *S) TestDeployToProvisionerImage(c *check.C) {
	a := App{
		Name:     "someApp",
		Platform: "django",
		Teams:    []string{s.team.Name},
	}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	s.provisioner.Provision(&a)
	defer s.provisioner.Destroy(&a)
	writer := &bytes.Buffer{}
	opts := "github.com/megamsys{App: &a, Image: "my-image-x"}
	_, err = deployToProvisioner(&opts, writer)
	c.Assert(err, check.IsNil)
	logs := writer.String()
	c.Assert(logs, check.Equals, "Image deploy called")
}

func (s *S) TestMarkDeploysAsRemoved(c *check.C) {
	s.createAdminUserAndTeam(c)
	a := App{Name: "someApp"}
	err := s.conn.Apps().Insert(a)
	c.Assert(err, check.IsNil)
	defer s.conn.Apps().Remove(bson.M{"name": a.Name})
	opts := "github.com/megamsys{
		App:     &a,
		Version: "version",
		Commit:  "1ee1f1084927b3a5db59c9033bc5c4abefb7b93c",
	}
	err = saveDeployData(&opts, "myid", "mylog", time.Second, nil)
	c.Assert(err, check.IsNil)
	c.Assert(err, check.IsNil)
	defer s.conn.Deploys().RemoveAll(bson.M{"app": a.Name})
	result, err := ListDeploys(nil, nil, s.admin, 0, 0)
	c.Assert(err, check.IsNil)
	c.Assert(result, check.HasLen, 1)
	c.Assert(result[0].Image, check.Equals, "myid")
	err = markDeploysAsRemoved(a.Name)
	c.Assert(err, check.IsNil)
	result, err = ListDeploys(nil, nil, s.admin, 0, 0)
	c.Assert(err, check.IsNil)
	c.Assert(result, check.HasLen, 0)
	var allDeploys []DeployData
	err = s.conn.Deploys().Find(nil).All(&allDeploys)
	c.Assert(err, check.IsNil)
	c.Assert(allDeploys, check.HasLen, 1)
	c.Assert(allDeploys[0].Image, check.Equals, "myid")
	c.Assert(allDeploys[0].RemoveDate.IsZero(), check.Equals, false)
}
