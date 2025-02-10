package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type Permission struct {
	ID int
}

func (u *Permission) AfterUpdate(tx *gorm.DB) (err error) {
	var permissions []Permission
	tx.Table("redis_permissions").Find(&permissions)
	marshalPermissions, err := json.Marshal(permissions)
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "permissions", marshalPermissions, 0).Err()
	if err != nil {
		panic(err)
	}
	return
}

func (u *Permission) AfterCreate(tx *gorm.DB) (err error) {
	var permissions []Permission
	tx.Table("redis_permissions").Find(&permissions)
	marshalPermissions, err := json.Marshal(permissions)
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "permissions", marshalPermissions, 0).Err()
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	dsn := "test:test@tcp(127.0.0.1:3306)/redis_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var permissions []Permission
	db.Table("redis_permissions").Find(&permissions)
	marshalPermissions, err := json.Marshal(permissions)
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "permissions", marshalPermissions, 0).Err()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/permissions", getPermissions)
	r.GET("/permissions-mysql", getPermissionsMySQL)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getPermissions(c *gin.Context) {
	val, err := rdb.Get(ctx, "permissions").Result()
	if err != nil {
		panic(err)
	}
	var unmarshalPermissions []Permission
	err = json.Unmarshal([]byte(val), &unmarshalPermissions)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, unmarshalPermissions)
}

func getPermissionsMySQL(c *gin.Context) {
	dsn := "test:test@tcp(127.0.0.1:3306)/redis_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var permissions []Permission
	db.Table("redis_permissions").Find(&permissions)
	c.JSON(http.StatusOK, permissions)
}
