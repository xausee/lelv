package controllers

import (
	"lelv/app/models/blog"
	"lelv/app/models/conversation"
	"lelv/app/models/user"
	qiniu "lelv/app/qiniu"
	qiniumock "lelv/app/qiniumock"
	"log"
	"os"
	"strings"
	"time"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

// User 用户控制器
type User struct {
	*revel.Controller
}

// Home 用户首页
func (c User) Home() revel.Result {
	id := c.Session["UserID"]

	u := user.User{}
	user, err := u.FindByID(id)
	if err != nil {
		log.Println(err)
		return c.Render()
	}

	b := blog.Blog{}
	blogs, err := b.FindByAuthorID(id)
	if err != nil {
		log.Println(err)
		return c.Render()
	}

	// 首页无需加载所有内容，清空content，减少数据传送
	for _, blog := range blogs {
		blog.Content = ""
	}

	c.RenderArgs["User"] = user
	c.RenderArgs["Blogs"] = blogs
	c.RenderArgs["BlogsCount"] = len(blogs)
	c.RenderArgs["CollectionCount"] = len(user.Collection)
	c.RenderArgs["FansCount"] = len(user.Fans)
	c.RenderArgs["WatchCount"] = len(user.Watches)
	c.RenderArgs["ConversationCount"] = len(user.ConversationIDs)

	return c.Render()
}

// Index 用户对其它用户开放的首页
func (c User) Index(idnum string) revel.Result {
	u := user.User{}
	us, err := u.FindByID(idnum)

	if err != nil {
		log.Println(err)
		return c.RenderText("该用户不存在，或者已经被注销")
	}

	b := blog.Blog{}
	blogs, err := b.FindByAuthorID(idnum)
	if err != nil {
		log.Println(err)
		return c.RenderText(err.Error())
	}

	// 首页无需加载所有内容，清空content，减少数据传送
	for _, blog := range blogs {
		blog.Content = ""
	}

	watched := false
	sid := c.Session["UserID"]

	if sid != "" && sid != "Guest" {
		su := user.User{}
		signinedUser, err := su.FindByID(sid)
		if err != nil {
			log.Println(err)
			return c.RenderText("该用户不存在，或者已经被注销")
		}

		for _, w := range signinedUser.Watches {
			if w == idnum {
				watched = true
			}
		}
	}

	c.RenderArgs["User"] = us
	c.RenderArgs["Blogs"] = blogs
	c.RenderArgs["BlogsCount"] = len(blogs)
	c.RenderArgs["FansCount"] = len(us.Fans)
	c.RenderArgs["WatchCount"] = len(us.Watches)
	c.RenderArgs["Watched"] = watched

	c.RenderArgs["SigninedUserID"] = sid

	return c.Render()
}

// Avatar 用户修改资料页面控制器
func (c User) Avatar() revel.Result {
	return c.Render()
}

// Profile 用户修改资料页面控制器
func (c User) Profile() revel.Result {
	id := c.Session["UserID"]

	u := user.User{}
	user, err := u.FindByID(id)
	if err != nil {
		log.Println(err)
		return c.Render()
	}

	token := qiniu.CreatUpToken()
	log.Println("生成七牛上传凭证：" + token)

	c.RenderArgs["UpToken"] = token
	c.RenderArgs["User"] = user

	return c.Render()
}

// PostProfile 用户修改资料Post请求控制器
func (c User) PostProfile() revel.Result {
	u := user.User{
		ID: c.Session["UserID"],
		//NickName:     c.Request.Form["NickName"][0],
		NickName:            c.Session["NickName"],
		Introduction:        c.Request.Form["Introduction"][0],
		LastUpdateTimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	user, err := u.FindByID(c.Session["UserID"])
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}
	u.Avatar = user.Avatar

	if c.Request.Form["Base64String"] != nil && c.Request.Form["Base64String"][0] != "" {
		// 剪裁过的文件是PNG格式
		fp := c.Session["UserID"] + ".PNG"
		err = qiniu.DecodeBase64(fp, c.Request.Form["Base64String"][0])
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}

		// 先移动（修改）七牛上原来的文件
		oldkey := strings.Replace(user.Avatar, qiniu.SPACE, "", -1)
		log.Println(c.Request.Form["Base64String"])
		ofp := "O_" + oldkey
		if user.Avatar != qiniu.DefaultMaleAvatar && user.Avatar != qiniu.DefaultFemaleAvatar {
			qiniu.Move(oldkey, ofp)
			if err != nil {
				log.Println(err)
				return c.Render(err)
			}
		}

		// 上传本地服务器文件到七牛
		key := fp
		err = qiniu.Upload(fp, key)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}

		// 上传成功后删除本地服务器临时头像文件
		err = os.Remove(fp)
		if err != nil {
			log.Println("删除临时头像文件失败:", err)
		}

		// 上传成功后从七牛上删除原来的头像图片
		if user.Avatar != qiniu.DefaultMaleAvatar && user.Avatar != qiniu.DefaultFemaleAvatar {
			qiniu.Delete(ofp)
		}

		u.Avatar = qiniu.SPACE + key
	}

	err = u.Update()
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	c.Session["NickName"] = u.NickName
	return c.Redirect(User.Profile)
}

