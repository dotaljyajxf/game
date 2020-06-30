package userTest

import (
	"game/module/codeTest"
	"game/pb"
	"testing"
)

func TestUserGet(t *testing.T) {
	player := &codeTest.UserClient{}

	arg := pb.TEmptyReq{}
	resp := pb.TUserLoginResponse{}

	player.SendReq("user.Login", &arg, &resp)

	t.Log(resp)
}
