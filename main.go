package main

import (
	// "fmt"
	"example.com/project1/db"
	"example.com/project1/handlers"
	"example.com/project1/logins"
	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	r := gin.Default()
	r.Static("static", "./static")
	//admin roles
	r.POST("/add-book", handlers.AddBook)
	r.DELETE("/remove-book", handlers.RemoveBook)
	r.PATCH("/update-book", handlers.UpdateBook)
	r.GET("/list-issue-request", handlers.ListIssueRequests)
	r.POST("/libraries", handlers.CreateLibrary)

	//readers role
	r.POST("/search/book", handlers.SearchBook)
	r.POST("/issue/request", handlers.RaiseIssueRequest)
	r.PUT("/approve-issue-request/:request_id", handlers.ApproveIssueRequest)
	r.DELETE("/reject-issue-request/:request_id", handlers.RejectIssueRequest)

	//starting page
	r.GET("/", handlers.StartingPage)

	//authentication endpoints
	r.POST("/login", handlers.Login)
	r.GET("/home", handlers.Home)

	//login with otp
	r.POST("/send-otp", logins.GenerateAndSendOTP)
	r.POST("/validate-otp/:otp", logins.ValidateOTP)
	r.Run(":8088")

}
