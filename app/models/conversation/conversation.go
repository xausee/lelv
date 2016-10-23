package conversation

import (
	"lelv/app/models/dbmgr"
	"log"

	"gopkg.in/mgo.v2/bson"
)

// 表示信息的读取状态
const (
	Read   = iota // 已读状态
	Unread        //未读状态
)

// From 消息发送方信息
type From struct {
	UserID       string // 发送消息的用户ID
	UserAvatar   string // 发送消息的用户头像地址
	UserNickName string // 发送消息的用户昵称
	Status       int    // 发送消息的用户读取状态 此状态应永远为Read已读
}

// To 消息接受方信息
type To struct {
	UserID       string // 接受消息的用户ID
	UserAvatar   string // 接受消息的用户头像地址
	UserNickName string // 接受消息的用户昵称
	Status       int    // 接受消息的用户读取状态
}

// Message 用户消息
type Message struct {
	ID        string // 消息ID
	From      From   // 发送者信息
	To        To     // 接受这信息
	Content   string // 消息内容
	TimeStamp string // 时间戳
}

// Conversation 会话数据结构
type Conversation struct {
	ID                string    // 会话ID号
	InitiatorID       string    // 发起人ID
	InitiatorNickName string    // 发起人昵称
	InitiatorAvatar   string    // 发起人头像地址
	AcceptorID        string    // 接收人ID
	AcceptorNickName  string    // 接收人昵称
	AcceptorAvatar    string    // 接收人头像地址
	Messages          []Message // 会话中的单条消息
	TimeStamp         string    // 时间戳
}

// ConversationDB 对Conversation数据表的操作struct
type ConversationDB struct {
}

// NewConversationDB 创建对ConversationDB对象
func NewConversationDB() *ConversationDB {
	return &ConversationDB{}
}

// Add 添加会话到数据库
func Add(cvs Conversation) error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Conversations)

	err = c.Insert(cvs)
	if err != nil {
		return err
	}

	return nil
}

// FindByID 根据会话id查找会话
func FindByID(id string) (cvs Conversation, err error) {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Conversations)

	err = c.Find(bson.M{"id": id}).One(&cvs)
	if err != nil {
		return cvs, err
	}

	return cvs, nil
}

// AddMessage 给会话添加一条消息
// cid 会话ID号
func (m *ConversationDB) AddMessage(cid string, msg Message) error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Conversations)
	var (
		old Conversation
		new Conversation
	)

	err = c.Find(bson.M{"id": cid}).One(&old)
	if err != nil {
		log.Println(err)
		return err
	}

	new = old
	new.Messages = append(new.Messages, msg)

	err = c.Update(old, new)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ClearMessageStatus 清除当前会话中所有消息的未读取状态
// cid 会话ID号
func ClearMessageStatus(cid string) error {
	db, err := dbmgr.NewDBManager()
	defer db.Close()

	c := db.Session.DB(dbmgr.Name).C(dbmgr.Conversations)
	var (
		old Conversation
		new Conversation
	)

	err = c.Find(bson.M{"id": cid}).One(&old)
	if err != nil {
		log.Println(err)
		return err
	}

	new = old
	var messages = []Message{}
	// 直接修改new.Messages也会改变old.Messages的值？
	// 而其他属性是不会改变到old的
	// 在此暂时使用new.Messages = append([]Message{}, messages...)的方式
	for _, msg := range new.Messages {
		msg.To.Status = Read
		messages = append(messages, msg)
	}
	new.Messages = append([]Message{}, messages...)

	err = c.Update(old, new)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
