package controllers

import (
	"lelv/app/models/admin"
	"lelv/app/models/blog"
	"lelv/app/models/csblog"
	"lelv/app/qiniu"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/revel/revel"
)

// Admin 管理员（编辑）控制器
type Admin struct {
	*revel.Controller
}

// Home 管理员首页
func (c Admin) Home() revel.Result {

	return c.Render()
}

// SearchForCarousel 站内搜索功能
func (c Admin) SearchForCarousel(key string) revel.Result {
	b := blog.Blog{}
	bs, err := b.FindByTag(key)
	if err != nil {
		log.Println(err)
	}

	c.RenderArgs["Blogs"] = bs

	return c.Render()
}

// SearchForModule 站内搜索功能
func (c Admin) SearchForModule(key string) revel.Result {
	b := blog.Blog{}
	bs, err := b.FindByTag(key)
	if err != nil {
		log.Println(err)
	}

	c.RenderArgs["Blogs"] = bs

	return c.Render()
}

// Carousel 编辑轮播博客页面
func (c Admin) Carousel() revel.Result {

	return c.Render()
}

// PostCarouselBlog 新增轮播博客信息
func (c Admin) PostCarouselBlog() revel.Result {
	id := c.Request.Form["ID"][0]
	b := csblog.CarouselBlog{
		ID:        id,
		Title:     c.Request.Form["Title"][0],
		TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	if c.Request.Form["CoverBase64String"][0] != "" {
		// 剪裁过的文件是PNG格式
		fp := id + ".PNG"
		err := qiniu.DecodeBase64(fp, c.Request.Form["CoverBase64String"][0])
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}

		key := strconv.FormatInt(time.Now().Unix(), 10) + "_" + fp
		err = qiniu.Upload(fp, key)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}

		// 删除临时头像文件
		err = os.Remove(fp)
		if err != nil {
			log.Println("删除临时头像文件失败:", err)
		}

		b.Cover = qiniu.SPACE + key
	}

	log.Println(b)
	b.Add()
	return c.Render()
}

// EditHomeBlog 编辑要在首页显示的博客
func (c Admin) EditHomeBlog(t admin.BlogType, title string) revel.Result {
	c.RenderArgs["BlogType"] = t
	c.RenderArgs["Title"] = title

	return c.Render()
}

// PostEditHomeBlog 编辑要在首页显示的博客
func (c Admin) PostEditHomeBlog() revel.Result {
	idstr := c.Request.Form["IDs"][0]
	idstrarr := strings.Split(idstr, ",")

	ids := []string{}
	for _, id := range idstrarr {
		ids = append(ids, id)
	}

	t, err := strconv.Atoi(c.Request.Form["BlogType"][0])
	if err != nil {
		log.Println(err)
		return c.Redirect(Admin.EditHomeBlog)
	}

	b := admin.HomeBlogID{
		TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	b.AddOrUpdate(admin.BlogType(t), ids)

	return c.Redirect(Admin.EditHomeBlog)
}