// SignUp 用户注册页，处理GET请求，加载页面内容
func (c User) SignUp() revel.Result {
	return c.Render()
}

// PostSignUp 用户注册页，处理POST请求，从前端获取数据并写入到数据库
func (c User) PostSignUp(mockuser user.MockUser) revel.Result {
	c.Validation.Required(mockuser.NickName).Message("用户昵称不能为空")
	c.Validation.Required(mockuser.Password).Message("密码不能为空")
	c.Validation.Required(mockuser.ConfirmPassword).Message("确认密码不能为空")
	c.Validation.MinSize(mockuser.Password, 6).Message("密码长度不短于6位")
	c.Validation.Required(mockuser.ConfirmPassword == mockuser.Password).Message("两次输入的密码不一致")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(User.SignUp)
	}

	var avatar string

	if revel.RunMode == "dev" {
		if mockuser.Gender == "男" {
			avatar = qiniumock.DefaultMaleAvatar
		} else {
			avatar = qiniumock.DefaultFemaleAvatar
		}
	} else {
		if mockuser.Gender == "男" {
			avatar = qiniu.DefaultMaleAvatar
		} else {
			avatar = qiniu.DefaultFemaleAvatar
		}
	}

	p, _ := bcrypt.GenerateFromPassword([]byte(mockuser.Password), bcrypt.DefaultCost)
	u := user.User{
		NickName:        mockuser.NickName,
		Gender:          mockuser.Gender,
		Avatar:          avatar,
		Introduction:    "这家伙有点懒，没有留下一点信息...",
		Password:        p,
		CreateTimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	u.ID = conversation.CreateObjectID()

	err := u.SignUp()
	if err != nil {
		c.Validation.Clear()
		log.Println(err)
		// 添加错误信息，显示在用户名下面
		e := revel.ValidationError{
			Message: err.Error(),
			Key:     "mockuser.NickName"}

		c.Validation.Errors = append(c.Validation.Errors, &e)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(User.SignUp)
	}

	return c.Redirect(User.SignIn)
}

// SignIn 用户登录页
func (c User) SignIn(redirect string) revel.Result {

	log.Println(redirect)
	if redirect != "" {
		c.RenderArgs["RedirectTo"] = redirect
	}

	return c.Render()
}

// PostSignIn 用户登录页，处理POST请求，从前端获取数据并写入到数据库
func (c User) PostSignIn(mockuser user.MockUser) revel.Result {
	c.Validation.Required(mockuser.NickName).Message("用户昵称不能为空")
	c.Validation.Required(mockuser.Password).Message("密码不能为空")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(User.SignIn)
	}

	u := user.User{
		NickName: mockuser.NickName,
		Password: []byte(mockuser.Password)}

	err := u.SignIn()
	if err != nil {
		c.Validation.Clear()
		log.Println(err)
		// 添加错误信息，显示在用户名下面
		e := revel.ValidationError{
			Message: err.Error(),
			Key:     "mockuser.NickName"}

		c.Validation.Errors = append(c.Validation.Errors, &e)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(User.SignIn)
	}

	c.Session["UserID"] = u.ID
	c.Session["NickName"] = u.NickName
	c.Session["Avatar"] = u.Avatar
	c.RenderArgs["UserID"] = u.ID
	c.RenderArgs["NickName"] = u.NickName
	c.RenderArgs["Avatar"] = u.Avatar

	if c.Session["NickName"] == "狂赞士之怒" {
		return c.Redirect(Admin.Home)
	}

	if c.Request.Form["RedirectTo"][0] != "" {
		redirectUrl := c.Request.Form["RedirectTo"][0]
		return c.Redirect(redirectUrl)
	}

	return c.Redirect(User.Home)
}

