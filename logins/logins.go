package logins

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"example.com/project1/db"
	"example.com/project1/models"
	"github.com/gin-gonic/gin"
)

//struct defination

type Credential struct {
	Email string `json:"email"`
}

// generating otp

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(999999))
}

var otp string
var email string

//sending otp

func SendOTP(email, otp string) error {
	from := "example@gmail.com"
	password := "password_detailss"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := "From:" + from + "\n" +
		"To:" + email + "\n" +
		"Subject: One time OTP for login\n\n" +
		"Your OTP is: " + otp

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg))
	if err != nil {
		return err
	}
	fmt.Println("Email sent successfully : " + email)
	return nil

}

// generate and sendotp
// var otp string

func GenerateAndSendOTP(c *gin.Context) {
	var request Credential
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var user models.User
	result := db.DB.Where("email = ?", request.Email).First(&user)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no record found"})
		return
	}
	email = request.Email

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	otp = GenerateOTP()
	err := SendOTP(email, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"otp": otp})
}
func ValidateOTP(c *gin.Context) {
	submittedOTP := otp
	serverOTP := c.Param("otp")

	if submittedOTP == serverOTP {
		from := "example@gmail.com"
		password := "password_detailss"

		smtpHost := "smtp.gmail.com"
		smtpPort := "587"

		msg := "From:" + from + "\n" +
			"To:" + email + "\n" +
			"Subject: Login Successful\n\n" +
			"Congratulations you have successfully logged in !!!"

		auth := smtp.PlainAuth("", from, password, smtpHost)

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "sorry no confirmation message delievered"})
		}
		c.JSON(200, gin.H{"message": "OTP is valid"})
	} else {
		from := "example@gmail.com"
		password := "password_detailss"

		smtpHost := "smtp.gmail.com"
		smtpPort := "587"

		msg := "From:" + from + "\n" +
			"To:" + email + "\n" +
			"Subject:  ****Login Failure ****\n\n" +
			"Sorry you are not authorized user  !!!"

		auth := smtp.PlainAuth("", from, password, smtpHost)

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "sorry no confirmation message delievered"})
		}
		c.JSON(400, gin.H{"error": "Invalid OTP"})
	}
}
