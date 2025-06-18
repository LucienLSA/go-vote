package session

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
)

var store = sessions.NewCookieStore([]byte("lucien-go-vote"))
var sessionName = "session-name"

func GetSession(c *gin.Context) map[interface{}]interface{} {
	session, _ := store.Get(c.Request, sessionName)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

func SetSession(c *gin.Context, name string, id int64) error {
	session, _ := store.Get(c.Request, sessionName)
	session.Values["name"] = name
	session.Values["id"] = id
	return session.Save(c.Request, c.Writer)
}

func FlushSession(c *gin.Context) error {
	session, _ := store.Get(c.Request, sessionName)
	fmt.Printf("session : %+v\n", session.Values)
	session.Values["name"] = ""
	session.Values["id"] = int64(0)
	return session.Save(c.Request, c.Writer)
}

var SessionStore *redisstore.RedisStore
var sessionNameV1 = "session-name-v1"

func GetSessionV1(c *gin.Context) map[interface{}]interface{} {
	session, _ := SessionStore.Get(c.Request, sessionNameV1)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

func SetSessionV1(c *gin.Context, name string, id int64) error {
	session, _ := SessionStore.Get(c.Request, sessionNameV1)
	session.Values["name"] = name
	session.Values["id"] = id
	return session.Save(c.Request, c.Writer)
}

func FlushSessionV1(c *gin.Context) error {
	session, _ := SessionStore.Get(c.Request, sessionNameV1)
	fmt.Printf("session : %+v\n", session.Values)
	session.Values["name"] = ""
	session.Values["id"] = ""
	return session.Save(c.Request, c.Writer)
}
