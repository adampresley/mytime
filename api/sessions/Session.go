package sessions

import "time"

type Session struct {
	SessionID     int       `json:"sessionID"`
	ClientID      int       `json:"clientID"`
	ProjectID     int       `json:"projectID"`
	CategoryID    int       `json:"categoryID"`
	StartDateTime time.Time `json:"startDateTime"`
	EndDateTime   time.Time `json:"endDateTime"`
	Notes         string    `json:"notes"`
	Invoiced      bool      `json:"invoiced"`
	InvoiceDate   time.Time `json:"invoiceDate"`
	Paid          bool      `json:"paid"`
	PaidDate      time.Time `json:"paidDate"`
}

func (s Session) ID() (string, interface{}) {
	return "sessionID", s.SessionID
}

type SessionCollection []Session

type SessionSearch struct {
	CategoryCode string
	ClientCode   string
	Invoiced     bool
	Paid         bool
	ProjectCode  string
	SessionID    int
	SessionIDs   []int
}

type ActiveSession struct {
	ActiveSessionID string    `json:"activeSessionID"`
	ProjectID       int       `json:"projectID"`
	ClientID        int       `json:"clientID"`
	CategoryID      int       `json:"categoryID"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	Notes           string    `json:"notes"`
}

func (as ActiveSession) ID() (string, interface{}) {
	return "activeSessionID", as.ActiveSessionID
}
