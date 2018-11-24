package main

import (
	"strings"
	"github.com/spf13/viper"
	"fmt"

	"gopkg.in/mgo.v2"

	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
)

func main() {
	
	viper.SetDefault("port", "1323")
	viper.SetDefault("mongo.host", "localhost:27017")
	viper.SetDefault("mongo.user", "root")
	viper.SetDefault("mongo.pass", "example")
	viper.AutomaticEnv()
	// yml
	// mongo:
	//   host:
	// env
	// MONGO_HOST
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	mongoHost := viper.GetString("mongo.host")//"localhost:27017"
	mongoUser := viper.GetString("mongo.user")//"root"
	mongoPass := viper.GetString("mongo.pass")//"example"
	port := viper.GetString("port")//"1323"

	connString := fmt.Sprintf("%s:%s@%s", mongoUser, mongoPass, mongoHost)
	session, err := mgo.Dial(connString)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	h := &handler{
		m: session,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.POST("/todos", h.create)
	e.GET("/todos", h.list)
	e.GET("/todos/:id", h.view)
	e.PUT("/todos/:id", h.edit)
	e.DELETE("/todos/:id", h.delete)
	e.Logger.Fatal(e.Start(":" + port))
}