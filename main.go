package main

import (
	"../go-qs-system/controller/account"
	"../go-qs-system/controller/answer"
	"../go-qs-system/controller/category"
	"../go-qs-system/controller/comment"
	"../go-qs-system/controller/favorite"
	"../go-qs-system/controller/question"
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
	router.GET("/api/category/list", category.GetCategoryListHandle)
	router.POST("/api/ask/submit", maccount.AuthMiddleware, question.QuestionSubmitHandle)
	router.GET("/api/question/list", category.GetQuestionListHandle)
	router.GET("/api/question/detail", question.QuestionDetailHandle)
	router.GET("/api/answer/list", answer.AnswerListHandle)

	//评论模块
	//commentGroup := router.Group("/api/comment/", maccount.AuthMiddleware)
	commentGroup := router.Group("/api/comment/")
	commentGroup.POST("/post_comment", comment.PostCommentHandle)
	commentGroup.POST("/post_reply", comment.PostReplyHandle)
	commentGroup.GET("/list", comment.CommentListHandle)
	commentGroup.GET("/reply_list", comment.ReplyListHandle)
	commentGroup.POST("/like", comment.LikeHandle)

	//收藏模块路由
	favoriteGroup := router.Group("/api/favorite/")
	favoriteGroup.POST("/add_dir", favorite.AddDirHandle)
	favoriteGroup.POST("/add", favorite.AddFavoriteHandle)
	favoriteGroup.GET("/dir_list", favorite.DirListHandle)
	favoriteGroup.GET("/list", favorite.FavoriteListHandle)


	router.Run(":9090")
}