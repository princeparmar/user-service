package repositories

import (
	"database/sql"
	"errors"
	"time"
)

type UserRole struct {
	UserID     int
	RoleID     int
	ExpiryDate time.Time
}

type UserRoleRepository interface {
	Create(*UserRole) error
	Get(int, int) (*UserRole, error)
	Update(*UserRole) error
	Delete(int, int) error
	GetAll() ([]*UserRole, error)
	GetRolesForUser(int) ([]*Role, error)
	GetAllAccess(userID int) ([]*Access, error)
}

type userRoleRepository struct {
	db *sql.DB
}

func NewUserRoleRepository(db *sql.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) Create(ur *UserRole) error {
	query := "INSERT INTO user_roles (user_id, role_id, expiry_date, created_date, updated_date) VALUES (?, ?, ?, NOW(), NOW())"
	result, err := r.db.Exec(query, ur.UserID, ur.RoleID, ur.ExpiryDate)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (r *userRoleRepository) Get(userID, roleID int) (*UserRole, error) {
	query := "SELECT user_id, role_id, expiry_date FROM user_roles WHERE user_id = ? AND role_id = ?"
	row := r.db.QueryRow(query, userID, roleID)
	userRole := &UserRole{}
	err := row.Scan(&userRole.UserID, &userRole.RoleID, &userRole.ExpiryDate)
	if err != nil {
		return nil, err
	}
	return userRole, nil
}

func (r *userRoleRepository) Update(ur *UserRole) error {
	query := "UPDATE user_roles SET expiry_date = ?, updated_date = NOW() WHERE user_id = ? AND role_id = ?"
	_, err := r.db.Exec(query, ur.ExpiryDate, ur.UserID, ur.RoleID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRoleRepository) Delete(userID, roleID int) error {
	query := "DELETE FROM user_roles WHERE user_id = ? AND role_id = ?"
	_, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRoleRepository) GetAll() ([]*UserRole, error) {
	query := "SELECT user_id, role_id, expiry_date FROM user_roles"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userRoles := []*UserRole{}

	for rows.Next() {
		userRole := &UserRole{}
		err := rows.Scan(&userRole.UserID, &userRole.RoleID, &userRole.ExpiryDate)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userRoles, nil
}

func (r *userRoleRepository) GetRolesForUser(userID int) ([]*Role, error) {
	query := "SELECT r.role_id, r.role_name FROM roles r INNER JOIN user_roles ur ON r.role_id = ur.role_id WHERE ur.user_id = ?"
	rows, err := r.db.Query(query, userID)
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

func (r *userRoleRepository) GetAllAccess(userID int) ([]*Access, error) {
	query := `
		SELECT DISTINCT access.access_id, access.access_name
		FROM user_role
		JOIN role ON user_role.role_id = role.role_id
		JOIN access_role ON role.role_id = access_role.role_id
		JOIN access ON access_role.access_id = access.access_id
		WHERE user_role.user_id = ?
	`
	rows, err := r.db.Query(query, userID)
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

	if len(accesses) == 0 {
		return nil, errors.New("no access found for the user")
	}

	return accesses, nil
}

func (r *userRoleRepository) CreateTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS user_roles (
            user_id INT NOT NULL,
            role_id INT NOT NULL,
            expiry_date DATETIME,
			created_date DATETIME NOT NULL DEFAULT NOW(),
			updated_date DATETIME NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, role_id),
            FOREIGN KEY (user_id) REFERENCES users(user_id),
            FOREIGN KEY (role_id) REFERENCES roles(role_id)
        )`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