// SignOut 用注销页
func (c User) SignOut() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}

	return c.Redirect(App.Home)
}

// Collect 收藏/取消收藏博客
func (c User) Collect() revel.Result {
	id := c.Request.Form["BlogID"][0]

	response := ""
	u := user.User{}
	u.ID = c.Session["UserID"]

	if c.Request.Form["Flag"][0] == "收藏" {
		err := u.Collect(id, true)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}

		response = "取消收藏"
	} else {
		err := u.Collect(id, false)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}
		response = "收藏"
	}

	return c.RenderText(response)
}

// Collection 获取所有收藏的博客
func (c User) Collection() revel.Result {
	uid := c.Session["UserID"]

	u := user.User{}
	user, err := u.FindByID(uid)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	b := blog.Blog{}
	blogs := []blog.Blog{}
	for _, id := range user.Collection {
		blog, err := b.FindByID(id)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}
		blogs = append(blogs, blog)
	}

	c.RenderArgs["Blogs"] = blogs
	return c.Render()
}

// AllBlogs 获取所有博客
func (c User) AllBlogs() revel.Result {
	// 用户访问自己的所有博客页面
	id := c.Session["UserID"]

	// 其他用户访问
	if len(c.Request.Form["UserID"]) != 0 && c.Request.Form["UserID"][0] != "" {
		id = c.Request.Form["UserID"][0]
	}

	b := blog.Blog{}
	blogs, err := b.FindByAuthorID(id)
	if err != nil {
		log.Println(err)
		return c.Render()
	}

	// 首页无需加载所有内容，清空content，减少数据传送
	for _, blog := range blogs {
		blog.Content = ""
	}

	c.RenderArgs["Blogs"] = blogs

	return c.Render()
}

// Watch 关注某个用户
func (c User) Watch() revel.Result {
	id := c.Request.Form["UserID"][0]
	uid := c.Session["UserID"]

	response := ""
	u := user.User{}
	u.ID = uid
	// currentUser, err := u.FindByID(uid)

	wu := user.User{}
	wu.ID = id

	if strings.Contains(c.Request.Form["Flag"][0], "取消关注") {
		// 将被取消关注的用户从此用户的关注列表中移出
		err := u.UpdateWatch(id, false)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}
		response = "关注"

		// 将此用户从被关注用户的粉丝列表中移出
		err = wu.UpdateFans(uid, false)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}
	} else {
		// 将被关注的用户添加到此用户的关注列表
		err := u.UpdateWatch(id, true)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}
		response = "取消关注"

		// 将此用户添加到被关注用户的粉丝列表
		err = wu.UpdateFans(uid, true)
		if err != nil {
			log.Println(err)
			return c.Render(err)
		}
	}

	return c.RenderText(response)
}

// Watches 获取所有关注的用户信息
func (c User) Watches() revel.Result {
	// 用户访问自己的所有粉丝
	id := c.Session["UserID"]

	u := user.User{}
	us, err := u.FindByID(id)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	watches := []user.User{}
	for _, userID := range us.Watches {
		f, _ := u.FindByID(userID)
		watches = append(watches, f)
	}

	c.RenderArgs["Watches"] = watches

	return c.Render()
}

// Fans 获取所有的粉丝信息
func (c User) Fans() revel.Result {
	// 用户访问自己的所有粉丝
	id := c.Session["UserID"]

	// 其他用户访问
	if len(c.Request.Form["UserID"]) != 0 && c.Request.Form["UserID"][0] != "" {
		id = c.Request.Form["UserID"][0]
	}

	u := user.User{}
	us, err := u.FindByID(id)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	fans := []user.User{}
	for _, userID := range us.Fans {
		f, _ := u.FindByID(userID)
		fans = append(fans, f)
	}

	c.RenderArgs["Fans"] = fans

	return c.Render()
}

