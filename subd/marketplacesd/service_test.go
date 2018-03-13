package marketplacesd

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct {
	service *Service
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	srv := NewService(nil, &Config{}, nil)
	s.service = srv
	c.Assert(srv, check.NotNil)
}
