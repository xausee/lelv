package controllers

import (
	"lelv/app/models/blog"
	"lelv/app/models/dbmgr"
	"lelv/app/models/user"
	qiniu "lelv/app/qiniu"
	qiniumock "lelv/app/qiniumock"
	"lelv/app/util"
	"log"

	"strings"
	"time"

	"github.com/revel/revel"
)

// Blog 博客控制器
type Blog struct {
	*revel.Controller
}

// Create 写博客页面
func (c Blog) Create() revel.Result {
	if revel.RunMode == "dev" {
		token := qiniumock.CreatUpToken()
		log.Println("生成七牛上传凭证：" + token)
		UpToken:=token
		CDN:=qiniumock.SPACE
		c.Render(UpToken,CDN)
	} else {
		token := qiniu.CreatUpToken()
		log.Println("生成七牛上传凭证：" + token)
		UpToken:=token
		CDN:=qiniu.CDN
		c.Render(UpToken,CDN)
	}

	return c.Render()
}

// PostBlog 写博客页面，处理POST请求，从前端获取数据并写入到数据库
func (c Blog) PostBlog(b blog.Blog) revel.Result {
	var t blog.Type
	switch c.Request.Form["type"][0] {
	case "Picture":
		t = blog.Picture
	case "Text":
		t = blog.Text
	case "Hybrid":
		t = blog.Hybrid
	}

	tags := strings.Split(c.Request.Form["tags"][0], ",")
	pictures := strings.Split(c.Request.Form["pictures"][0], ",")

	b = blog.Blog{
		ID:                  util.CreateObjectID(),
		AuthorID:            c.Session["UserID"],
		Author:              c.Session["NickName"],
		Tags:                tags,
		Type:                t,
		Title:               c.Request.Form["title"][0],
		Cover:               c.Request.Form["cover"][0],
		BriefText:           c.Request.Form["briefText"][0],
		Content:             c.Request.Form["content"][0],
		Pictures:            pictures,
		CreateTimeStamp:     time.Now().Format("2006-01-02 15:04:05"),
		LastUpdateTimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	err := blog.Add(b)
	if err != nil {
		log.Println(err)
		return c.RenderText("博客创建失败：" + err.Error())
	}

	return c.RenderText(b.ID)
}

// View 查看博客
func (c Blog) View(id string) revel.Result {
	// 获取博客信息
	b := blog.Blog{}
	b.ID = id
	blog, err := blog.FindByID(id)
	if err != nil {
		log.Println("获取博客失败， ID：" + id)
	}

	// 获取作者信息
	aid := blog.AuthorID
	author, err := user.FindByID(aid)
	if err != nil {
		log.Println(err)
		return c.Render()
	}

	collected := false
	if c.Session["UserID"] != "" && c.Session["UserID"] != dbmgr.Guest {
		user, err := user.FindByID(c.Session["UserID"])
		if err != nil {
			log.Println(err)
			return c.RenderText("查找用户失败")
		}
		for i := range user.Collection {
			if user.Collection[i] == id {
				collected = true
			}
		}
	}

	isAuthor := false
	if aid == c.Session["UserID"] {
		isAuthor = true
	}

	b.UpdateView()
	ViewCount:= blog.ViewCount
	Collected:= collected
	Author:= author
	IsAuthor:= isAuthor
	Blog:= blog

	SigninedUserID:= c.Session["UserID"]

	return c.Render(ViewCount,Collected,Author,IsAuthor,Blog,SigninedUserID)
}

// Edit 编辑博客
func (c Blog) Edit(id string) revel.Result {
	// 获取博客信息
	b := blog.Blog{}
	b.ID = id
	blog, err := blog.FindByID(id)
	if err != nil {
		log.Println("获取博客失败， ID：" + id)
	}

	if revel.RunMode == "dev" {
		token := qiniumock.CreatUpToken()
		log.Println("生成七牛上传凭证：" + token)
		UpToken:=token
		CDN:=qiniumock.SPACE
		c.Render(UpToken,CDN)
	} else {
		token := qiniu.CreatUpToken()
		log.Println("生成七牛上传凭证：" + token)
		UpToken:=token
		CDN:=qiniu.CDN
		c.Render(UpToken,CDN)
	}

Blog:= blog
	return c.Render(Blog)
}

// PostEdit 编辑博客 Post请求
func (c Blog) PostEdit() revel.Result {
	// 获取博客信息
	id := c.Request.Form["id"][0]
	b := blog.Blog{}
	b.ID = id
	OldBlog, err := blog.FindByID(id)
	if err != nil {
		log.Println("获取博客失败， ID：" + id)
	}

	var t blog.Type
	switch c.Request.Form["type"][0] {
	case "Picture":
		t = blog.Picture
	case "Text":
		t = blog.Text
	case "Hybrid":
		t = blog.Hybrid
	}

	tags := strings.Split(c.Request.Form["tags"][0], ",")
	pictures := strings.Split(c.Request.Form["pictures"][0], ",")
	editedblog := OldBlog
	editedblog.Tags = tags
	editedblog.Type = t
	editedblog.Title = c.Request.Form["title"][0]
	editedblog.BriefText = c.Request.Form["briefText"][0]
	editedblog.Content = c.Request.Form["content"][0]
	editedblog.Pictures = pictures
	editedblog.LastUpdateTimeStamp = time.Now().Format("2006-01-02 15:04:05")
	if c.Request.Form["cover"][0] != "" {
		editedblog.Cover = c.Request.Form["cover"][0]
	} else {
		editedblog.Cover = OldBlog.Cover
	}

	err = blog.Update(OldBlog, editedblog)
	if err != nil {
		log.Println(err)
		return c.RenderText("更新失败：" + err.Error())
	}

	return c.RenderText(editedblog.ID)
}

// Delete 删除博客
func (c Blog) Delete(id string) revel.Result {
	a, err := blog.FindByID(id)
	if err != nil {
		log.Println(err)
		return c.RenderText("查找博客失败：" + err.Error())
	}

	// 删除云存储上的博客图片
	for _, p := range a.Pictures {
		if revel.RunMode == "dev" {
			qiniumock.Delete(p)
		} else {
			qiniu.Delete(p)
		}
	}

	// 删除博客文本
	err = blog.Delete(id)
	if err != nil {
		log.Println(err)
		return c.RenderText("删除失败：" + err.Error())
	}

	return c.RenderText("成功删除博客")
}

// PostComment 发表评论, POST 数据处理
func (c Blog) PostComment(comment blog.Comment) revel.Result {
	user, err := user.FindByID(c.Session["UserID"])
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	blogID := c.Request.Form["BlogIDForComment"][0]
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	b := blog.Blog{}
	b.ID = blogID
	blog, err := blog.FindByID(blogID)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	comment.ID = util.CreateObjectID()
	comment.CommenterID = c.Session["UserID"]
	comment.CommenterAvatar = user.Avatar
	comment.CommenterNickName = c.Session["NickName"]
	comment.TimeStamp = time.Now().Format("2006-01-02 15:04:05")

	err = blog.AddComment(comment)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	return c.Redirect("/Blog/View?id=" + blogID)
}
