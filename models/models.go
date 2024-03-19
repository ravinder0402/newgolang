package models

import "time"

type Library struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"unique"`
}

type User struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string `json:"name" gorm:"unique"`
	Email         string `json:"email" gorm:"unique"`
	ContactNumber string `json:"contact_number"`
	Role          string `json:"role" gorm:"type:enum('Admin', 'Reader')"`
	LibID         uint   `json:"lib_id"`
}

type BookInventory struct {
	ISBN            string `json:"isbn" gorm:"primaryKey;size:20"`
	LibID           uint   `json:"lib_id"`
	Title           string `json:"title"`
	Authors         string `json:"authors"`
	Publisher       string `json:"publisher"`
	Version         int    `json:"version"`
	TotalCopies     int    `json:"total_copies"`
	AvailableCopies int    `json:"available_copies"`
}

type RequestEvent struct {
	ReqID        uint      `json:"req_id" gorm:"primaryKey;autoIncrement"`
	BookID       string    `json:"book_id" gorm:"primaryKey;size:20"`
	ReaderID     uint      `json:"reader_id"`
	RequestDate  time.Time `json:"request_date" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	ApprovalDate time.Time `json:"approval_date" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	ApproverID   uint      `json:"approver_id"`
	RequestType  string    `json:"request_type" gorm:"type:enum('Issue', 'Return')"`
}

type IssueRegistry struct {
	IssueID            uint      `json:"issue_id" gorm:"primaryKey;autoIncrement"`
	ISBN               string    `json:"isbn" gorm:"primaryKey;size:20"`
	ReaderID           uint      `json:"reader_id"`
	IssueApproverID    uint      `json:"issue_approver_id"`
	IssueStatus        string    `json:"issue_status" gorm:"type:enum('Approved', 'Rejected', 'Pending')"`
	IssueDate          time.Time `json:"issue_date" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	ExpectedReturnDate time.Time `json:"expected_return_date" gorm:"type:date"`
	ReturnDate         time.Time `json:"return_date" gorm:"type:date"`
	ReturnApproverID   uint      `json:"return_approver_id"`
}
