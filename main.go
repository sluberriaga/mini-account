package main

import (
	"account/account"
	"fmt"
	"github.com/gin-contrib/static"
	_ "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"time"
)

var loggerMiddleware = gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
})

func main() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(account.TransactionBodyValidation, account.TransactionBody{})
	}

	router := gin.New()
	router.Use(loggerMiddleware)
	router.Use(gin.Recovery())

	mainAccount := account.New()
	accountService := account.NewService(mainAccount)

	router.Use(static.Serve("/", static.LocalFile("./static", true)))
	router.RedirectTrailingSlash = false
	apiRouter := router.Group("/api")
	{
		apiRouter.POST("/account/transaction", accountService.PostTransactionHandler)
		apiRouter.GET("/account/balance", accountService.GetBalance)
		apiRouter.GET("/account/transaction", accountService.GetTransactions)
		apiRouter.GET("/account/transaction/*id", accountService.GetTransactionByID)
	}

	err := router.Run(":8080")
	if err != nil {
		log.Panic(err)
	}
}
