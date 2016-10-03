package controllers

import (
	m "lelv/app/models"
	"log"

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
	cn1 := 3  // 第一部分博客数
	cn2 := 10 // 第二部分博客数
	cn3 := 20 // 第三部分博客数
	cn4 := 10 // 阅读排行榜博客数
	cn5 := 10 // 评论排行榜博客数
	cn6 := 10 // 轮播博客数

	// 轮播图博客取值
	cb := m.CarouselBlog{}
	carouselBlogs, err := cb.FindLast(cn6)
	if err != nil {
		log.Println(err)
	}

	var (
		latestBlogs []m.Blog // 最新博客集合
		part1       []m.Blog // 首页第一部分错层博客
		part2Left   []m.Blog // 首页第二部分左边列表博客
		part2Right  []m.Blog // 首页第二部分右边列表博客
		part3Left   []m.Blog // 首页第三部分左边列表博客
		part3Right  []m.Blog // 首页第三部分右边列表博客
	)
	blog := m.Blog{}

	n := cn1 + cn2 + cn3
	latestBlogs, err = blog.FindLast(n)
	if err != nil {
		log.Println(err)
	}
	log.Println(latestBlogs)
	count := len(latestBlogs)

	// 第一部分错层博客取值
	if count > cn1 {
		part1 = latestBlogs[0:cn1]
	} else {
		part1 = latestBlogs[0:count]
	}

	// 第二部分列表博客取值
	if count > (cn1 + cn2) {
		part2Left = latestBlogs[cn1 : cn1+cn2/2]
		part2Right = latestBlogs[cn1+cn2/2 : cn1+cn2]
	} else {
		half1 := 0
		if (count-cn1)%2 == 0 {
			half1 = (count-cn1)/2 + cn1
		} else {
			half1 = (count-cn1)/2 + 1 + cn1
		}
		part2Left = latestBlogs[cn1:half1]
		part2Right = latestBlogs[half1:]
	}

	// 第三部分列表博客取值
	if count > (cn1 + cn2 + cn3) {
		part3Left = latestBlogs[cn1+cn2 : cn1+cn2+cn3/2]
		part3Right = latestBlogs[cn1+cn2+cn3/2 : cn1+cn2+cn3]
	} else if count > (cn1 + cn2) {
		half := 0
		if (count-cn1-cn2)%2 == 0 {
			half = (count-cn1-cn2)/2 + cn1 + cn2
		} else {
			half = (count-cn1-cn2)/2 + 1 + cn1 + cn2
		}
		part3Left = latestBlogs[cn1+cn2 : half]
		part3Right = latestBlogs[half:]
	}

	// 阅读排行榜取值
	viewCountBlogs, err := blog.FindAndSortBy("-viewcount", cn4)
	if err != nil {
		log.Println(err)
	}

	// 评论排行榜取值
	commentsBlogs, err := blog.FindAndSortByComments(cn5)
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

// Home 首页
// func (c App) Home() revel.Result {
// 	cb := m.CarouselBlog{}
// 	carouselBlogs, err := cb.FindLast(3)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	u := m.User{}
// 	FamousUsers, err := u.FindLast(3)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	var (
// 		HeadBlogs        []m.Blog // 今日博客头条博客集合
// 		HeadFamousBlog   []m.Blog // 今日推荐名博集合
// 		NaturalBlogs     []m.Blog // 纯美大自然博客集合
// 		HistoryBlogs     []m.Blog // 历史.有韵味博客集合
// 		CustomsBlogs     []m.Blog // 世界.奇异风情博客集合
// 		ShareBlogs       []m.Blog // 达人分享博客集合
// 		LatestBlogs      []m.Blog // 及时更新博客集合
// 		LatestBlogsLeft  []m.Blog // 及时更新博客集合 LatestBlogs 的前一半
// 		LatestBlogsRight []m.Blog // 及时更新博客集合 LatestBlogs 的后一半
// 	)

// 	hbls := m.HomeBlogID{}
// 	//blogIds, err := hbls.FindByTimeStamp(time.Now().Format("2006-01-02"))
// 	blogIds, err := hbls.GetLast()
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	const N = 6
// 	ch := make(chan int, N)
// 	go getBlogs(blogIds.HeadBlogs, ch, &HeadBlogs)
// 	go getBlogs(blogIds.HeadFamousBlog, ch, &HeadFamousBlog)
// 	go getBlogs(blogIds.NaturalBlogs, ch, &NaturalBlogs)
// 	go getBlogs(blogIds.HistoryBlogs, ch, &HistoryBlogs)
// 	go getBlogs(blogIds.CustomsBlogs, ch, &CustomsBlogs)
// 	go getBlogs(blogIds.ShareBlogs, ch, &ShareBlogs)

// 	for i := 0; i < N; i++ {
// 		<-ch
// 	}

// 	blog := m.Blog{}

// 	n := 20
// 	LatestBlogs, err = blog.FindLast(n)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	log.Println(LatestBlogs)

// 	l := len(LatestBlogs)
// 	half := 0
// 	if l%2 == 0 {
// 		half = l / 2
// 	} else {
// 		half = l/2 + 1
// 	}
// 	LatestBlogsLeft = LatestBlogs[0:half]
// 	LatestBlogsRight = LatestBlogs[half:]

// 	var StaggeredBlogs []m.Blog
// 	if l > 3 {
// 		StaggeredBlogs = LatestBlogs[0:3]
// 	} else {
// 		StaggeredBlogs = LatestBlogs[0 : l]
// 	}

// 	viewCountBlogs, err := blog.FindAndSortBy("-viewcount", 10)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	commentsBlogs, err := blog.FindAndSortByComments(10)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	c.RenderArgs["ViewCountBlogs"] = viewCountBlogs
// 	c.RenderArgs["CommentsBlogs"] = commentsBlogs

// 	c.RenderArgs["CarouselBlog"] = carouselBlogs
// 	c.RenderArgs["HeadBlogs"] = HeadBlogs
// 	c.RenderArgs["HeadFamousBlog"] = HeadFamousBlog
// 	c.RenderArgs["NaturalBlogs"] = NaturalBlogs
// 	c.RenderArgs["HistoryBlogs"] = HistoryBlogs
// 	c.RenderArgs["CustomsBlogs"] = CustomsBlogs
// 	c.RenderArgs["ShareBlogs"] = ShareBlogs
// 	c.RenderArgs["StaggeredBlogs"] = StaggeredBlogs
// 	c.RenderArgs["LatestBlogsLeft"] = LatestBlogsLeft
// 	c.RenderArgs["LatestBlogsRight"] = LatestBlogsRight
// 	c.RenderArgs["FamousUsers"] = FamousUsers

// 	return c.Render()
// }

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
	//blogIds, err := hbls.FindByTimeStamp(time.Now().Format("2006-01-02"))
	blogIds, err := hbls.GetLast()
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
	//blogIds, err := hbls.FindByTimeStamp(time.Now().Format("2006-01-02"))
	blogIds, err := hbls.GetLast()
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
