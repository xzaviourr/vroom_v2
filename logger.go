package main

import "time"

type LogEntry struct {
	RequestId          string    // Request Id
	TaskIdentifier     string    // Requested task
	Deadline           float32   // Deadline (in milli second)
	Accuracy           float32   // Accuracy requirement (in percentage)
	RequestSize        int       // Number of inference queries
	RegistrationTs     time.Time // Timestamp when this request got registered
	DeployInstanceTs   time.Time // Timestamp when instance deployment was initiated
	SentForExecutionTs time.Time // Timestamp when API call was made to execute this request
	ResponseTs         time.Time // Timestamp when the response was sent back
	SelectedNode       string    // Id of the node on which this request was executed
	SelectedVariantId  string    // Id of the variant selected
	FinalState         string    // State (completed, error)
	ErrorMessage       string    // Error message if the request ended in an error

	TotalTimeTaken  float32 // Total time taken to execute the request
	VariantAccuracy float32 // Accuracy of the variant on which request was executed
}

type Logger struct {
	DatabaseManager *DatabaseManager
}

func (l *Logger) newLog(funcReq *FuncReq, variantId string, variantAccuracy float32, errorMessage string) {
	totalTimeTaken := funcReq.ResponseTs.Sub(funcReq.RegistrationTs).Milliseconds()
	log := LogEntry{
		RequestId:          funcReq.Uid,
		TaskIdentifier:     funcReq.TaskIdentifier,
		Deadline:           funcReq.Deadline,
		Accuracy:           funcReq.Accuracy,
		RequestSize:        funcReq.RequestSize,
		RegistrationTs:     funcReq.RegistrationTs,
		DeployInstanceTs:   funcReq.DeployInstanceTs,
		SentForExecutionTs: funcReq.SentForExecutionTs,
		ResponseTs:         funcReq.ResponseTs,
		SelectedNode:       funcReq.SelectedNode,
		SelectedVariantId:  variantId,
		FinalState:         funcReq.State,
		ErrorMessage:       errorMessage,
		TotalTimeTaken:     float32(totalTimeTaken),
		VariantAccuracy:    variantAccuracy,
	}
	l.DatabaseManager.insertLogInDb(&log)
}
