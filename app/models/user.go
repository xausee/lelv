package models

import (
	"errors"
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

// User 用户数据结构
type User struct {
	ID                  string   // ID号
	NickName            string   // 昵称
	Email               string   // 邮箱
	Gender              string   // 性别
	Password            []byte   // 密码
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
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)
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
func (m *User) SignIn() error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)
	i, _ := c.Find(bson.M{"nickname": m.NickName}).Count()
	if i == 0 {
		err = errors.New("用户不存在")
		return err
	}

	var u *User
	c.Find(bson.M{"nickname": m.NickName}).One(&u)
	err = bcrypt.CompareHashAndPassword(u.Password, m.Password)
	if err != nil {
		err = errors.New("密码不正确")
	}

	m.ID = u.ID
	m.NickName = u.NickName
	m.Email = u.Email
	m.Gender = u.Gender
	m.Password = u.Password
	m.Avatar = u.Avatar
	m.Introduction = u.Introduction

	return nil
}

// Update 根据ID查找用户
func (m *User) Update() error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)

	var oldUser User
	err = c.Find(bson.M{"id": m.ID}).One(&oldUser)
	if err != nil {
		return err
	}

	newuser := oldUser
	newuser.ID = m.ID
	newuser.NickName = m.NickName
	newuser.Introduction = m.Introduction
	newuser.Avatar = m.Avatar
	newuser.LastUpdateTimeStamp = m.LastUpdateTimeStamp

	err = c.Update(oldUser, newuser)
	if err != nil {
		return err
	}

	return nil
}

// FindByID 根据ID查找用户
func (m *User) FindByID(id string) (u User, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)
	err = c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return u, err
	}

	return
}

// FindLast 查找最新的n个记录
func (m *User) FindLast(n int) (r []User, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)
	type Items map[string]string
	err = c.Find(bson.M{}).Sort("-createtimestamp").Limit(n).All(&r)

	return r, nil
}

// Count 获取所有用户数量
func (m *User) Count() (int, error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)
	n, err := c.Find(bson.M{}).Count()

	if err != nil {
		return 0, err
	}

	return n, err
}

// Collect 收藏博客
// flag 为true：收藏， false：取消收藏
func (m *User) Collect(id string, flag bool) error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)

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
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)

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
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)

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
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Users)

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
