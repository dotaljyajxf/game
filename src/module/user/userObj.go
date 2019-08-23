package user

import (
	"data/db"
	"netserver"
	"strconv"
)

type UserObj struct {
	userTbl db.User
	context *netserver.TContext
}

func initUser(uid int64, context *netserver.TContext, userObj *UserObj) {
	userObj.userTbl.Uid = uid
	userObj.userTbl.Uname = "game" + strconv.Itoa(int(uid))
	userObj.userTbl.Level = 1
	userObj.userTbl.Pid = strconv.Itoa(int(uid))

	_, err := context.DB().Insert(userObj.userTbl)
	if err != nil {
		context.GetLogger().Fatal("UserInsert error:%s", err)
	}
}

func GetUserObj(uid int64, context *netserver.TContext) *UserObj {

	userObj := new(UserObj)
	log := context.GetLogger()
	log.Info("getUserObj user: %s", userObj)
	has, err := context.DB().ID(uid).Get(&userObj.userTbl)
	if err != nil {
		log.Fatal("getUserObj error uid:%d err:%s", uid, err)
	}
	if !has {
		initUser(uid, context, userObj)
	}

	return userObj
}

func (this *UserObj) GetAllInfo() *db.User {
	return &this.userTbl
}

func (this *UserObj) GetUname() string {
	return this.userTbl.Uname
}

func (this *UserObj) GetUid() int64 {
	return this.userTbl.Uid
}

func (this *UserObj) GetPid() string {
	return this.userTbl.Pid
}

func (this *UserObj) GetLevel() int64 {
	return this.userTbl.Level
}

func (this *UserObj) SetUname(name string) {
	this.userTbl.Uname = name
}

func (this *UserObj) SetUid(uid int64) {
	this.userTbl.Uid = uid
}

func (this *UserObj) SetLevel(level int64) {
	this.userTbl.Level = level
}

func (this *UserObj) Update() {
	_, err := this.context.DB().Update(this)
	if err != nil {
		this.context.GetLogger().Fatal("User Update error :%s", err)
	}
	this.context.GetLogger().Info("User Update Finish :%s", this)
}
