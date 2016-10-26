package app

import (
	"lelv/app/controllers"
	"lelv/app/models/dbmgr"
	"log"

	"github.com/revel/revel"
)

func check(c *revel.Controller) revel.Result {
	if c.Session["UserID"] != "" && c.Session["NickName"] != "" && c.Session["UserID"] != dbmgr.Guest {
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

	if c.Session["UserID"] == dbmgr.Guest && (c.Action == "User.Watch" ||
		c.Action == "User.Collect" || c.Action == "User.ConversationWith") {
		log.Println("游客访问 " + c.Action)
		c.RenderArgs["SigninedUserID"] = dbmgr.Guest
		return c.Redirect(controllers.User.SignIn)
	}

	if c.Action == "User.SignIn" || c.Action == "User.PostSignIn" || c.Action == "User.AllBlogs" ||
		c.Action == "User.SignUp" || c.Action == "User.PostSignUp" ||
		c.Action == "User.Index" || c.Action == "User.SignOut" || c.Action == "User.Fans" ||
		c.Action == "Blog.View" || c.Action == "App.Home" || c.Action == "App.ForStaticHome" ||
		c.Action == "App.Articles" || c.Action == "App.Pictures" ||
		c.Action == "App.About" {

		c.RenderArgs["UserID"] = dbmgr.Guest
		c.RenderArgs["NickName"] = dbmgr.Guest
		c.RenderArgs["Avatar"] = dbmgr.Guest
		c.Session["UserID"] = dbmgr.Guest
		c.Session["NickName"] = dbmgr.Guest
		c.Session["Avatar"] = dbmgr.Guest

		log.Println("游客访问 " + c.Action)
		return nil
	}

	return c.Redirect(controllers.App.Home)
}
