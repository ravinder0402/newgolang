package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"example.com/project1/db"
	"example.com/project1/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// for starting page on web
func StartingPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello your Website is running"})

}

func AddBook(c *gin.Context) {
	var request struct {
		Book       models.BookInventory `json:"book"`
		AdminEmail string               `json:"email"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//checking for email of admin whether entered email matches with admin email or not
	var existingEmail models.User
	result1 := db.DB.Where("email = ?", request.AdminEmail).First(&existingEmail)
	if result1.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, "Admin email is incorrect")
		return
	}

	//checking presence of book

	var existingBook models.BookInventory
	result := db.DB.Where("isbn = ?", request.Book.ISBN).First(&existingBook)
	if result.RowsAffected > 0 {
		existingBook.TotalCopies++
		if err := db.DB.Save(&existingBook).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update book"})
			return
		}

		if err := sendTeamsNotification("https://xenonstack1.webhook.office.com/webhookb2/98c349a5-2f06-4612-945e-ad6ab51b9667@7ff914bc-ca07-4c28-8277-73e20a4966c7/IncomingWebhook/cb6509e4206748eab2abdb68bf300679/45af8d51-3b65-4dc2-a683-972d57f8c8a0", "Book Added", "A copy ofexisting book has been added", request.Book.Title, request.Book.Authors); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification on team"})
			return
		}
		c.JSON(http.StatusOK, existingBook)
		return
	}

	//book is not present

	request.Book.TotalCopies = 1
	if err := db.DB.Create(&request.Book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}
	c.JSON(http.StatusCreated, request.Book)
	if err := sendTeamsNotification("https://xenonstack1.webhook.office.com/webhookb2/98c349a5-2f06-4612-945e-ad6ab51b9667@7ff914bc-ca07-4c28-8277-73e20a4966c7/IncomingWebhook/cb6509e4206748eab2abdb68bf300679/45af8d51-3b65-4dc2-a683-972d57f8c8a0", "Book Added", "A new book has been added", request.Book.Title, request.Book.Authors); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification on team"})
		return
	}

}

// function to send notification on email
func sendTeamsNotification(webhookURL, title, subtitle, bookTitle, bookAuthor string) error {
	teamsPayload := map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "http://schema.org/extensions",
		"themeColor": "0076D7",
		"summary":    title,
		"sections": []map[string]string{
			{"activityTitle": title,
				"activitySubtitle": subtitle,
				"text":             "Title:" + bookTitle + "\n Author:" + bookAuthor,
			},
		},
	}
	//converting payload to josn format

	payloadJSON, err := json.Marshal(teamsPayload)
	if err != nil {
		return err
	}

	// final statement to send teams message
	_, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}
	return nil
}

//removing book

func RemoveBook(c *gin.Context) {
	var request struct {
		ISBN string `json:"isbn"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	//check whether book is present or not

	var book models.BookInventory
	result := db.DB.Where("isbn = ?", request.ISBN).First(&book)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no book found"})
		return
	}
	//checking whether no of toalcopies < 0
	if book.TotalCopies == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No copies availabe to delete"})
		return
	}

	// check if any copies issued

	if book.AvailableCopies < book.TotalCopies {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove copy"})
		return
	}

	// decrease number of copies

	book.TotalCopies--

	if err := db.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book removed successfully"})

}

//Update book

