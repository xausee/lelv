package models

import (
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

// CreateObjectID 创建一个唯一标识Id
func CreateObjectID() string {
	return bson.NewObjectId().Hex()
}

// Add 添加会话到数据库
func (m ConversationDB) Add(conversation Conversation) error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Conversations)

	err = c.Insert(conversation)
	if err != nil {
		return err
	}

	return nil
}

// FindByID 根据会话id查找会话
func (m ConversationDB) FindByID(id string) (con Conversation, err error) {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Conversations)

	err = c.Find(bson.M{"id": id}).One(&con)
	if err != nil {
		return con, err
	}

	return con, nil
}

// AddMessage 给会话添加一条消息
func (m *ConversationDB) AddMessage(conversionID string, msg Message) error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Conversations)
	var (
		oldConver Conversation
		newConver Conversation
	)

	err = c.Find(bson.M{"id": conversionID}).One(&oldConver)
	if err != nil {
		log.Println(err)
		return err
	}

	newConver = oldConver
	newConver.Messages = append(newConver.Messages, msg)

	err = c.Update(oldConver, newConver)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ClearMessageStatus 清除当前会话中所有消息的未读取状态
func (m *ConversationDB) ClearMessageStatus(conversionID string) error {
	db, err := NewDBManager()
	defer db.Close()

	c := db.session.DB(Name).C(Conversations)
	var (
		oldConver Conversation
		newConver Conversation
	)

	err = c.Find(bson.M{"id": conversionID}).One(&oldConver)
	if err != nil {
		log.Println(err)
		return err
	}

	newConver = oldConver
	var messages = []Message{}
	// 直接修改newConver.Messages也会改变oldConver.Messages的值？
	// 而其他属性是不会改变到oldConver的
	// 在此暂时使用newConver.Messages = append([]Message{}, messages...)的方式
	for _, msg := range newConver.Messages {
		msg.To.Status = Read
		messages = append(messages, msg)
	}
	newConver.Messages = append([]Message{}, messages...)

	err = c.Update(oldConver, newConver)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
