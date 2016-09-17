package controllers

import (
	m "lelv/app/models"
	"log"
	"time"

	"github.com/revel/revel"
)

func getBlogs(ids []string, c chan int, r *[]m.Blog) {
	blog := m.Blog{}
	for _, id := range ids {
		b, err := blog.FindByID(id)
		if err != nil {
			log.Println(err)
		}
		*r = append(*r, b)
	}
	c <- 1
}

// App app控制器
type App struct {
	*revel.Controller
}

// Home 首页
func (c App) Home() revel.Result {
	cb := m.CarouselBlog{}
	carouselBlogs, err := cb.FindLast(3)
	if err != nil {
		log.Println(err)
	}

	u := m.User{}
	FamousUsers, err := u.FindLast(3)
	if err != nil {
		log.Println(err)
	}

	var (
		HeadBlogs        []m.Blog // 今日博客头条博客集合
		HeadFamousBlog   []m.Blog // 今日推荐名博集合
		NaturalBlogs     []m.Blog // 纯美大自然博客集合
		HistoryBlogs     []m.Blog // 历史.有韵味博客集合
		CustomsBlogs     []m.Blog // 世界.奇异风情博客集合
		ShareBlogs       []m.Blog // 达人分享博客集合
		LatestBlogs      []m.Blog // 及时更新博客集合
		LatestBlogsLeft  []m.Blog // 及时更新博客集合 LatestBlogs 的前一半
		LatestBlogsRight []m.Blog // 及时更新博客集合 LatestBlogs 的后一半
	)

	hbls := m.HomeBlogID{}
	blogIds, err := hbls.FindByTimeStamp(time.Now().Format("2006-01-02"))
	if err != nil {
		log.Println(err)
	}

	const N = 6
	ch := make(chan int, N)
	go getBlogs(blogIds.HeadBlogs, ch, &HeadBlogs)
	go getBlogs(blogIds.HeadFamousBlog, ch, &HeadFamousBlog)
	go getBlogs(blogIds.NaturalBlogs, ch, &NaturalBlogs)
	go getBlogs(blogIds.HistoryBlogs, ch, &HistoryBlogs)
	go getBlogs(blogIds.CustomsBlogs, ch, &CustomsBlogs)
	go getBlogs(blogIds.ShareBlogs, ch, &ShareBlogs)

	for i := 0; i < N; i++ {
		<-ch
	}

	blog := m.Blog{}

	n := 20
	LatestBlogs, err = blog.FindLast(n)
	if err != nil {
		log.Println(err)
	}

	l := len(LatestBlogs)
	half := 0
	if l%2 == 0 {
		half = l / 2
	} else {
		half = l/2 + 1
	}
	LatestBlogsLeft = LatestBlogs[0:half]
	LatestBlogsRight = LatestBlogs[half:]

	viewCountBlogs, err := blog.FindAndSortBy("-viewcount", 10)
	if err != nil {
		log.Println(err)
	}

	commentsBlogs, err := blog.FindAndSortByComments(10)
	if err != nil {
		log.Println(err)
	}
	c.RenderArgs["ViewCountBlogs"] = viewCountBlogs
	c.RenderArgs["CommentsBlogs"] = commentsBlogs

	c.RenderArgs["CarouselBlog"] = carouselBlogs
	c.RenderArgs["HeadBlogs"] = HeadBlogs
	c.RenderArgs["HeadFamousBlog"] = HeadFamousBlog
	c.RenderArgs["NaturalBlogs"] = NaturalBlogs
	c.RenderArgs["HistoryBlogs"] = HistoryBlogs
	c.RenderArgs["CustomsBlogs"] = CustomsBlogs
	c.RenderArgs["ShareBlogs"] = ShareBlogs
	c.RenderArgs["LatestBlogsLeft"] = LatestBlogsLeft
	c.RenderArgs["LatestBlogsRight"] = LatestBlogsRight
	c.RenderArgs["FamousUsers"] = FamousUsers

	return c.Render()
}

// Search 站内搜索功能
func (c App) Search(key string) revel.Result {
	b := m.Blog{}
	bs, err := b.FindByTag(key)
	if err != nil {
		log.Println(err)
	}

	c.RenderArgs["Blogs"] = bs

	return c.Render()
}

// Pictures 图片博客
func (c App) Pictures() revel.Result {
	u := m.User{}
	FamousUsers, err := u.FindLast(3)
	if err != nil {
		log.Println(err)
	}

	var PicturesBlogs []m.Blog
	hbls := m.HomeBlogID{}
	blogIds, err := hbls.FindByTimeStamp(time.Now().Format("2006-01-02"))
	if err != nil {
		log.Println(err)
	}

	blog := m.Blog{}
	for _, id := range blogIds.PicturesBlogs {
		b, err := blog.FindByID(id)
		if err != nil {
			log.Println(err)
		}
		PicturesBlogs = append(PicturesBlogs, b)
	}

	viewCountBlogs, err := blog.FindByAndSortBy(m.Picture, "-viewcount", 10)
	if err != nil {
		log.Println(err)
	}

	commentsBlogs, err := blog.FindByAndSortByComments(m.Picture, 10)
	if err != nil {
		log.Println(err)
	}
	c.RenderArgs["ViewCountBlogs"] = viewCountBlogs
	c.RenderArgs["CommentsBlogs"] = commentsBlogs

	c.RenderArgs["PicturesBlogs"] = PicturesBlogs
	c.RenderArgs["FamousUsers"] = FamousUsers

	return c.Render()
}

// Articles 文字博客
func (c App) Articles() revel.Result {
	u := m.User{}
	FamousUsers, err := u.FindLast(3)
	if err != nil {
		log.Println(err)
	}

	var ArticlesBlogs []m.Blog
	hbls := m.HomeBlogID{}
	blogIds, err := hbls.FindByTimeStamp(time.Now().Format("2006-01-02"))
	if err != nil {
		log.Println(err)
	}

	blog := m.Blog{}
	for _, id := range blogIds.ArticlesBlogs {
		b, err := blog.FindByID(id)
		if err != nil {
			log.Println(err)
		}
		ArticlesBlogs = append(ArticlesBlogs, b)
	}

	viewCountBlogs, err := blog.FindByAndSortBy(m.Text, "-viewcount", 10)
	if err != nil {
		log.Println(err)
	}

	commentsBlogs, err := blog.FindByAndSortByComments(m.Text, 10)
	if err != nil {
		log.Println(err)
	}
	c.RenderArgs["ViewCountBlogs"] = viewCountBlogs
	c.RenderArgs["CommentsBlogs"] = commentsBlogs

	c.RenderArgs["ArticlesBlogs"] = ArticlesBlogs
	c.RenderArgs["FamousUsers"] = FamousUsers

	return c.Render()
}

// About 关于和联系我们
func (c App) About() revel.Result {
	return c.Render()
}

func getBlogByTags(tags []string) []m.Blog {
	mo := m.Blog{}
	var blogs []m.Blog

	for _, tag := range tags {
		t, err := mo.FindByTag(tag)
		if err != nil {
			log.Println(err)
		}
		blogs = append(blogs, t...)
	}
	return blogs
}
