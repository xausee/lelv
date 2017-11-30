package app

import (
	"lelv/app/controllers"
	"lelv/app/models/dbmgr"
	"lelv/app/models/user"
	"log"

	"strconv"

	"github.com/revel/revel"
)

func check(c *revel.Controller) revel.Result {
	if c.Session["UserID"] != "" && c.Session["NickName"] != "" &&
		c.Session["Role"] != "" && c.Session["UserID"] != dbmgr.Guest {

		UserID := c.Session["UserID"]
		NickName := c.Session["NickName"]
		Avatar := c.Session["Avatar"]
		Role := c.Session["Role"]
		c.ViewArgs["UserID"] = UserID
		c.ViewArgs["NickName"] = NickName
		c.ViewArgs["Avatar"] = Avatar
		c.ViewArgs["Role"] = Role

		if c.Action == "Admin.Home" {
			i, err := strconv.Atoi(c.Session["Role"])
			if err != nil {
				log.Println(err)
				return c.Redirect(controllers.App.Home)
			}

			if i == (int)(user.Super) || i == (int)(user.Admin) {
				return nil
			}

			log.Println("用户 " + c.Session["NickName"] + " 尝试访问未授权Admin页面，自动跳转到首页")
			return c.Redirect(controllers.App.Home)
		}

		count := controllers.GetUnreadMsgCount(c.Session["UserID"])
		if count > 0 {
			UnreadMsgCount := count
			c.ViewArgs["UnreadMsgCount"] = UnreadMsgCount
		}
		log.Println("用户 " + c.Session["NickName"] + " 访问 " + c.Action)
		return nil
	}

	if c.Session["UserID"] == dbmgr.Guest && (c.Action == "User.Watch" ||
		c.Action == "User.Collect" || c.Action == "User.ConversationWith") {
		log.Println("游客访问 " + c.Action)
		SigninedUserID := dbmgr.Guest
		c.ViewArgs["SigninedUserID"] = SigninedUserID
		return c.Redirect(controllers.User.SignIn)
	}

	if c.Action == "User.SignIn" || c.Action == "User.PostSignIn" || c.Action == "User.AllBlogs" ||
		c.Action == "User.SignUp" || c.Action == "User.PostSignUp" ||
		c.Action == "User.Index" || c.Action == "User.SignOut" || c.Action == "User.Fans" ||
		c.Action == "Blog.View" || c.Action == "App.Home" || c.Action == "App.DynamicHome" ||
		c.Action == "App.Articles" || c.Action == "App.Pictures" ||
		c.Action == "App.About" {

		UserID := dbmgr.Guest
		NickName := dbmgr.Guest
		Avatar := dbmgr.Guest
		c.Session["UserID"] = UserID
		c.Session["NickName"] = NickName
		c.Session["Avatar"] = Avatar

		c.ViewArgs["UserID"] = UserID
		c.ViewArgs["NickName"] = NickName
		c.ViewArgs["Avatar"] = Avatar

		log.Println("游客访问 " + c.Action)
		return nil
	}

	return c.Redirect(controllers.App.Home)
}
