package data

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"os"
)

var mData *xorm.Engine

func NewData(dbSrc string, logPath string) *xorm.Engine {
	if mData == nil {
		var err error
		mData, err = xorm.NewEngine("mysql", dbSrc)
		if err != nil {
			panic("NewData NewEngine err")
		}
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic("create logFile err")
		}
		mData.Logger().SetLevel(core.LOG_DEBUG)
		mData.SetLogger(xorm.NewSimpleLogger(f))

		//采用了LRU算法的一个缓存，缓存方式是存放到内存中，缓存struct的记录数为200条
		cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 200)
		mData.SetDefaultCacher(cacher)
		//mData.ClearCacheBean()

		return mData
	} else {
		return mData
	}
}
