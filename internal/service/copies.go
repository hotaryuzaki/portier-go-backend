package service

import (
	"context"
	"fmt"
	"log"
	"portier/pkg/db"
	"time"
)

type Copy struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	KeyID     int       `json:"key_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by,omitempty"` // Optional, nullable field
	IsActive  bool      `json:"is_active"`
}

// GetAllCopies fetches all copies
func GetAllCopies() ([]Copy, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	var copies []Copy

	rows, err := dbConn.Query(ctx, "SELECT id, name, key_id, created_at, created_by, is_active FROM copies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var copy Copy
		if err := rows.Scan(&copy.ID, &copy.Name, &copy.KeyID, &copy.CreatedAt, &copy.CreatedBy, &copy.IsActive); err != nil {
			return nil, err
		}
		copies = append(copies, copy)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return copies, nil
}

// GetCopyByID fetches a copy by its ID
func GetCopyByID(id int) (Copy, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Add context

	var copy Copy

	query := `SELECT id, name, key_id, created_at, created_by, is_active FROM copies WHERE id=$1`
	err := dbConn.QueryRow(ctx, query, id).Scan(&copy.ID, &copy.Name, &copy.KeyID, &copy.CreatedAt, &copy.CreatedBy, &copy.IsActive)
	if err != nil {
		return Copy{}, err
	}

	return copy, nil
}

// CreateCopy creates a new copy
func CreateCopy(copy Copy) (Copy, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	// Explicitly set the default value for IsActive
	copy.IsActive = true

	query := `INSERT INTO copies (name, key_id, created_at, created_by, is_active) 
						VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int
	err := dbConn.QueryRow(ctx, query, copy.Name, copy.KeyID, time.Now(), copy.CreatedBy, copy.IsActive).Scan(&id)
	if err != nil {
		log.Printf("Error creating copy: %v", err)
		return Copy{}, fmt.Errorf("failed to create copy: %v", err)
	}

	copy.ID = id
	return copy, nil
}

// UpdateCopy updates a copy's information
func UpdateCopy(id int, copy Copy) (Copy, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	query := `UPDATE copies SET name=$1, key_id=$2, created_by=$3, is_active=$4 WHERE id=$5`
	_, err := dbConn.Exec(ctx, query, copy.Name, copy.KeyID, copy.CreatedBy, copy.IsActive, id)
	if err != nil {
		return Copy{}, err
	}

	copy.ID = id
	return copy, nil
}

// DeleteCopy deletes a copy
func DeleteCopy(id int) error {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	query := `DELETE FROM copies WHERE id=$1`
	_, err := dbConn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}