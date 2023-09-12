package repo

import (
	"log"

	"github.com/Coreychen4444/shortvideo_ms-video/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbRepository struct {
	db *gorm.DB
}

func NewDbRepository(db *gorm.DB) *DbRepository {
	return &DbRepository{db: db}
}

// mysql 初始化
func InitMysql() *gorm.DB {
	// 连接数据库(用户名和密码自己改)
	dsn := "root:44447777@tcp(127.0.0.1:3306)/tiktok_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error() + ", failed to connect database")
	}
	// 自动迁移
	err = db.AutoMigrate(&model.Video{}, &model.VideoLike{}, &model.Comment{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}
	log.Println("成功连接mysql数据库!")
	return db
}
