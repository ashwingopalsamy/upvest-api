package handler

const (
	// Error Titles
	ErrTitleDatabaseError   = "Database Error"
	ErrTitleInvalidRequest  = "Invalid Request"
	ErrTitleKafkaError      = "Kafka Error"
	ErrTitleValidationError = "Validation Error"

	// Error Messages
	ErrMsgCreateUserFailed   = "failed to create user"
	ErrMsgEmitEventFailed    = "failed to emit user creation event"
	ErrMsgInvalidRequestBody = "request body could not be parsed"
	ErrMsgMarshalEventFailed = "failed to marshal event"
)
