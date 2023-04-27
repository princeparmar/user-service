package repositories

import (
	"database/sql"
	"errors"
)

type User struct {
	ID       int
	UserName string
	Mobile   string
	EmailID  string
}

// UserRepository defines a struct for User data storage and retrieval.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance using the provided database connection.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new User record into the database.
func (r *UserRepository) Create(user *User) error {
	query := "INSERT INTO users (user_name, mobile, created_date, updated_date, email_id) VALUES (?, ?, NOW(), NOW(), ?)"
	result, err := r.db.Exec(query, user.UserName, user.Mobile, user.EmailID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)

	return nil
}

// Get retrieves a User record from the database by ID.
func (r *UserRepository) Get(id int) (*User, error) {
	query := "SELECT user_id, user_name, mobile, email_id FROM users WHERE user_id = ?"
	row := r.db.QueryRow(query, id)
	user := &User{}
	err := row.Scan(&user.ID, &user.UserName, &user.Mobile, &user.EmailID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetAll() ([]*User, error) {
	query := "SELECT user_id, user_name, email_id, mobile FROM user"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.UserName, &user.EmailID, &user.Mobile)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates an existing User record in the database.
func (r *UserRepository) Update(user *User) error {
	query := "UPDATE users SET user_name = ?, mobile = ?, updated_date = NOW(), email_id = ? WHERE user_id = ?"
	result, err := r.db.Exec(query, user.UserName, user.Mobile, user.EmailID, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were affected during the update")
	}

	return nil
}

// Delete removes a User record from the database by ID.
func (r *UserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE user_id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were affected during the delete")
	}

	return nil
}

// List retrieves a list of all User records from the database.
func (r *UserRepository) List() ([]*User, error) {
	query := "SELECT user_id, user_name, mobile, email_id FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.UserName, &user.Mobile, &user.EmailID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetPassword retrieves the password of a user from the database by user_id.
func (r *UserRepository) GetPassword(id int) (string, error) {
	query := "SELECT password FROM users WHERE user_id = ?"
	row := r.db.QueryRow(query, id)
	var password string
	err := row.Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid user id")
		}
		return "", err
	}
	return password, nil
}

// GetUserByUserName retrieves a User record from the database by user_name.
func (r *UserRepository) GetUserByUserName(userName string) (*User, error) {
	query := "SELECT user_id, user_name, mobile, email_id FROM users WHERE user_name = ?"
	row := r.db.QueryRow(query, userName)
	user := &User{}
	err := row.Scan(&user.ID, &user.UserName, &user.Mobile, &user.EmailID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid user name")
		}
		return nil, err
	}
	return user, nil
}

// UpdatePassword updates the password of an existing User record in the database.
func (r *UserRepository) UpdatePassword(userID int, password string) error {
	query := "UPDATE users SET password = ?, updated_date = NOW() WHERE user_id = ?"
	result, err := r.db.Exec(query, password, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were affected during the update")
	}

	return nil
}

// CreateTable creates the 'users' table in the database.
func (ur *UserRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INT AUTO_INCREMENT PRIMARY KEY,
		user_name VARCHAR(255) NOT NULL UNIQUE,
		mobile VARCHAR(10) NOT NULL,
		email_id VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_date DATETIME NOT NULL DEFAULT NOW(),
		updated_date DATETIME NOT NULL DEFAULT NOW()
	)	
	`

	_, err := ur.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
