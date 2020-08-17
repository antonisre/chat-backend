package config

import (
	"chat-backend/app/models"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// Insert fake data
func Insert() {
	// Insert two role data
	if err := DB.Debug().Table("roles").Create(&models.Role{Name: "High Admin"}).Error; err != nil {
		panic(err)
	}
	if err := DB.Debug().Table("roles").Create(&models.Role{Name: "Normal Admin"}).Error; err != nil {
		panic(err)
	}

	// Insert user
	if err := DB.Debug().Table("users").Create(&models.User{
		Name:     "High Admin",
		Email:    "highadmin@gmail.com",
		RoleID:   1,
		Password: "password",
		ImageURL: "",
	}).Error; err != nil {
		panic(err)
	}

	if err := DB.Debug().Table("users").Create(&models.User{
		Name:     "Normal Admin",
		Email:    "normaladmin@gmail.com",
		RoleID:   2,
		Password: "password",
	}).Error; err != nil {
		panic(err)
	}
}

// Connect to database
func Connect(DbUser, DbPassword, DbName string) {
	var err error
	DBURI := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbName)

	DB, err = gorm.Open("mysql", DBURI)
	if err != nil {
		fmt.Println("Failed connecting to database")
		panic(err)
	}
	// Migrate the models
	DB.Debug().AutoMigrate(&models.User{}, &models.Role{}, &models.Verification{})

	// Insert fake data
	role := &models.Role{}
	user := &models.UserJSON{}
	countRole, err := role.CountRoles(DB)
	if err != nil {
		panic(err)
	}

	countUser, err := user.CountUsers(DB)
	if err != nil {
		panic(err)
	}

	if countRole == 0 && countUser == 0 {
		Insert()
	}
}
