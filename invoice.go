package main

import "time"

type Invoice struct {
	InvoiceID   string    `json:"invoiceID"`
	InvoiceDate time.Time `json:"invoiceDate"`
	ClientID    string    `json:"clientID"`
	Sessions    []string  `json:"sessions"`
	Paid        bool      `json:"paid"`
	PaidDate    time.Time `json:"paidDate"`
}
