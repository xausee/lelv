package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/mgo.v2/bson"
)

// CreateObjectID 创建一个唯一标识Id
func CreateObjectID() string {
	return bson.NewObjectId().Hex()
}

// GetAndSaveHTML 获取网页内容并保存到指定文件
// fn文件路径（文件名）
func GetAndSaveHTML(url, fn string) error {
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	body, _ := ioutil.ReadAll(res.Body) //转换byte数组
	//io.Copy(os.Stdout, res.Body)//写到输出流，
	bodystr := string(body) //转换字符串

	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	absPath, _ := filepath.Abs(pwd + fn)

	file1, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, os.ModeType)
	defer file1.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	// 往创建的文件中写入内容
	_, err = file1.WriteString(bodystr)
	log.Println(bodystr)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
