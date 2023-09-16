package model

// SQS returns a map of attributes rather than a struct
type SQSQueue struct {
	Name string
	QueueUrl string
	IsFifo bool
}
