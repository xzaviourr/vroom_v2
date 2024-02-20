package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func setup_db() {
	db, err := sql.Open("mysql", "vroom:vroom@tcp(localhost:3306)/")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Create the vroom database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS vroom")
	if err != nil {
		panic(err.Error())
	}

	// Use the vroom database
	_, err = db.Exec("USE vroom")
	if err != nil {
		panic(err.Error())
	}

	// Create the variants table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS variants (
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
		panic(err.Error())
	}

	fmt.Println("Database 'vroom' and table 'variants' created successfully.")

	stmt, _ := db.Prepare("INSERT INTO variants (task_identifier, gpu_memory, gpu_cores, image, latency, accuracy, batch_size) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()

	_, _ = stmt.Exec("image-rec", 8, 50, "synergcseiitb/image-rec-resnet:1.6", 2.0, 80.0, 500)
	_, _ = stmt.Exec("image-rec", 12, 50, "synergcseiitb/image-rec-resnet:1.6", 1.5, 80.0, 700)
	_, _ = stmt.Exec("image-rec", 8, 80, "synergcseiitb/image-rec-resnet:1.6", 1.5, 80.0, 500)
	_, _ = stmt.Exec("image-rec", 4, 50, "synergcseiitb/image-rec-resnet:1.6", 3.0, 80.0, 200)
}

func getVariantsForReq(funcReq FuncReq) ([]FuncInfo, error) {
	db, _ := sql.Open("mysql", "vroom:vroom@tcp(localhost:3306)/")
	_, _ = db.Exec("USE vroom")

	current_ts := time.Now()

	// Time remaining before SLO miss
	remaining_time := funcReq.deadline - float32(current_ts.Sub(funcReq.timestamp)/time.Millisecond)

	// Query to fetch all relevant resource variants for the given task
	query := "SELECT * FROM variants WHERE task_identifier = '?' AND accuracy >= ? AND latency <= ?;"

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

	return variants, nil
}
