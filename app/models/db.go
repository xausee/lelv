package models

import (
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
)

// 数据库信息
const (
	Name          = "travelblog"    // 数据库名
	PV            = "pv"            // page view数据表
	Users         = "users"         // 用户信息数据表
	Blogs         = "blogs"         // 博客数据表
	Conversations = "conversations" // 用户间的消息数据表
	HomeBlogIDs   = "homeblogids"   // 首页各个板块的博客id号集合
	CarouselBlogs = "carouselblogs" // 首页轮播图片数据表
	Guest         = "Guest"         // 游客名字常量
)

// DBManager 数据库管理器
type DBManager struct {
	session *mgo.Session
}

// NewDBManager 创建数据库管理器对象
func NewDBManager() (*DBManager, error) {
	revel.Config.SetSection("db")
	ip, found := revel.Config.String("ip")
	if !found {
		revel.ERROR.Fatal("Cannot load database ip from app.conf")
	}

	session, err := mgo.Dial(ip)
	if err != nil {
		return nil, err
	}

	// revel.Config.SetSection("db")
	// uri, found := revel.Config.String("mongodburi")
	// if !found {
	// 	revel.ERROR.Fatal("找不到配置：mongodb uri")
	// }

	// dialinfo, err := mgo.ParseURL(uri)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(dialinfo)

	// session, err := mgo.Dial(uri)
	// if err != nil {
	// 	return nil, err
	// }

	return &DBManager{session}, nil
}

// SetDB 根据数据库名字，创建数据库连接
func (manager *DBManager) SetDB(name string) *mgo.Database {
	return manager.session.DB(name)
}

// Coll 根据数据库表名，返回表对象
func (manager *DBManager) Coll(name string) *mgo.Collection {
	return manager.session.DB(Name).C(name)
}

// Close 关闭数据库
func (manager *DBManager) Close() {
	manager.session.Close()
}