// ConversationWith 给某个用户发私信会话
func (c User) ConversationWith(uid string) revel.Result {
	localUserID := c.Session["UserID"]
	remoteUserID := uid

	// 不允许自己给自己发信息
	if localUserID == remoteUserID {
		return c.Redirect(User.Home)
	}

	u := user.User{}
	localUser, err := u.FindByID(localUserID)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	// 检查之前是否有过会话，是则跳转到该会话
	conDB := conversation.NewConversationDB()
	for _, id := range localUser.ConversationIDs {
		conversation, err := conDB.FindByID(id)
		if err != nil {
			log.Println(err)
		}

		if conversation.InitiatorID == uid || conversation.AcceptorID == uid {
			return c.Redirect("Conversation?id=" + id)
		}
	}

	remoteUser, err := u.FindByID(remoteUserID)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	// 创建一个消息内容为空的新会话
	conversation, err := createEmptyConversation(localUser, remoteUser)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	return c.Redirect("Conversation?id=" + conversation.ID)
}

func createEmptyConversation(initiator, acceptor user.User) (conversation.Conversation, error) {
	conDB := conversation.NewConversationDB()
	conver := conversation.Conversation{
		ID:                conversation.CreateObjectID(),
		InitiatorID:       initiator.ID,
		InitiatorNickName: initiator.NickName,
		InitiatorAvatar:   initiator.Avatar,
		AcceptorID:        acceptor.ID,
		AcceptorNickName:  acceptor.NickName,
		AcceptorAvatar:    acceptor.Avatar,
		Messages:          []conversation.Message{},
		TimeStamp:         time.Now().Format("2006-01-02 15:04:05"),
	}

	err := conDB.Add(conver)
	if err != nil {
		return conversation.Conversation{}, err
	}

	err = initiator.UpdateConversationIDs(conver.ID, true)
	if err != nil {
		return conversation.Conversation{}, err
	}

	err = acceptor.UpdateConversationIDs(conver.ID, true)
	if err != nil {
		return conversation.Conversation{}, err
	}

	return conver, nil
}

// Conversation 单个会话页面
func (c User) Conversation(id string) revel.Result {
	conDB := conversation.NewConversationDB()

	conversation, err := conDB.FindByID(id)
	if err != nil {
		log.Println(err)
		return c.RenderText("会话错误")
	}

	LocalUserID, LocalUserAvatar, LocalUserNickName := "", "", ""
	RemoteUserID, RemoteUserAvatar, RemoteUserNickName := "", "", ""

	if c.Session["UserID"] == conversation.InitiatorID {
		LocalUserID = conversation.InitiatorID
		LocalUserAvatar = conversation.InitiatorAvatar
		LocalUserNickName = conversation.InitiatorNickName

		RemoteUserID = conversation.AcceptorID
		RemoteUserAvatar = conversation.AcceptorAvatar
		RemoteUserNickName = conversation.AcceptorNickName
	} else {
		LocalUserID = conversation.AcceptorID
		LocalUserAvatar = conversation.AcceptorAvatar
		LocalUserNickName = conversation.AcceptorNickName

		RemoteUserID = conversation.InitiatorID
		RemoteUserAvatar = conversation.InitiatorAvatar
		RemoteUserNickName = conversation.InitiatorNickName
	}

	conDB.ClearMessageStatus(id)

	c.RenderArgs["Conversation"] = conversation
	c.RenderArgs["UserID"] = c.Session["UserID"]
	c.RenderArgs["RemoteUserNickName"] = RemoteUserNickName

	c.Session["LocalUserID"] = LocalUserID
	c.Session["LocalUserAvatar"] = LocalUserAvatar
	c.Session["LocalUserNickName"] = LocalUserNickName
	c.Session["RemoteUserID"] = RemoteUserID
	c.Session["RemoteUserAvatar"] = RemoteUserAvatar
	c.Session["RemoteUserNickName"] = RemoteUserNickName

	return c.Render()
}

