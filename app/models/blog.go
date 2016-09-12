package models

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

// Comment 评论
type Comment struct {
	ID                string // 评论ID号
	CommenterID       string // 评论人ID
	CommenterAvatar   string // 评论人头像地址
	CommenterNickName string // 评论人昵称
	Body              string // 评论内容
	TimeStamp         string // 时间戳
}

// Type 博客类型
type Type int

// 博客类型
const (
	Picture Type = iota
	Text
	Hybrid
)

// Blog 博客结构
type Blog struct {
	ID                  string    // ID号
	AuthorID            string    // 作者ID号
	Author              string    // 作者名称
	Tags                []string  // 标签
	Title               string    // 标题
	Type                Type      // 类型
	Cover               string    // 博客封面图片地址
	BriefText           string    // 截取的文本内容
	Content             string    // 内容
	ViewCount           int       // 阅读次数
	PraiseCount         int       // 点赞次数
	Comments            []Comment // 评论集合
	CreateTimeStamp     string    // 创建时间戳
	LastUpdateTimeStamp string    // 最后更新时间戳
}

// Add 新增
func (b *Blog) Add() error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)

	err = c.Insert(b)
	if err != nil {
		log.Println("创建博客失败：")
		log.Println(b)
		return err
	}

	return nil
}

// FindByID 根据博客ID查找
func (b *Blog) FindByID(id string) (blo Blog, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)

	err = c.Find(bson.M{"id": id}).One(&blo)
	if err != nil {
		return blo, err
	}

	return blo, nil
}

// FindByAuthorID 根据作者ID查找
func (b *Blog) FindByAuthorID(id string) (blo []Blog, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)

	err = c.Find(bson.M{"authorid": id}).All(&blo)
	if err != nil {
		return nil, err
	}

	return blo, nil
}

// GetCountByAuthorID 根据作者ID查找
func (b *Blog) GetCountByAuthorID(id string) (count int, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)

	count, err = c.Find(bson.M{"authorid": id}).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByTag 根据博客标签查找
func (b *Blog) FindByTag(tag string) (r []Blog, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)
	type Items map[string]string
	//err = c.Find(bson.M{"tags": Items{"$elemMatch": tag}}).All(&r)
	err = c.Find(bson.M{"tags": tag}).Limit(100).All(&r)

	return r, nil
}

// FindLast 查找最新的n个记录
// func (b *Blog) FindLast(n int) (r []Blog, err error) {
// 	db, err := NewDBManager()
// 	defer db.Close()

// 	c := db.session.DB(Name).C(Blogs)
// 	type Items map[string]string
// 	err = c.Find(nil).Sort("_id").Limit(n).All(&r)
// 	//err = c.Find(bson.M{"$orderby": {"$natural": -1}}).Sort("_id").Limit(n).All(&r)

// 	return r, nil
// }

// FindLast 查找最新的n个记录
func (b *Blog) FindLast(n int) (r []Blog, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)
	type Items map[string]string
	err = c.Find(nil).Sort("timestamp").Limit(n).All(&r)
	//err = c.Find(bson.M{"$orderby": {"$natural": -1}}).Sort("_id").Limit(n).All(&r)

	return r, nil
}

// Count 获取所有博客数量
func (b *Blog) Count() (int, error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)
	n, err := c.Find(bson.M{}).Count()

	if err != nil {
		return 0, err
	}

	return n, err
}

// UpdateView 更新阅读数量
func (b *Blog) UpdateView() {
	db, _ := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)

	var old Blog
	c.Find(bson.M{"id": b.ID}).One(&old)

	n := old
	n.ViewCount++

	c.Update(old, n)
}

// AddComment 更新阅读数量
func (b *Blog) AddComment(comment Comment) error {
	db, _ := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)

	var old Blog
	c.Find(bson.M{"id": b.ID}).One(&old)

	n := old
	n.Comments = append(n.Comments, comment)

	return c.Update(old, n)
}

// GetViewCount 更新阅读数量
func (b *Blog) GetViewCount() int {
	blog, _ := b.FindByID(b.ID)
	return blog.ViewCount
}

// FindAndSortBy 查找所有记录， 按条件排序，取排序后的最新的n个记录
func (b *Blog) FindAndSortBy(field string, n int) (r []Blog, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)
	type Items map[string]string
	err = c.Find(nil).Sort(field).Limit(n).All(&r)

	return r, nil
}

// FindByAndSortBy 查找所有记录， 按条件排序，取排序后的最新的n个记录
func (b *Blog) FindByAndSortBy(t Type, field string, n int) (r []Blog, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)
	type Items map[string]string
	err = c.Find(bson.M{"type": t}).Sort(field).Limit(n).All(&r)

	return r, nil
}

// FindAndSortByComments 按评论数量由大到小排序，取前n个记录
// TODO: 数据太大时不能全部记录一次取出，资源耗费太大，需采取别的方式
func (b *Blog) FindAndSortByComments(n int) (r []Blog, err error) {
	db, err := NewDBManager()
	defer db.Close()

	var d []Blog
	c := db.session.DB(Name).C(Blogs)
	type Items map[string]string
	err = c.Find(nil).All(&d)

	for i := 0; i < len(d); i++ {
		for j := 0; j < len(d)-i-1; j++ {
			if len(d[j].Comments) < len(d[j+1].Comments) {
				d[j], d[j+1] = d[j+1], d[j]
			}
		}
	}

	if len(d) <= n {
		n = len(d)
	}

	r = d[0:n]

	return r, nil
}

// FindByAndSortByComments 按评论数量由大到小排序，取前n个记录
// TODO: 数据太大时不能全部记录一次取出，资源耗费太大，需采取别的方式
func (b *Blog) FindByAndSortByComments(t Type, n int) (r []Blog, err error) {
	db, err := NewDBManager()
	defer db.Close()

	var d []Blog
	c := db.session.DB(Name).C(Blogs)
	type Items map[string]string
	err = c.Find(bson.M{"type": t}).All(&d)

	for i := 0; i < len(d); i++ {
		for j := 0; j < len(d)-i-1; j++ {
			if len(d[j].Comments) < len(d[j+1].Comments) {
				d[j], d[j+1] = d[j+1], d[j]
			}
		}
	}

	if len(d) <= n {
		n = len(d)
	}

	r = d[0:n]

	return r, nil
}

// Update 修改
func (b *Blog) Update(old, new Blog) error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)

	err = c.Update(old, new)
	if err != nil {
		return err
	}

	return nil
}

// Delete 删除
func (b *Blog) Delete(id string) error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Blogs)
	type Items map[string]string
	err = c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}

	return nil
}
