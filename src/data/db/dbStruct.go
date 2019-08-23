package db

type User struct {
	Uid   int64  `xorm:"pk notnull"`
	Uname string `xorm:"index"`
	Pid   string
	Level int64 `xorm:"index"`
}

var DbMap []interface{} = []interface{}{
	&User{},
}
