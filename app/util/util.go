package util

import "gopkg.in/mgo.v2/bson"

// CreateObjectID 创建一个唯一标识Id
func CreateObjectID() string {
	return bson.NewObjectId().Hex()
}
