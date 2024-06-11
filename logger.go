package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

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
	SelectedInstanceId string    // Id of the instance that executed the request
	FinalState         string    // State (completed, error)
	ErrorMessage       string    // Error message if the request ended in an error

	TotalTimeTaken  float32 // Total time taken to execute the request
	VariantAccuracy float32 // Accuracy of the variant on which request was executed
}

type Logger struct {
	DatabaseManager *DatabaseManager
}

func (l *Logger) newLog(funcReq *FuncReq, variantId string, variantAccuracy float32, errorMessage string, instanceId string) {
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
		SelectedInstanceId: instanceId,
		FinalState:         funcReq.State,
		ErrorMessage:       errorMessage,
		TotalTimeTaken:     float32(totalTimeTaken),
		VariantAccuracy:    variantAccuracy,
	}
	// l.DatabaseManager.insertLogInDb(&log)
	l.saveLogToCSV(&log)
}

func (l *Logger) saveLogToCSV(log *LogEntry) {
	file, err := os.OpenFile("log_entries.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		log.RequestId,
		log.TaskIdentifier,
		strconv.FormatFloat(float64(log.Deadline), 'f', -1, 32),
		strconv.FormatFloat(float64(log.Accuracy), 'f', -1, 32),
		strconv.Itoa(log.RequestSize),
		log.RegistrationTs.Format(time.RFC3339),
		log.DeployInstanceTs.Format(time.RFC3339),
		log.SentForExecutionTs.Format(time.RFC3339),
		log.ResponseTs.Format(time.RFC3339),
		log.SelectedNode,
		log.SelectedVariantId,
		log.SelectedInstanceId,
		log.FinalState,
		log.ErrorMessage,
		strconv.FormatFloat(float64(log.TotalTimeTaken), 'f', -1, 32),
		strconv.FormatFloat(float64(log.VariantAccuracy), 'f', -1, 32),
	}

	if err := writer.Write(record); err != nil {
		panic(err)
	}
}
