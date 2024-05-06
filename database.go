package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Database handler
type DatabaseManager struct {
}

func initDatabaseManager() *DatabaseManager {
	databaseManager := DatabaseManager{}
	databaseManager.setupDb()
	return &databaseManager
}

// Returns a connection object with the db
func (dm *DatabaseManager) connectDb() *sql.DB {
	db, err := sql.Open("mysql", "vroom:vroom@tcp(localhost:3306)/")
	if err != nil {
		panic(err.Error())
	}
	db.Exec("USE vroom;")
	return db
}

// Creates necessary db schema
func (dm *DatabaseManager) setupDb() {
	db := dm.connectDb()

	// Create the variants table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Variant (
        Id VARCHAR(255) PRIMARY KEY,
        TaskId VARCHAR(255),
        GpuMemory INT,
        GpuCores INT,
        Image VARCHAR(255),
        StartupLatency FLOAT,
		MinLatency FLOAT,
		MeanLatency FLOAT,
		MaxLatency FLOAT,
        Accuracy FLOAT,
		BatchSize INT,
		EndPoint VARCHAR(255),
		Port INT,
		Capacity FLOAT
    );`)

	if err != nil {
		db.Close()
		panic(err.Error())
	}

	// Create the logs table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Logs (
		RequestId VARCHAR(255) NOT NULL,
		TaskIdentifier VARCHAR(255),
		Deadline FLOAT,
		Accuracy FLOAT,
		RequestSize INT,
		RegistrationTs TIMESTAMP,
		DeployInstanceTs TIMESTAMP,
		SentForExecutionTs TIMESTAMP,
		ResponseTs TIMESTAMP,
		SelectedNode VARCHAR(255),
		SelectedVariantId VARCHAR(255),
		FinalState VARCHAR(255),
		ErrorMessage TEXT,
		TotalTimeTaken FLOAT,
		VariantAccuracy FLOAT,
		PRIMARY KEY (RequestId)
	);`)

	if err != nil {
		db.Close()
		panic(err.Error())
	}

	fmt.Println("Database initialized successfully")
	db.Close()
}

// Insert a variant into the Db
func (dm *DatabaseManager) insertVariantInDb(variant *Variant) {
	db := dm.connectDb()

	stmt, _ := db.Prepare("INSERT INTO Variant (Id, TaskId, GpuMemory, GpuCores, " +
		"Image, StartupLatency, MinLatency, MeanLatency, MaxLatency, Accuracy, " +
		"BatchSize, EndPoint, Port, Capacity) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	_, err := stmt.Exec(
		variant.Id,
		variant.TaskId,
		variant.GpuMemory,
		variant.GpuCores,
		variant.Image,
		variant.StartupLatency,
		variant.MinLatency,
		variant.MeanLatency,
		variant.MaxLatency,
		variant.Accuracy,
		variant.BatchSize,
		variant.EndPoint,
		variant.Port,
		variant.Capacity,
	)
	if err != nil {
		db.Close()
		panic(err.Error())
	}

	fmt.Println("Variant ", variant.Id, " inserted in the database")
	db.Close()
}

func (dm *DatabaseManager) insertLogInDb(log *LogEntry) {
	db := dm.connectDb()

	stmt, _ := db.Prepare("INSERT INTO Logs (RequestId, TaskIdentifier, Deadline, Accuracy, " +
		"RequestSize, RegistrationTs, DeployInstanceTs, SentForExecutionTs, ResponseTs, SelectedNode, " +
		"SelectedVariantId, FinalState, ErrorMessage, TotalTimeTaken, VariantAccuracy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	_, err := stmt.Exec(
		log.RequestId,
		log.TaskIdentifier,
		log.Deadline,
		log.Accuracy,
		log.RequestSize,
		log.RegistrationTs,
		log.DeployInstanceTs,
		log.SentForExecutionTs,
		log.ResponseTs,
		log.SelectedNode,
		log.SelectedVariantId,
		log.FinalState,
		log.ErrorMessage,
		log.TotalTimeTaken,
		log.VariantAccuracy,
	)
	if err != nil {
		db.Close()
		panic(err.Error())
	}

	fmt.Println("Response logged")
	db.Close()
}

// Fetches all the variants from the db
func (dm *DatabaseManager) loadAllVariantsFromDb() map[string]*Variant {
	db := dm.connectDb()

	// Query to fetch all the variants stored in the database
	query := "SELECT * FROM Variant;"

	// Execute query on the sql db
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
	}
	defer rows.Close()

	variants := make(map[string]*Variant)

	for rows.Next() {
		var variant Variant
		if err := rows.Scan(
			&variant.Id,
			&variant.TaskId,
			&variant.GpuMemory,
			&variant.GpuCores,
			&variant.Image,
			&variant.StartupLatency,
			&variant.MinLatency,
			&variant.MeanLatency,
			&variant.MaxLatency,
			&variant.Accuracy,
			&variant.BatchSize,
			&variant.EndPoint,
			&variant.Port,
			&variant.Capacity,
		); err == nil {
			variants[variant.Id] = &variant
		}
	}

	fmt.Println("Total number of variants loaded from DB: ", len(variants))
	db.Close()
	return variants
}
