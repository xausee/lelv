package controllers

import (
	"lelv/app/models/admin"
	"lelv/app/models/blog"
	"lelv/app/models/csblog"
	"lelv/app/models/user"
	"log"

	"github.com/revel/revel"
)

// App app控制器
type App struct {
	*revel.Controller
}

// Home 首页
func (c App) Home() revel.Result {
	cn1 := 3  // 第一部分博客数
	cn2 := 10 // 第二部分博客数
	cn3 := 20 // 第三部分博客数
	cn4 := 10 // 阅读排行榜博客数
	cn5 := 10 // 评论排行榜博客数
	cn6 := 10 // 轮播博客数

	// 轮播图博客取值
	cb := csblog.CarouselBlog{}
	carouselBlogs, err := cb.FindLast(cn6)
	if err != nil {
		log.Println(err)
	}

	var (
		latestBlogs []blog.Blog // 最新博客集合
		part1       []blog.Blog // 首页第一部分错层博客
		part2Left   []blog.Blog // 首页第二部分左边列表博客
		part2Right  []blog.Blog // 首页第二部分右边列表博客
		part3Left   []blog.Blog // 首页第三部分左边列表博客
		part3Right  []blog.Blog // 首页第三部分右边列表博客
	)

	n := cn1 + cn2 + cn3
	latestBlogs, err = blog.FindLast(n)
	if err != nil {
		log.Println(err)
	}
	count := len(latestBlogs)

	// 第一部分错层博客取值
	if count > cn1 {
		part1 = latestBlogs[0:cn1]
	} else {
		part1 = latestBlogs[0:count]
	}

	// 第二部分列表博客取值
	if count > cn1 && count < (cn1+cn2) {
		half1 := 0
		if (count-cn1)%2 == 0 {
			half1 = (count-cn1)/2 + cn1
		} else {
			half1 = (count-cn1)/2 + 1 + cn1
		}
		part2Left = latestBlogs[cn1:half1]
		part2Right = latestBlogs[half1:]
	} else if count >= (cn1 + cn2) {
		part2Left = latestBlogs[cn1 : cn1+cn2/2]
		part2Right = latestBlogs[cn1+cn2/2 : cn1+cn2]
	}

	// 第三部分列表博客取值
	if count > (cn1+cn2) && count < (cn1+cn2+cn3) {
		half := 0
		if (count-cn1-cn2)%2 == 0 {
			half = (count-cn1-cn2)/2 + cn1 + cn2
		} else {
			half = (count-cn1-cn2)/2 + 1 + cn1 + cn2
		}
		part3Left = latestBlogs[cn1+cn2 : half]
		part3Right = latestBlogs[half:]
	} else if count >= (cn1 + cn2 + cn3) {
		part3Left = latestBlogs[cn1+cn2 : cn1+cn2+cn3/2]
		part3Right = latestBlogs[cn1+cn2+cn3/2 : cn1+cn2+cn3]
	}

	// 阅读排行榜取值
	viewCountBlogs, err := blog.FindAndSortBy("-viewcount", cn4)
	if err != nil {
		log.Println(err)
	}

	// 评论排行榜取值
	commentsBlogs, err := blog.FindALLSortByCommentsNum(cn5)
	if err != nil {
		log.Println(err)
	}
	c.RenderArgs["ViewCountBlogs"] = viewCountBlogs
	c.RenderArgs["CommentsBlogs"] = commentsBlogs
	c.RenderArgs["CarouselBlog"] = carouselBlogs
	c.RenderArgs["Part1Blogs"] = part1
	c.RenderArgs["Part2LeftBlogs"] = part2Left
	c.RenderArgs["Part2RightBlogs"] = part2Right
	c.RenderArgs["Part3LeftBlogs"] = part3Left
	c.RenderArgs["Part3RightBlogs"] = part3Right

	return c.Render()
}

// Search 站内搜索功能
func (c App) Search(key string) revel.Result {
	bs, err := blog.FindByTag(key)
	if err != nil {
		log.Println(err)
	}

	c.RenderArgs["Blogs"] = bs

	return c.Render()
}

// Pictures 图片博客
func (c App) Pictures() revel.Result {
	u := user.User{}
	FamousUsers, err := u.FindLast(3)
	if err != nil {
		log.Println(err)
	}

	var PicturesBlogs []blog.Blog
	hbls := admin.HomeBlogID{}
	//blogIds, err := hbls.FindByTimeStamp(time.Now().Format("2006-01-02"))
	blogIds, err := hbls.GetLast()
	if err != nil {
		log.Println(err)
	}

	for _, id := range blogIds.PicturesBlogs {
		b, err := blog.FindByID(id)
		if err != nil {
			log.Println(err)
		}
		PicturesBlogs = append(PicturesBlogs, b)
	}

	viewCountBlogs, err := blog.FindByAndSortBy(blog.Picture, "-viewcount", 10)
	if err != nil {
		log.Println(err)
	}

	commentsBlogs, err := blog.FindByTypeAndSortByCoNum(blog.Picture, 10)
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
	u := user.User{}
	FamousUsers, err := u.FindLast(3)
	if err != nil {
		log.Println(err)
	}

	var ArticlesBlogs []blog.Blog
	hbls := admin.HomeBlogID{}
	//blogIds, err := hbls.FindByTimeStamp(time.Now().Format("2006-01-02"))
	blogIds, err := hbls.GetLast()
	if err != nil {
		log.Println(err)
	}

	for _, id := range blogIds.ArticlesBlogs {
		b, err := blog.FindByID(id)
		if err != nil {
			log.Println(err)
		}
		ArticlesBlogs = append(ArticlesBlogs, b)
	}

	viewCountBlogs, err := blog.FindByAndSortBy(blog.Text, "-viewcount", 10)
	if err != nil {
		log.Println(err)
	}

	commentsBlogs, err := blog.FindByTypeAndSortByCoNum(blog.Text, 10)
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

func getBlogByTags(tags []string) (bs []blog.Blog) {
	for _, tag := range tags {
		t, err := blog.FindByTag(tag)
		if err != nil {
			log.Println(err)
		}
		bs = append(bs, t...)
	}
	return
}
