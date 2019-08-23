package data

import (
	"data/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"os"
)

var mData *xorm.Engine

//dbSrc "root:123@/test?charset=utf8"
//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
//user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true

func InitDb(dbSrc string, logPath string) error {

	var err error
	mData, err = xorm.NewEngine("mysql", dbSrc)
	if err != nil {
		panic("NewData NewEngine err")
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("create logFile err")
	}
	//TODO
	mData.Logger().SetLevel(1)
	mData.SetLogger(xorm.NewSimpleLogger(f))

	//采用了LRU算法的一个缓存，缓存方式是存放到内存中，缓存struct的记录数为500条
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 500)
	mData.SetDefaultCacher(cacher)
	//mData.ClearCacheBean()

	err = mData.Sync2(db.DbMap)

	return err

}

func NewEngine() *xorm.Engine {
	if mData == nil {
		panic("db need init")
	}
	return mData
}
