package user

import (
	"fmt"
	"game/gateway"
	"game/pb"
)

type UserRpc struct {
	gateway.TContextHandler
}

func init() {
	var dummyUser = UserRpc{}
	gateway.RegisterFunc("user", &dummyUser)
}

func (this *UserRpc) Login(arg *pb.TEmptyReq) (ret *pb.TUserLoginResponse) {
	ret = new(pb.TUserLoginResponse)

	context := this.GetContext()

	fmt.Println(context)
	log := context.GetLogger()

	log.Info("Begin Login %s", context.GetUid())

	userObj := GetUserObj(10001, context)
	uname := userObj.GetUname()
	ret.Name = &uname

	log.Info("UserInfo  %s ", userObj.GetAllInfo())
	log.Info("Finish Login %s ", context.GetUid())
	return ret
}
