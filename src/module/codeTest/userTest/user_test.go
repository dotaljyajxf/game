package userTest

import (
	"module/codeTest"
	"pb"
	"testing"
)

func TestUserGet(t *testing.T) {
	player := &codeTest.UserClient{}

	arg := pb.TEmptyReq{}
	resp := pb.TUserLoginResponse{}

	player.SendReq("user.Login", &arg, &resp)

	t.Log(resp)
}