// PostMessage 私信会话中发送单条消息, POST 数据处理
func (c User) PostMessage() revel.Result {
	from := conversation.From{
		UserID:       c.Session["UserID"],
		UserAvatar:   c.Session["LocalUserAvatar"],
		UserNickName: c.Session["LocalUserNickName"],
		Status:       conversation.Read,
	}

	to := conversation.To{
		UserID:       c.Session["RemoteUserID"],
		UserAvatar:   c.Session["RemoteUserAvatar"],
		UserNickName: c.Session["RemoteUserNickName"],
		Status:       conversation.Unread,
	}

	message := conversation.Message{
		ID:        conversation.CreateObjectID(),
		From:      from,
		To:        to,
		Content:   c.Request.Form["Content"][0],
		TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	conDB := conversation.NewConversationDB()
	err := conDB.AddMessage(c.Request.Form["ConversationID"][0], message)
	if err != nil {
		log.Println(err)
		return c.RenderText("消息发送失败")
	}

	response := `<li class="even">
                                    <a class="user" href="{{url "User.Index" ` + from.UserID + `">
                                        <img class="img-responsive avatar_" src="` + from.UserAvatar + `" alt="" style="width:36px">
                                        <span class="user-name">` + from.UserNickName + `</span>
                                    </a>
                                    <div class="reply-content-box">
                                        <span class="reply-time">` + message.TimeStamp + `</span>
                                        <div class="reply-content pr">
                                            <span class="arrow">&nbsp;</span>` + message.Content + `
                                        </div>
                                    </div>
                                </li>`
	return c.RenderText(response)
}

// Conversations 获取所有会话
func (c User) Conversations() revel.Result {
	id := c.Session["UserID"]

	u := user.User{}
	user, err := u.FindByID(id)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	conversations := []conversation.Conversation{}
	conDB := conversation.NewConversationDB()
	for _, id := range user.ConversationIDs {
		log.Println(id)
		convertion, err := conDB.FindByID(id)
		if err != nil {
			log.Println(err)
		} else {
			conversations = append(conversations, convertion)
		}
	}

	c.RenderArgs["Conversations"] = conversations

	return c.Render()
}

// UnreadConversations 获取所有包含有未读消息的会话
func (c User) UnreadConversations() revel.Result {
	id := c.Session["UserID"]

	u := user.User{}
	user, err := u.FindByID(id)
	if err != nil {
		log.Println(err)
		return c.Render(err)
	}

	type ConversationWithUnreadMessageCount struct {
		Conversation       conversation.Conversation
		UnreadMessageCount int
	}
	totalUnreadCount, conversations := 0, []ConversationWithUnreadMessageCount{}
	conDB := conversation.NewConversationDB()

	for _, id := range user.ConversationIDs {
		convertion, err := conDB.FindByID(id)
		if err != nil {
			log.Println(err)
		} else {
			unreadCount := getUnreadMsgCount(id, convertion)
			totalUnreadCount += unreadCount
			if unreadCount != 0 {
				conversations = append(conversations, ConversationWithUnreadMessageCount{convertion, unreadCount})
			}
		}
	}

	c.RenderArgs["Conversations"] = conversations
	c.RenderArgs["TotalUnread"] = totalUnreadCount

	return c.Render()
}

func getUnreadMsgCount(userID string, conver conversation.Conversation) int {
	sum := 0
	for _, message := range conver.Messages {
		if message.To.UserID == userID && message.To.Status == conversation.Unread {
			sum++
		}
	}
	return sum
}

// GetUnreadMsgCount 获取用户未读消息数
func GetUnreadMsgCount(userID string) int {
	u := user.User{}
	user, err := u.FindByID(userID)
	if err != nil {
		log.Println(err)
	}

	unreadCount, conDB := 0, conversation.NewConversationDB()

	for _, id := range user.ConversationIDs {
		convertion, err := conDB.FindByID(id)
		if err != nil {
			log.Println("查找会话" + id + "发生错误：" + err.Error())
		} else {
			unreadCount += getUnreadMsgCount(userID, convertion)
		}
	}
	return unreadCount
}

// GetUnreadMessages 获取会话中未读取的消息
// 每隔一面查询一次数据库数据，一旦发现有未读消息就返回，否则继续循环查询，最大查询次数不超过60次（一分钟）
func (c User) GetUnreadMessages(conversationID string) revel.Result {
	userID := c.Session["UserID"]

	var messages []conversation.Message
	messages = make([]conversation.Message, 0)
	conDB := conversation.NewConversationDB()

	for i := 0; i < 60; i++ {
		conver, err := conDB.FindByID(conversationID)
		if err != nil {
			log.Println("查找会话" + conversationID + "发生错误：" + err.Error())
		}

		for _, message := range conver.Messages {
			if message.To.UserID == userID && message.To.Status == conversation.Unread {
				messages = append(messages, message)
			}

			if len(messages) > 0 {
				c.RenderArgs["Messages"] = messages
				conDB.ClearMessageStatus(conversationID)
				return c.Render()
			}
		}
		time.Sleep(1 * time.Second)
	}

	c.RenderArgs["Messages"] = messages
	conDB.ClearMessageStatus(conversationID)
	return c.Render()
}
