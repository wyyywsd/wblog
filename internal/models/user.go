package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
	"log"
	"time"
)

type User struct {
	gorm.Model
	Username             string `gorm:"not null;unique"`
	PassWord             string `gorm:"not null"`
	Email                string `gorm:"unique_index;default:null"`
	ProfilePhoto         string `gorm:"default:'/file/default.jpeg'"`
	UserRegistrationTime time.Time
	UserBirthday         time.Time
	UserAge              uint
	UserTelephoneNumber  string `gorm:"unique_index;default:null"`
	UserNickname         string
	IsAdmin              bool `gorm:"default:false"`
	GithubUrl            string
	Articles             []*Article    `gorm:"ForeignKey:UserID;"`
	Comments             []*Comment    `或者cgorm:"ForeignKey:UserID;"`
	UserFriends          []*UserFriend `gorm:"ForeignKey:UserId"`
}

func UpdateUser(user User, updateMap map[string]interface{}) error {
	fmt.Println(user.Username)
	fmt.Println(updateMap)
	err := db.W_Db.Model(&user).Updates(updateMap).Error
	return err

}
func FindUserByUserName(name string) (*User, bool, error) {
	log.Println("000000000")
	log.Println(name)
	var err error
	user := &User{}
	//err = db.W_Db.Where("username = ?", name).Find(&user).Error
	err = db.W_Db.Table("users").Where("username = ?", name).Find(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			//如果并不是找不到用户的操作
			return nil, false, err
		}
		//如果是找不到的问题，直接返回
		return nil, false, nil
	}
	return user, true, err
}

func FindUserById(id string) User {
	user := User{}
	db.W_Db.Where("id = ?", id).First(&user)
	return user
}

func FindUserByEmail(email string) (*User, bool) {
	user := &User{}
	var err error
	err = db.W_Db.Table("users").Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, false
	}
	return user, true
}

func CreateUser(username string, password string, email string) (*User, error) {
	user := &User{Username: username, PassWord: password, Email: email, ProfilePhoto: "/file/default.jpg"}
	err := db.W_Db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
func FindUserByArticle(article *Article) (*User, bool) {

	var err error
	user := &User{}
	//err = db.W_Db.Where("username = ?", name).Find(&user).Error
	err = db.W_Db.Table("users").Where("id = ?", article.UserId).Find(&user).Error
	if err != nil {
		return nil, false
	}
	return user, true
}

//func AddUser(user User)(User,error) {
//	user = User{UserName: user.UserName,PassWord: user.PassWord}
//	err := db.W_Db.Create(&user).Error
//	if err != nil {
//		return _,err
//	}
//	return user,nil
//}
