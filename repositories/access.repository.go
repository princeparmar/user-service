package repositories

import (
	"database/sql"
)

type Access struct {
	ID   int
	Name string
}

// AccessRepository provides access to the access store
type AccessRepository struct {
	db *sql.DB
}

// NewAccessRepository returns a new instance of AccessRepository
func NewAccessRepository(db *sql.DB) *AccessRepository {
	return &AccessRepository{db: db}
}

// Create creates a new access in the database and sets its ID
func (r *AccessRepository) Create(access *Access) error {
	// Prepare the query to insert a new access object
	query := "INSERT INTO access (access_name) VALUES (?)"
	// Execute the query with the access name parameter
	result, err := r.db.Exec(query, access.Name)
	if err != nil {
		return err
	}

	// Get the ID of the newly inserted access object
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID of the access object
	access.ID = int(id)

	return nil
}

// Get retrieves an access object with the given ID from the database
func (r *AccessRepository) Get(id int) (*Access, error) {
	// Prepare the query to select an access object by ID
	query := "SELECT access_id, access_name FROM access WHERE access_id = ?"
	// Execute the query with the ID parameter
	row := r.db.QueryRow(query, id)
	access := &Access{}
	err := row.Scan(&access.ID, &access.Name)
	return access, err
}

// Update updates an access object in the database with the new data
func (r *AccessRepository) Update(access *Access) error {
	// Prepare the query to update an access object by ID
	query := "UPDATE access SET access_name = ? WHERE access_id = ?"
	// Execute the query with the access name and ID parameters
	_, err := r.db.Exec(query, access.Name, access.ID)
	return err
}

// Delete deletes an access object with the given ID from the database
func (r *AccessRepository) Delete(id int) error {
	// Prepare the query to delete an access object by ID
	query := "DELETE FROM access WHERE access_id = ?"
	// Execute the query with the ID parameter
	_, err := r.db.Exec(query, id)
	return err
}

// GetAll retrieves all access objects from the database
func (r *AccessRepository) GetAll() ([]*Access, error) {
	// Prepare the query to select all access objects
	query := "SELECT access_id, access_name FROM access"
	// Execute the query
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accesses := []*Access{}

	for rows.Next() {
		access := &Access{}
		err := rows.Scan(&access.ID, &access.Name)
		if err != nil {
			return nil, err
		}
		accesses = append(accesses, access)
	}

	return accesses, rows.Err()
}

// CreateTable creates the access table in the database
func (r *AccessRepository) CreateTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS access (
            access_id INT AUTO_INCREMENT PRIMARY KEY,
            access_name VARCHAR(255) NOT NULL UNIQUE,
			created_date DATETIME NOT NULL DEFAULT NOW(),
			updated_date DATETIME NOT NULL DEFAULT NOW()
        )`
	_, err := r.db.Exec(query)
	return err
}
