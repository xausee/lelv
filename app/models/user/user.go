package user

import (
	"errors"
	"lelv/app/models/dbmgr"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// MockUser 用户数据结构，用于注册
type MockUser struct {
	NickName        string // 昵称
	Email           string // 邮箱
	Gender          string // 性别
	Password        string // 密码
	ConfirmPassword string // 校验密码
}

// Role 用户角色
type Role int

// 用户角色列表
const (
	Super  Role = iota // 超级用户
	Admin              // 管理员
	RgUser             // 普通注册用户
	Guest              // 游客
)

// User 用户数据结构
type User struct {
	ID                  string   // ID号
	NickName            string   // 昵称
	Email               string   // 邮箱
	Gender              string   // 性别
	Password            []byte   // 密码
	Role                Role     // 用户权限类型
	Avatar              string   // 头像地址
	Introduction        string   // 自我介绍
	Fans                []string // 粉丝ID集
	Watches             []string // 关注的用户ID集
	ConversationIDs     []string // 消息ID集合
	Collection          []string // 收藏的博客ID集
	CreateTimeStamp     string   // 创建时间戳
	LastUpdateTimeStamp string   // 最后更新时间戳
}

// SignUp 用户注册
func (m *User) SignUp() error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	i, _ := c.Find(bson.M{"nickname": m.NickName}).Count()
	if i != 0 {
		return errors.New("该账号已被注册")
	}

	err = c.Insert(m)
	if err != nil {
		log.Println("创建账号失败：")
		log.Println(m)
		return err
	}
	return nil
}

// SignIn 用户登录
func SignIn(nickName string, pwd []byte) (u *User, err error) {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	i, err := c.Find(bson.M{"nickname": nickName}).Count()
	if err != nil {
		return nil, err
	}

	if i == 0 {
		return nil, errors.New("用户不存在")
	}

	c.Find(bson.M{"nickname": nickName}).One(&u)

	err = bcrypt.CompareHashAndPassword(u.Password, pwd)
	if err != nil {
		return nil, errors.New("密码不正确")
	}

	return u, nil
}

// Update 根据ID查找用户
func (m *User) Update() error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	var old User
	err = c.Find(bson.M{"id": m.ID}).One(&old)
	if err != nil {
		return err
	}

	new := old
	new.ID = m.ID
	new.NickName = m.NickName
	new.Introduction = m.Introduction
	new.Avatar = m.Avatar
	new.LastUpdateTimeStamp = m.LastUpdateTimeStamp

	err = c.Update(old, new)
	if err != nil {
		return err
	}

	return nil
}

// FindByID 根据ID查找用户
func FindByID(id string) (u User, err error) {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	err = c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return u, err
	}

	return
}

// FindLast 查找最新的n个记录
func (m *User) FindLast(n int) (r []User, err error) {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	type Items map[string]string
	err = c.Find(bson.M{}).Sort("-createtimestamp").Limit(n).All(&r)

	return r, nil
}

// Count 获取所有用户数量
func (m *User) Count() (int, error) {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	n, err := c.Find(bson.M{}).Count()

	if err != nil {
		return 0, err
	}

	return n, err
}

// Collect 收藏博客
// flag 为true：收藏， false：取消收藏
func (m *User) Collect(id string, flag bool) error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	var old User
	err = c.Find(bson.M{"id": m.ID}).One(&old)
	if err != nil {
		return err
	}

	nu := old
	contains := false
	if flag {
		for i := range nu.Collection {
			if nu.Collection[i] == id {
				contains = true
			}
		}
		if !contains {
			nu.Collection = append(nu.Collection, id)
		}
	} else {
		for i := range nu.Collection {
			if nu.Collection[i] == id {
				nu.Collection = append(nu.Collection[:i], nu.Collection[i+1:]...)
			}
		}
	}

	err = c.Update(old, nu)
	if err != nil {
		return err
	}

	return nil
}

// UpdateWatch 关注
// flag 为true：收藏， false：取消收藏
func (m *User) UpdateWatch(id string, flag bool) error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	var old User
	err = c.Find(bson.M{"id": m.ID}).One(&old)
	if err != nil {
		return err
	}

	nu := old
	contains := false
	if flag {
		for _, userid := range nu.Watches {
			if userid == id {
				contains = true
			}
		}
		if !contains {
			nu.Watches = append(nu.Watches, id)
		}
	} else {
		for index, userid := range nu.Watches {
			if userid == id {
				nu.Watches = append(nu.Watches[:index], nu.Watches[index+1:]...)
			}
		}
	}

	err = c.Update(old, nu)
	if err != nil {
		return err
	}

	return nil
}

// UpdateFans 粉丝
// flag 为true：添加粉丝， false：取消粉丝
func (m *User) UpdateFans(id string, flag bool) error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	var old User
	err = c.Find(bson.M{"id": m.ID}).One(&old)
	if err != nil {
		return err
	}

	nu := old
	contains := false
	if flag {
		for _, userid := range nu.Fans {
			if userid == id {
				contains = true
			}
		}
		if !contains {
			nu.Fans = append(nu.Fans, id)
		}
	} else {
		for index, userid := range nu.Fans {
			if userid == id {
				nu.Fans = append(nu.Fans[:index], nu.Fans[index+1:]...)
			}
		}
	}

	err = c.Update(old, nu)
	if err != nil {
		return err
	}

	return nil
}

// UpdateConversationIDs 更新会话ID数据
// flag 为true：添加会话， false：删除会话
func (m *User) UpdateConversationIDs(id string, flag bool) error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Users)

	var old User
	log.Println(m.ID)
	err = c.Find(bson.M{"id": m.ID}).One(&old)
	if err != nil {
		return err
	}

	nu := old
	if flag {
		nu.ConversationIDs = append(nu.ConversationIDs, id)
	} else {
		for index, conversationID := range nu.ConversationIDs {
			if conversationID == id {
				nu.ConversationIDs = append(nu.ConversationIDs[:index], nu.ConversationIDs[index+1:]...)
			}
		}
	}

	err = c.Update(old, nu)
	if err != nil {
		return err
	}

	return nil
}
