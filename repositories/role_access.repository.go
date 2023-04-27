package repositories

import (
	"database/sql"
)

type RoleAccess struct {
	RoleID   int
	AccessID int
}

// RoleAccessRepository is a struct that handles all database operations related to RoleAccess
type RoleAccessRepository struct {
	db *sql.DB
}

// NewRoleAccessRepository creates a new RoleAccessRepository with the given db instance
func NewRoleAccessRepository(db *sql.DB) *RoleAccessRepository {
	return &RoleAccessRepository{db}
}

// Create creates a new role access object in the database
func (r *RoleAccessRepository) Create(ra *RoleAccess) error {
	query := "INSERT INTO access_role (role_id, access_id, created_date, updated_date) VALUES (?, ?, NOW(), NOW())"
	result, err := r.db.Exec(query, ra.RoleID, ra.AccessID)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves a role access object with the given role ID and access ID from the database
func (r *RoleAccessRepository) Get(roleID, accessID int) (*RoleAccess, error) {
	query := "SELECT role_id, access_id FROM access_role WHERE role_id = ? AND access_id = ?"
	row := r.db.QueryRow(query, roleID, accessID)
	roleAccess := &RoleAccess{}
	err := row.Scan(&roleAccess.RoleID, &roleAccess.AccessID)
	if err != nil {
		return nil, err
	}
	return roleAccess, nil
}

// Delete deletes a role access object with the given role ID and access ID from the database
func (r *RoleAccessRepository) Delete(roleID, accessID int) error {
	query := "DELETE FROM access_role WHERE role_id = ? AND access_id = ?"
	_, err := r.db.Exec(query, roleID, accessID)
	if err != nil {
		return err
	}

	return nil
}

// GetAll retrieves all role access objects from the database
func (r *RoleAccessRepository) GetAll() ([]*RoleAccess, error) {
	query := "SELECT role_id, access_id FROM access_role"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roleAccesses := []*RoleAccess{}

	for rows.Next() {
		roleAccess := &RoleAccess{}
		err := rows.Scan(&roleAccess.RoleID, &roleAccess.AccessID)
		if err != nil {
			return nil, err
		}
		roleAccesses = append(roleAccesses, roleAccess)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roleAccesses, nil
}

func (r *RoleRepository) GetAccessesForRole(roleID int) ([]*Access, error) {
	query := "SELECT a.access_id, a.access_name FROM access a INNER JOIN access_role ar ON a.access_id = ar.access_id WHERE ar.role_id = ?"
	rows, err := r.db.Query(query, roleID)
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accesses, nil
}

func (r *RoleAccessRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS access_role (
		role_id INT NOT NULL,
		access_id INT NOT NULL,
		created_date DATETIME NOT NULL DEFAULT NOW(),
		updated_date DATETIME NOT NULL DEFAULT NOW(),
		PRIMARY KEY (role_id, access_id),
		FOREIGN KEY (role_id) REFERENCES role(role_id) ON DELETE CASCADE,
		FOREIGN KEY (access_id) REFERENCES access(access_id) ON DELETE CASCADE
	)	
`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
