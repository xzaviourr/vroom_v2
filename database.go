package main

import (
	"database/sql"
	"fmt"
)

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/go-sql-driver/mysql"
// )

func connectDb() *sql.DB {
	db, err := sql.Open("mysql", "vroom:vroom@tcp(localhost:3306)/")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func setupDb() {
	db := connectDb()
	defer db.Close()
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
		Port INT
    )`)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Table 'Variant' created successfully")
}

func insertVariantInDb(variant *Variant) {
	db := connectDb()

	stmt, _ := db.Prepare("INSERT INTO Variant (TaskId, GpuMemory, GpuCores, " +
		"Image, StartupLatency, MinLatency, MeanLatency, MaxLatency, Accuracy, " +
		"BatchSize, EndPoint, Port) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()

	_, err := stmt.Exec(
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
	)
	if err != nil {
		db.Close()
		panic(err.Error())
	}

	fmt.Println("Insert 'Variant' executed sucessfully")
	db.Close()
}

func initializeVariants(resourceManager *ResourceManager) {
	db := connectDb()

	// Query to fetch all the variants stored in the database
	query := "SELECT * FROM variants;"

	// Execute query on the sql db
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var variant Variant
		if err := rows.Scan(
			&variant,
		); err == nil {
			resourceManager.variantStore.Variants[variant.Id] = &variant
		}
	}

	db.Close()
}

// func getVariantsForReq(funcReq FuncReq, remaining_time float32) ([]FuncInfo, error) {
// 	db := connectDb()

// 	// Query to fetch all relevant resource variants for the given task
// 	query := "SELECT * FROM variants WHERE task_identifier = ? AND accuracy >= ? AND latency <= ?;"

// 	// Execute query on the sql db
// 	rows, err := db.Query(query, funcReq.TaskIdentifier, funcReq.Accuracy, remaining_time)
// 	if err != nil {
// 		fmt.Println("Error executing query:", err)
// 	}
// 	defer rows.Close()

// 	var variants []FuncInfo

// 	for rows.Next() {
// 		var variant FuncInfo
// 		if err := rows.Scan(
// 			&variant.VariantId,
// 			&variant.TaskIdentifier,
// 			&variant.GpuMemory,
// 			&variant.GpuCores,
// 			&variant.Image,
// 			&variant.Latency,
// 			&variant.Accuracy,
// 			&variant.BatchSize,
// 			&variant.Port,
// 		); err == nil {
// 			variants = append(variants, variant)
// 		}
// 	}

// 	db.Close()
// 	return variants, nil
// }

// func getMinimumLatencyVariantForReq(funcReq FuncReq) (FuncInfo, error) {
// 	db := connectDb()

// 	// Query to fetch all relevant resource variants for the given task
// 	query := "SELECT * FROM variants WHERE task_identifier = ? AND accuracy >= ? ORDER BY latency ASC LIMIT 1;"

// 	// Execute query on the sql db
// 	rows, err := db.Query(query, funcReq.TaskIdentifier, funcReq.Accuracy)
// 	if err != nil {
// 		fmt.Println("Error executing query:", err)
// 	}
// 	defer rows.Close()

// 	var variant FuncInfo

// 	for rows.Next() {
// 		if err := rows.Scan(
// 			&variant.VariantId,
// 			&variant.TaskIdentifier,
// 			&variant.GpuMemory,
// 			&variant.GpuCores,
// 			&variant.Image,
// 			&variant.Latency,
// 			&variant.Accuracy,
// 			&variant.BatchSize,
// 			&variant.Port,
// 		); err != nil {
// 			db.Close()
// 			panic(err.Error())
// 		}
// 		break
// 	}

// 	db.Close()
// 	return variant, nil
// }

// func generateTestDb() {
// 	insertDb(FuncInfo{"null", "object-detection", 2, 50, "synergcseiitb/object-detection-resnet:v1", 40000, 80.0, 200, 5123})
// 	insertDb(FuncInfo{"null", "object-detection", 2, 75, "synergcseiitb/object-detection-resnet:v1", 33000, 80.0, 200, 5123})
// 	insertDb(FuncInfo{"null", "object-detection", 2, 100, "synergcseiitb/object-detection-resnet:v1", 24000, 80.0, 200, 5123})
// 	insertDb(FuncInfo{"null", "object-detection", 2, 50, "synergcseiitb/object-detection-yolos:v1", 32000, 60.0, 50, 5126})
// 	insertDb(FuncInfo{"null", "object-detection", 2, 75, "synergcseiitb/object-detection-yolos:v1", 23000, 60.0, 50, 5126})
// 	insertDb(FuncInfo{"null", "object-detection", 2, 100, "synergcseiitb/object-detection-yolos:v1", 17000, 60.0, 50, 5126})
// }
