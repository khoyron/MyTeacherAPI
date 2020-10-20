package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"myTeacherEndPoint/model"
)

// DBInit create connection to database
//127.0.0.1
func DBInit(table string) *gorm.DB {
	db, err := gorm.Open("mysql", "root:pengusaha@tcp(localhost:3306)/myteacher?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Printf("conn fail: %s\n", err)
		//panic("failed to connect to database")
	}


	if table=="user" {
		db.AutoMigrate(model.User{})
		return db
	}else if table=="category" {
		db.AutoMigrate(model.Category{})
		return db
	}else if table=="services"{
		db.AutoMigrate(model.Services{})
		return db
	}else if table=="media"{
		db.AutoMigrate(model.Media{})
		return db
	}else if table=="review"{
		db.AutoMigrate(model.Review{})
		return db
	}else if table=="transaction"{
		db.AutoMigrate(model.Transaction{})
		return db
	}else if table== "schedule"{
		db.AutoMigrate(model.Schedule{})
		return db
	}else if table=="clasess"{
		db.AutoMigrate(model.Classes{})
		return db
	}else if table== "receipt"{
		db.AutoMigrate(model.Receipt{})
		return db
	}else if table== "city"{
		db.AutoMigrate(model.City{})
		return db
	}else if table== "profincy"{
		db.AutoMigrate(model.Profincy{})
		return db
	}else if table== "bank"{
		db.AutoMigrate(model.BankDetail{})
		return db
	}else {
		return db
	}
}

//db, err := gorm.Open("mysql", "myteacher:pengusaha@tcp(myteacher.cvjxzuxfdqk4.us-east-2.rds.amazonaws.com:3306)/myteacher?charset=utf8&parseTime=True&loc=Local")
