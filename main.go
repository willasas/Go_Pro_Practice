package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

// 定义一个model
type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(110);not null;unique"`
	Password string `gorm:"size:255;not null"`
}

func main(){
	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/API/auth/register", func(c *gin.Context) {
		// 获取参数
		name := c.PostForm("name")
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")
		// 数据验证
		if len(telephone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}
		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}

		// 如果名称没有传，给一个10位的随机字符串
		if len(name) == 0 {
			name = RandomString(10)
		}
		log.Println(name, telephone, password)

		// 判断手机号是否存在
		if isTelephoneExist(db, telephone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
			return
		}

		// 创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		// 返回结果
		c.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	panic(r.Run())	//listen and serve on 0.0.0.0:8080
}

// 验证手机号函数
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

// 创建一个随机的字符串
func RandomString(n int) string {
	var letters  = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

// 通过GORM连接数据库
func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "ROOT"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	// 自动创建数据表
	db.AutoMigrate(&User{})
	return db
}