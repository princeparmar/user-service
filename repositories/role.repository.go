package repositories

import (
	"database/sql"
)

type Role struct {
	ID   int
	Name string
}

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *Role) error {
	query := "INSERT INTO roles (role_name, created_date, updated_date) VALUES (?, NOW(), NOW())"
	result, err := r.db.Exec(query, role.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	role.ID = int(id)

	return nil
}

func (r *RoleRepository) Get(id int) (*Role, error) {
	query := "SELECT role_id, role_name FROM roles WHERE role_id = ?"
	row := r.db.QueryRow(query, id)
	role := &Role{}
	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) Update(role *Role) error {
	query := "UPDATE roles SET role_name = ?, updated_date = NOW() WHERE role_id = ?"
	_, err := r.db.Exec(query, role.Name, role.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoleRepository) Delete(id int) error {
	query := "DELETE FROM roles WHERE role_id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoleRepository) GetAll() ([]*Role, error) {
	query := "SELECT role_id, role_name FROM roles"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := []*Role{}

	for rows.Next() {
		role := &Role{}
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *RoleRepository) CreateTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS roles (
            role_id INT AUTO_INCREMENT PRIMARY KEY,
            role_name VARCHAR(255) NOT NULL UNIQUE,
			created_date DATETIME NOT NULL DEFAULT NOW(),
			updated_date DATETIME NOT NULL DEFAULT NOW()
			)`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
