package main

import (
	"strings"
	"github.com/spf13/viper"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"net/http"

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

type todo struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
	Topic string `json:"topic" bson:"topic"`
	Done bool `json:"done" bson:"done"`
}


type handler struct {
	m *mgo.Session
}

func (h *handler) delete(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	col := session.DB("workshop").C("todos")
	if err := col.RemoveId(id); err != nil {
		return err
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"result": "success",
	})
}

func (h *handler) edit(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	var ot todo
	col := session.DB("workshop").C("todos")
	if err := col.FindId(id).One(&ot); err != nil {
		return err
	}
	ot.Topic = t.Topic
	ot.Done = t.Done

	if err := col.UpdateId(id, ot); err != nil {
		return err
	}
	t.ID = ot.ID

	return c.JSON(http.StatusOK, t)
}

func (h *handler) view(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	var t todo
	col := session.DB("workshop").C("todos")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}


	return c.JSON(http.StatusOK, t)
}

func (h *handler) list(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var ts []todo
	col := session.DB("workshop").C("todos")
	if err := col.Find(nil).All(&ts); err != nil {
		return err
	}


	return c.JSON(http.StatusOK, ts)
}

func (h *handler) create(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	t.ID = bson.NewObjectId()
	
	col := session.DB("workshop").C("todos")
	if err := col.Insert(t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}