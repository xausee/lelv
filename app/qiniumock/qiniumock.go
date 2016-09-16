package qiniumock

import (
	"encoding/base64"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

const (
	// BUCKET 七牛云空间
	BUCKET = "lelvmock"
	// SPACE 云空间对外url空间
	SPACE = "http://odkkucuql.bkt.clouddn.com/"
	// CDN 加速绑定域名
	//CDN = "http://file.lelvboke.com/"
	// AK access key
	AK = "ZJ-KESjY89JYeZGT3dCvgQIngru-Qnkzt9PScvH1"
	// SK secret key
	SK = "gRmeHtZ20UhTsUh0AfnlMHcUstwK7fR3vW-BSQQe"
)

const (
	// DefaultMaleAvatar 缺省男生头像
	// 头像一：C7399E2CAB400001631110E0CD5A1AD6.jpg
	// 头像二：C7399E2E77500001C1B51C80917815EE.jpg
	// 头像三：C7399E30312000014AFA169A1F70166B.jpg
	DefaultMaleAvatar = SPACE + "C7399E2CAB400001631110E0CD5A1AD6.jpg"
	// DefaultFemaleAvatar 缺省女生头像
	// 头像一：C7399E25C400000163D1A46082D0AAA0.jpg
	// 头像二：C7399E28B1000001C05A7FA01C70DB90.jpg
	DefaultFemaleAvatar = SPACE + "C7399E25C400000163D1A46082D0AAA0.jpg"
)

// CreatUpToken 创建上传凭证
func CreatUpToken() string {
	//初始化AK，SK
	bucket := BUCKET
	conf.ACCESS_KEY = AK
	conf.SECRET_KEY = SK

	// //创建一个Client
	c := kodo.New(0, nil)

	// //设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: bucket,
		//设置Token过期时间
		Expires: 3600,
	}
	//生成一个上传token
	token := c.MakeUptoken(policy)

	return token
}

// Upload 将本地文件上传到七牛云存储，返回文件在七牛上的路径。 fp 是本地文件路径， key 是七牛上文件名
func Upload(fp, key string) error {
	kodo.SetMac(AK, SK)

	zone := 0                // 空间(Bucket)所在的区域
	c := kodo.New(zone, nil) // 用默认配置创建 Client
	b := c.Bucket(BUCKET)
	ctx := context.Background()

	err := b.PutFile(ctx, nil, key, fp, nil)
	if err != nil {
		return err
	}
	return nil
}

// DecodeBase64 解码Base64 图片
func DecodeBase64(f, d string) error {
	b, _ := base64.StdEncoding.DecodeString(d) //成图片文件并把文件写入到buffer
	err := ioutil.WriteFile(f, b, 0666)        //buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除七牛上的文件， 可以是文件名
func Delete(key string) {
	conf.ACCESS_KEY = AK
	conf.SECRET_KEY = SK

	//new一个Bucket管理对象
	c := kodo.New(0, nil)
	p := c.Bucket(BUCKET)

	//调用Delete方法删除文件
	res := p.Delete(nil, key)
	//打印返回值以及出错信息
	if res == nil {
		log.Println("Delete success")
	} else {
		log.Println(res)
	}
}

// Move 移动七牛上的文件
// keySrc  是要移动的文件的旧路径。
// keyDest 是要移动的文件的新路径。
func Move(keySrc, keyDest string) error {
	conf.ACCESS_KEY = AK
	conf.SECRET_KEY = SK

	//new一个Bucket管理对象
	c := kodo.New(0, nil)
	p := c.Bucket(BUCKET)

	//调用Move方法移动文件
	err := p.Move(nil, keySrc, keyDest)
	if err != nil {
		return err
	}

	return nil
}
