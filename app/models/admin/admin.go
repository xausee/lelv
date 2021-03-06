package admin

import (
	"lelv/app/models/dbmgr"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// HomeBlogID 首页各个板块的博客id号集合
type HomeBlogID struct {
	HeadBlogs      []string // 今日博客头条博客ID号集合
	HeadFamousBlog []string // 今日推荐名博ID
	NaturalBlogs   []string // 纯美大自然博客ID号集合
	HistoryBlogs   []string // 历史.有韵味博客ID号集合
	CustomsBlogs   []string // 世界.奇异风情博客ID号集合
	ShareBlogs     []string // 达人分享博客ID号集合
	LatestBlogs    []string // 及时更新博客ID号集合

	PicturesBlogs []string // 图片博客ID号集合
	ArticlesBlogs []string // 文字博客ID号集合
	TimeStamp     string   // 时间戳
}

// BlogType 首页博客ID所属模块类型
type BlogType int

// 博客ID的分类模块类型
const (
	HeadBlogs BlogType = iota
	HeadFamousBlog
	NaturalBlogs
	HistoryBlogs
	CustomsBlogs
	ShareBlogs
	LatestBlogs

	PicturesBlogs
	ArticlesBlogs
)

// AddOrUpdate 新增或者更新
func (b *HomeBlogID) AddOrUpdate(t BlogType, ids []string) error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.HomeBlogIDs)

	if HasTodayRecord() {
		old, err := FindByTimeStamp(time.Now().Format("2006-01-02"))
		if err != nil {
			return err
		}

		new := old

		switch t {
		case HeadBlogs:
			new.HeadBlogs = ids
		case HeadFamousBlog:
			new.HeadFamousBlog = ids
		case NaturalBlogs:
			new.NaturalBlogs = ids
		case HistoryBlogs:
			new.HistoryBlogs = ids
		case CustomsBlogs:
			new.CustomsBlogs = ids
		case ShareBlogs:
			new.ShareBlogs = ids
		case LatestBlogs:
			new.LatestBlogs = ids
		case PicturesBlogs:
			new.PicturesBlogs = ids
		case ArticlesBlogs:
			new.ArticlesBlogs = ids
		}

		err = c.Update(old, new)
		if err != nil {
			return err
		}

		return nil
	}

	var new HomeBlogID
	switch t {
	case HeadBlogs:
		new.HeadBlogs = ids
	case HeadFamousBlog:
		new.HeadFamousBlog = ids
	case NaturalBlogs:
		new.NaturalBlogs = ids
	case HistoryBlogs:
		new.HistoryBlogs = ids
	case CustomsBlogs:
		new.CustomsBlogs = ids
	case ShareBlogs:
		new.ShareBlogs = ids
	case LatestBlogs:
		new.LatestBlogs = ids
	}

	new.TimeStamp = time.Now().Format("2006-01-02")
	err = c.Insert(new)
	if err != nil {
		return err
	}

	return nil
}

// FindByTimeStamp 根据时间戳查询
func FindByTimeStamp(t string) (HomeBlogID, error) {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.HomeBlogIDs)

	var hbi HomeBlogID
	err = c.Find(bson.M{"timestamp": t}).One(&hbi)

	return hbi, err
}

// GetLast 根据时间戳查询最新的记录
func GetLast() (HomeBlogID, error) {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.HomeBlogIDs)

	var ids HomeBlogID
	err = c.Find(nil).Sort("-timestamp").Limit(1).One(&ids)

	return ids, err
}

// HasTodayRecord 是否有当天的编辑记录
func HasTodayRecord() bool {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.HomeBlogIDs)

	flag := false
	n, err := c.Find(bson.M{"timestamp": time.Now().Format("2006-01-02")}).Count()
	if err != nil {
		log.Println("查询当天编辑记录失败")
		flag = false
	}

	if n == 1 {
		flag = true
	}

	return flag
}
