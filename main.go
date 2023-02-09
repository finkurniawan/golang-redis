package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
	"redis-golang/db"
)

func main() {
	//initialize redis
	db.RedisInit()

	//initialize server
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/insert", Insert)
	e.GET("/get", Get)

	e.Logger.Fatal(e.Start(":1323"))
}

type RespJson struct {
	Data   interface{}
	Status string
}

type RequestRedis struct {
	Name string
	Age  string
}

var key string = "app_key"

func Insert(c echo.Context) error {
	id := c.QueryParam("id")
	name := c.QueryParam("name")
	age := c.QueryParam("age")

	//redis connection
	rdb := db.RedisConnect()

	reqRedis := RequestRedis{
		Name: name,
		Age:  age,
	}

	req, _ := json.Marshal(reqRedis)

	err := rdb.HSet(key, id, req).Err()

	if err != nil {
		return fmt.Errorf("error set redis %s", err)
	}

	//response

	resp := RespJson{
		Data:   id,
		Status: "succes",
	}

	return c.JSON(http.StatusOK, resp)
}

func Get(c echo.Context) error {
	id := c.QueryParam("id")

	//connection
	rdb := db.RedisConnect()

	val, err := rdb.HGet(key, id).Result()

	if err != redis.Nil {
		return c.JSON(http.StatusNotFound, "Data not found")
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("error get redis %s", err.Error()))
	}

	var requestRedis RequestRedis

	err = json.Unmarshal([]byte(val), &requestRedis)

	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("error unmarshal redis %s", err.Error()))
	}

	resp := RespJson{
		Data:   requestRedis,
		Status: "succes",
	}
	return c.JSON(http.StatusOK, resp)
}
