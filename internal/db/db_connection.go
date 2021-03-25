package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var W_Db *gorm.DB

//
//type DBServer struct {
//	Host string `toml:"host"`
//	Port int `toml:"port"`
//	Dbname string `toml:"dbname"`
//	User string `toml:"user"`
//	Password string `toml:"password"`
//}
//
//
//func (m DBServer)ConnectString() string{
//	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
//		m.Host,m.Port,m.User,m.Dbname)
//}
//
//
////初始化gorm的连接
//func (m DBServer) NewGormDB(openConnection int) (db *gorm.DB,err error){
//	db,err=gorm.Open("postgres",m.ConnectString())
//	if err != nil {
//		return
//	}
//	//设置最大连接数
//	db.DB().SetMaxOpenConns(openConnection)
//	return db,err
//}

func InitDbConnection(dbKey string) {
	host := viper.GetString(dbKey + ".host")
	user := viper.GetString(dbKey + ".user")
	dbname := viper.GetString(dbKey + ".name")
	password := viper.GetString(dbKey + ".password")
	url := "host=" + host + " user=" + user + " dbname=" + dbname + " sslmode=disable password=" + password + ""
	const postgres = "postgres"
	var err error
	W_Db, err = gorm.Open(postgres, url)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("连接成功")
		//fmt.Println(db.HasTable(&User{}))
	}
}

////初始化redis数据库的配置
//type RedisServer struct {
//	Addr string `toml:"addr"`
//	Password string `toml:"password"`
//	DB int `toml:"db"`
//}
//
////初始化redis连接池
//func (c RedisServer) NewRedisPool(maxIdle int) *redis.Pool {
//	return &redis.Pool{
//		MaxIdle:     maxIdle,
//		IdleTimeout: 240 * time.Second,
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.Dial("tcp", c.Addr, redis.DialDatabase(c.DB), redis.DialPassword(c.Password))
//			if err != nil {
//				return nil, err
//			}
//			return c, err
//		},
//		TestOnBorrow: func(c redis.Conn, t time.Time) error {
//			_, err := c.Do("PING")
//			return err
//		},
//	}
//}
//
//
//
//
//
//
