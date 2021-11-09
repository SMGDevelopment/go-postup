package postup

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewPostUp(t *testing.T) {
	const (
		user     = "USER"
		password = "PASSWORD"
	)

	var (
		is = is.New(t)
		pu = NewPostUp(user, password, nil)
	)

	is.True(pu.client != nil)
	is.Equal(pu.username, user)
	is.Equal(pu.password, password)
}
