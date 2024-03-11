package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func connectDb() *sql.DB {
	db, err := sql.Open("mysql", "vroom:vroom@tcp(localhost:3306)/")
	if err != nil {
		panic(err.Error())
	}

	// Use the vroom database
	_, err = db.Exec("USE vroom")
	if err != nil {
		// Create the vroom database
		_, _ = db.Exec("CREATE DATABASE IF NOT EXISTS vroom")
		fmt.Println("database vroom created successfully")

		setupDb(db)

		_, _ = db.Exec("USE vroom")
	}

	return db
}

func setupDb(db *sql.DB) {
	// Create the variants table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS variants (
        variant_id INT AUTO_INCREMENT PRIMARY KEY,
        task_identifier VARCHAR(255),
        gpu_memory INT,
        gpu_cores INT,
        image VARCHAR(255),
        latency FLOAT,
        accuracy FLOAT,
		batch_size INT
    )`)
	if err != nil {
		db.Close()
		panic(err.Error())
	}

	fmt.Println("table 'variants' created successfully")
	db.Close()
}

func insertDb(funcInfo FuncInfo) {
	db := connectDb()

	stmt, _ := db.Prepare("INSERT INTO variants (task_identifier, gpu_memory, gpu_cores, image, latency, accuracy, batch_size) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()

	_, err := stmt.Exec(
		funcInfo.task_identifier,
		funcInfo.gpu_memory,
		funcInfo.gpu_cores,
		funcInfo.image,
		funcInfo.latency,
		funcInfo.accuracy,
		funcInfo.batch_size,
	)
	if err != nil {
		db.Close()
		panic(err.Error())
	}

	fmt.Println("entry variant recorded successfully")
	db.Close()
}

func getVariantsForReq(funcReq FuncReq, remaining_time float32) ([]FuncInfo, error) {
	db := connectDb()

	// Query to fetch all relevant resource variants for the given task
	query := "SELECT * FROM variants WHERE task_identifier = ? AND accuracy >= ? AND latency <= ?;"

	// Execute query on the sql db
	rows, err := db.Query(query, funcReq.task_identifier, funcReq.accuracy, remaining_time)
	if err != nil {
		fmt.Println("Error executing query:", err)
	}
	defer rows.Close()

	var variants []FuncInfo

	for rows.Next() {
		var variant FuncInfo
		if err := rows.Scan(
			&variant.variant_id,
			&variant.task_identifier,
			&variant.gpu_memory,
			&variant.gpu_cores,
			&variant.image,
			&variant.latency,
			&variant.accuracy,
			&variant.batch_size,
		); err == nil {
			variants = append(variants, variant)
		}
	}

	db.Close()
	return variants, nil
}

func getMinimumLatencyVariantForReq(funcReq FuncReq) (FuncInfo, error) {
	db := connectDb()

	// Query to fetch all relevant resource variants for the given task
	query := "SELECT * FROM variants WHERE task_identifier = ? AND accuracy >= ? ORDER BY latency ASC LIMIT 1;"

	// Execute query on the sql db
	rows, err := db.Query(query, funcReq.task_identifier, funcReq.accuracy)
	if err != nil {
		fmt.Println("Error executing query:", err)
	}
	defer rows.Close()

	var variant FuncInfo

	for rows.Next() {
		if err := rows.Scan(
			&variant.variant_id,
			&variant.task_identifier,
			&variant.gpu_memory,
			&variant.gpu_cores,
			&variant.image,
			&variant.latency,
			&variant.accuracy,
			&variant.batch_size,
		); err != nil {
			db.Close()
			panic(err.Error())
		}
		break
	}

	db.Close()
	return variant, nil
}

func generateTestDb() {
	insertDb(FuncInfo{"null", "image-rec", 4, 25, "synergcseiitb/image-rec-resnet:1.6", 5000, 80.0, 200})
	insertDb(FuncInfo{"null", "image-rec", 4, 50, "synergcseiitb/image-rec-resnet:1.6", 3000, 80.0, 200})
	insertDb(FuncInfo{"null", "image-rec", 4, 100, "synergcseiitb/image-rec-resnet:1.6", 1000, 80.0, 200})
}
