package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-routeros/routeros"
	"net/http"
	"strings"
)

func main() {
	router := gin.Default()
	router.Static("/", "./template")
	router.POST("/login", func(c *gin.Context) {
		address := c.PostForm("address")
		username := c.PostForm("username")
		password := c.PostForm("password")
		command := c.PostForm("command")
		tls := c.PostForm("checkbox")
		async := c.PostForm("async")
		conn, err := dial(address, username, password, tls)
		if err != nil {
			c.String(http.StatusOK, fmt.Sprintf("服务器错误：%s",err))
			return
		}
		defer conn.Close()

		if async == "1" {
			conn.Async()
		}
		r, err := conn.RunArgs(strings.Split(command, " "))
		if err != nil {
			log.Fatal(err)
		}
		log.Print(r)
		c.String(http.StatusOK, fmt.Sprintf("IP地址:%s 用户名:%s 密码:%s 命令:%s", address, username, password, command))

	})
	router.Run(":5555")
}

func dial(address, username, password, tls string) (*routeros.Client, error) {
	if tls == "1" {
		return routeros.DialTLS(address, username, password, nil)
	}
	return routeros.Dial(address, username, password)
}
