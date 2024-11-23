package tests

import (
	"sso/tests/suite"
	"testing"
)

const (
	emptyAppID     = 0
	appID          = 1
	appSecret      = "test-secret"
	pathDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
}
