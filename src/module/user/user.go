package user

import (
	"netserver"
)
import "pb"

type UserRpc struct {
	netserver.TContextHandler
}

func init() {
	var dummyUser = UserRpc{}
	netserver.RegisterFunc("user", &dummyUser)
}

func (this *UserRpc) Login(arg *pb.TEmptyReq) (ret *pb.TUserLoginResponse) {
	ret = new(pb.TUserLoginResponse)

	context := this.GetContext()
	log := context.GetLogger()

	log.Info("Begin Login", context.GetUid())

	name := "haha"
	ret.Name = &name

	log.Info("Finish Login", context.GetUid())
	return ret
}
