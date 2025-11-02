package api

import "fmt"

// API Response Structure matching SMS.ir API
type APIResponse[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// GetStatusMessage returns the status message
func (r *APIResponse[T]) GetStatusMessage() string {
	statusMessages := map[int]string{
		0: "Failed",
		1: "Success",
	}

	if message, exists := statusMessages[r.Status]; exists {
		return message
	}
	return fmt.Sprintf("Unknown status: %d", r.Status)
}

// IsSuccess checks if the API response indicates success
func (r *APIResponse[T]) IsSuccess() bool {
	return r.Status == 1
}

// CreditResponse for GET /v1/credit
type CreditResponse float64

// LinesResponse for GET /v1/line
type LinesResponse []int64

// BulkSendRequest for POST /v1/send/bulk
type BulkSendRequest struct {
	LineNumber   int64    `json:"lineNumber"`
	MessageText  string   `json:"messageText"`
	Mobiles      []string `json:"mobiles"`
	SendDateTime *int64   `json:"sendDateTime,omitempty"` // Unix timestamp for scheduling
}

// BulkSendResponse
type BulkSendResponse struct {
	PackID     string  `json:"packId"`
	MessageIds []int32 `json:"messageIds"`
	Cost       float64 `json:"cost"`
}

// SendMessageReport for GET /v1/send/{messageId}
type ReportSendMessageResponse struct {
	MessageId        int32   `json:"messageId"`
	Mobile           int64   `json:"mobile"`
	MessageText      string  `json:"messageText"`
	SendDateTime     int64   `json:"sendDateTime"`
	LineNumber       int64   `json:"lineNumber"`
	Cost             float64 `json:"cost"`
	DeliveryState    int     `json:"deliveryState"`
	DeliveryDateTime int64   `json:"deliveryDateTime"`
	Status           string  `json:"status"`
}

// SendPackReport for GET /v1/send/pack/{packId}
type ReportSendPackResponse struct {
	PackID      string                      `json:"packId"`
	TotalCount  int32                       `json:"totalCount"`
	SentCount   int32                       `json:"sentCount"`
	FailedCount int32                       `json:"failedCount"`
	Messages    []ReportSendMessageResponse `json:"messages"`
}

// RemoveScheduledResponse for DELETE /v1/send/scheduled/{packId}
type RemoveScheduledResponse struct {
	ReturnedCreditCount float64 `json:"returnedCreditCount"`
	SmsCount            int32   `json:"smsCount"`
	Success             bool    `json:"success"`
	Message             string  `json:"message"`
}
