package main

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ctx = context.Background()

type Permission struct {
	ID int
}

func main() {

	dsn := "test:test@tcp(127.0.0.1:3306)/dip?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	var permissions []Permission
	db.Table("permissions").Find(&permissions)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	marshalPermissions, _ := json.Marshal(permissions)

	err := rdb.Set(ctx, "permissions", marshalPermissions, 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "permissions").Result()
	if err != nil {
		panic(err)
	}

	var unmarshalPermissions []Permission
	err = json.Unmarshal([]byte(val), &unmarshalPermissions)
	if err != nil {
		panic(err)
	}
	println(val)
}
