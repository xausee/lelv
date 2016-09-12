package app

import (
	"log"
	"lelv/app/controllers"
	model "lelv/app/models"

	"github.com/revel/revel"
)

func check(c *revel.Controller) revel.Result {
	if c.Session["UserID"] != "" && c.Session["NickName"] != "" && c.Session["UserID"] != model.Guest {
		c.RenderArgs["UserID"] = c.Session["UserID"]
		c.RenderArgs["NickName"] = c.Session["NickName"]
		c.RenderArgs["Avatar"] = c.Session["Avatar"]

		count := controllers.GetUnreadMsgCount(c.Session["UserID"])
		if count > 0 {
			c.RenderArgs["UnreadMsgCount"] = count
		}
		log.Println("用户 " + c.Session["NickName"] + " 访问 " + c.Action)
		return nil
	}

	if c.Session["UserID"] == model.Guest && (c.Action == "User.Watch" ||
		c.Action == "User.Collect" || c.Action == "User.ConversationWith") {
		log.Println("游客访问 " + c.Action)
		c.RenderArgs["SigninedUserID"] = model.Guest
		return c.Redirect(controllers.User.SignIn)
	}

	if c.Action == "User.SignIn" || c.Action == "User.PostSignIn" || c.Action == "User.AllBlogs" ||
		c.Action == "User.SignUp" || c.Action == "User.PostSignUp" ||
		c.Action == "User.Index" || c.Action == "User.SignOut" || c.Action == "User.Fans" ||
		c.Action == "Blog.View" || c.Action == "App.Home" ||
		c.Action == "App.Articles" || c.Action == "App.Pictures" ||
		c.Action == "App.About" {

		c.RenderArgs["UserID"] = model.Guest
		c.RenderArgs["NickName"] = model.Guest
		c.RenderArgs["Avatar"] = model.Guest
		c.Session["UserID"] = model.Guest
		c.Session["NickName"] = model.Guest
		c.Session["Avatar"] = model.Guest

		log.Println("游客访问 " + c.Action)
		return nil
	}

	return c.Redirect(controllers.App.Home)
}
