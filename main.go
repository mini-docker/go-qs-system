package main

import (
	"../go-qs-system/controller/account"
	"../go-qs-system/dal/db"
	"../go-qs-system/filter"
	"../go-qs-system/id_gen"
	"../go-qs-system/logger"
	maccount "../go-qs-system/middleware/account"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
)

func initTemplate(router *gin.Engine) {
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.Static("/css/", "./static/css/")
	router.Static("/fonts/", "./static/fonts/")
	router.Static("/img/", "./static/img/")
	router.Static("/js/", "./static/js/")
}

func  initDb() (err error) {
	dns := "root:root@tcp(localhost:33061)/mercury?parseTime=true"
	err = db.Init(dns)
	if err != nil {
		return
	}
	return
}

func initSession() (err error) {
	//err = maccount.InitSession("redis", "localhost:6379")
	err = maccount.InitSession("memory", "")
	return
}

func initFilter() (err error) {

	err = filter.Init("./data/filter.dat.txt")
	if err != nil {
		logger.Error("init filter failed, err:%v", err)
		return
	}

	logger.Debug("init filter succ")
	return
}


func main()  {
	router := gin.Default()
	config := make(map[string]string)
	config["log_level"] = "debug"
	logger.InitLogger("console", config)

	err := initDb()
	if err != nil {
		panic(err)
	}

	err = initFilter()
	if err != nil {
		panic(err)
	}
	err = initSession()
	if err != nil {
		panic(err)
	}

	err = id_gen.Init(1)
	if err != nil {
		panic(err)
	}

	ginpprof.Wrapper(router)
	initTemplate(router)
	router.POST("/api/user/register", account.RegisterHandle)
	router.POST("/api/user/login", account.LoginHandle)




	router.Run(":9090")
}