func UpdateBook(c *gin.Context) {
	var request struct {
		ISBN           string               `json:"isbn"`
		UpdatedDetails models.BookInventory `json:"updated_details"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//book exist?

	var book models.BookInventory
	result := db.DB.Where("isbn = ?", request.ISBN).First(&book)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no specific isbn book present"})
		return
	}

	//update details

	if err := db.DB.Model(&book).Updates(request.UpdatedDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update details in book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

//Listing issue request

func ListIssueRequests(c *gin.Context) {
	var issueRequests []models.RequestEvent

	//fetch request from database

	if err := db.DB.Find(&issueRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch issue requests"})
		return
	}
	c.JSON(http.StatusOK, issueRequests)
}

//request approval

func ApproveIssueRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	var request models.RequestEvent
	if err := db.DB.First(&request, requestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	//Update request details

	request.ApprovalDate = time.Now()
	request.ApproverID = 1
	if err := db.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve request"})
		return
	}
	var issueRegisty models.IssueRegistry
	//putting values in issue registry
	result := db.DB.Where("reader_id = ?", requestID).First(&issueRegisty)
	fmt.Println(requestID)
	if result.RowsAffected == 0 {
		issueRegistry := models.IssueRegistry{
			ISBN:               request.BookID,
			ReaderID:           request.ReaderID,
			IssueApproverID:    request.ApproverID,
			IssueStatus:        "APPROVED",
			IssueDate:          time.Now(),
			ExpectedReturnDate: time.Now(),
			ReturnDate:         time.Now(),
			ReturnApproverID:   request.ApproverID,
		}
		if err := db.DB.Create(&issueRegistry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create issue request"})
			return
		}

	} else {
		if err := db.DB.Where("reader_id = ?", request.ReaderID).First(&issueRegisty).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find issue registory"})
			return
		}

		issueRegisty.IssueStatus = "Approved"
		issueRegisty.IssueDate = time.Now()

		if err := db.DB.Save(&issueRegisty).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update issue registry"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue request approved successfully"})

}

// rejectissueregistry

func RejectIssueRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errror": "Invalid request id"})
		return
	}

	var request models.RequestEvent
	if err := db.DB.First(&request, requestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	// Delete request
	if err := db.DB.Delete(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "issue request rejected successfully"})
}

// creating a library
var num uint

func CreateLibrary(c *gin.Context) {
	var request struct {
		LibraryName string `json:"library_name"`
		OwnerName   string `json:"owner_name"`
		OwnerEmail  string `json:"owner_email"`
		OwnerPhone  string `json:"owner_phone"`
		OwnerRole   string `json:"owner_role"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if library already exists

	var existingLibrary models.Library

	result := db.DB.Where("name = ?", request.LibraryName).First(&existingLibrary)
	if request.OwnerRole == "Admin" {
		if result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Library name already exist"})
			return
		}
	}
	//creating new library
	newLibrary := models.Library{
		Name: request.LibraryName,
	}
	if request.OwnerRole != "Reader" {

		if err := db.DB.Create(&newLibrary).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create librarry"})
			return
		}
	}
	//fetching library id from library table

	if request.OwnerRole == "Reader" {
		if result.RowsAffected > 0 {
			num = existingLibrary.ID

		}

	} else {
		num = newLibrary.ID
	}

	//create the owner user
	var mod models.User
	newOwner := models.User{
		Name:          request.OwnerName,
		Email:         request.OwnerEmail,
		ContactNumber: request.OwnerPhone,
		Role:          request.OwnerRole,
		LibID:         num,
	}
	if err := db.DB.Create(&newOwner).Error; err != nil {
		mod.ID--
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create owner"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Library created successfully", "library_id": num})

}

//search for books matching

func SearchBook(c *gin.Context) {
	var request struct {
		Title     string `json:"title"`
		Author    string `json:"author"`
		Publisher string `json:"publisher"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// search for boks matching

	var books models.BookInventory
	result := db.DB.Where("title LIKE ?", "%"+request.Title+"%").
		Where("authors LIKE ?", "%"+request.Author+"%").
		Where("publisher LIKE ?", "%"+request.Publisher+"%").Find(&books)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search books"})
		return
	}
	c.JSON(http.StatusOK, books)
	//avaliability of book checked
	// if books == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "No such Book exist"})
	// 	return
	// }

	//availability and expected return

	// for i, book := range books {
	// 	var issueRegistry models.IssueRegistry
	// 	if err := db.DB.Where("isbn = ? AND issue_status = ?", book.ISBN, "Pending").
	// 		First(&issueRegistry).Error; err == nil {
	// 		books[i].AvailableCopies = 0

	// 	}
	// }

}

// Raise an issue request
func RaiseIssueRequest(c *gin.Context) {
	var request struct {
		BookID string `json:"book_id"`
		Email  string `json:"email"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//checking whether validate user that is currently logged in
	// if request.Email != Data {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Token has been expired please login again"})
	// 	return
	// }

	// Check if the book exists
	var book models.BookInventory
	result := db.DB.Where("isbn = ?", request.BookID).First(&book)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book  is not found"})
		return
	}
	//providing approver id
	var col models.User
	result1 := db.DB.Where("email = ?", request.Email).First(&col)
	if result1.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Check if the book is available
	if book.AvailableCopies == 0 {
		// Create a new issue request
		issueRequest := models.RequestEvent{
			BookID:      request.BookID,
			ReaderID:    col.ID,
			RequestDate: time.Now(),
			RequestType: "Issue",
		}

		if err := db.DB.Create(&issueRequest).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create issue request"})
			return
		}

		//print issueRequest in json format
		c.JSON(http.StatusCreated, issueRequest)
		return
	}

	// Decrement available copies of the book
	book.AvailableCopies--
	if err := db.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book availability"})
		return
	}

}

// for authentication
var Data string
var jwtKey = []byte("secret_key")

type Credentials struct {
	Email string `json:"email"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	result := db.DB.Where("email = ?", credentials.Email).First(&user)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no record found"})
		return
	}

	expirationTime := time.Now().Add(time.Minute * 15)
	claims := &Claims{
		Email: credentials.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
}

// home

func Home(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No cookie available in token"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	//parsing tokken
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"eror": err.Error()})
		return
	}
	if !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	Data = claims.Email
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello,%s", claims.Email)})
	//printing token email

}
