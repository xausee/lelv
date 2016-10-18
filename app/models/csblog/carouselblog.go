package csblog

import (
	"lelv/app/models/dbmgr"
	"log"
)

// CarouselBlog 用于头条的轮播博客信息，只需包含ID号，标题和博客封面图片
type CarouselBlog struct {
	ID        string // ID号
	Title     string // 标题
	Cover     string // 博客封面图片地址，图片大小一致
	TimeStamp string // 时间戳
}

// Add 新增
func (b *CarouselBlog) Add() error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.CarouselBlogs)

	err = c.Insert(b)
	if err != nil {
		log.Println("创建博客失败：")
		log.Println(b)
		return err
	}

	return nil
}

// FindLast 查找最新的n个记录
func (b *CarouselBlog) FindLast(n int) (r []CarouselBlog, err error) {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.CarouselBlogs)
	err = c.Find(nil).Sort("-timestamp").Limit(n).All(&r)

	return r, nil
}
