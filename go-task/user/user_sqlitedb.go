// user_db.go
package user

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUserConfig struct {
	db    *sql.DB
	mutex sync.Mutex
}

func InitSQLiteUserApp(databasePath string) (*SQLiteUserConfig, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}

	// Create the users table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteUserConfig{
		db: db,
	}, nil
}

func (app *SQLiteUserConfig) CreateUser(user User) (User, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	if len(user.Name) == 0 || len(user.Email) == 0 {
		return User{}, errors.New("User data is invalid. Either name or email is empty")
	}

	// Check if the user with the same name or email already exists
	var count int
	err := app.db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ? OR email = ?", user.Name, user.Email).Scan(&count)
	if err != nil {
		return User{}, err
	}

	if count > 0 {
		return User{}, errors.New("User with the same name/email already exists")
	}

	// Insert the user into the database
	result, err := app.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return User{}, err
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}

	user.ID = int(lastInsertID)
	return user, nil
}

func (app *SQLiteUserConfig) ReadUser(id int) (User, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	var user User
	err := app.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("No user exists with ID: [%d]", id))
	}

	return user, nil
}

func (app *SQLiteUserConfig) UpdateUser(toUpdateUser User) (User, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	// Check if the user exists
	var count int
	err := app.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", toUpdateUser.ID).Scan(&count)
	if err != nil {
		return User{}, err
	}

	if count == 0 {
		return User{}, errors.New("User not found")
	}

	// Update the user in the database
	_, err = app.db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", toUpdateUser.Name, toUpdateUser.Email, toUpdateUser.ID)
	if err != nil {
		return User{}, err
	}

	return toUpdateUser, nil
}

func (app *SQLiteUserConfig) DeleteUser(id int) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	// Check if the user exists
	var count int
	err := app.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New(fmt.Sprintf("No user exists with ID: [%d]", id))
	}

	// Delete the user from the database
	_, err = app.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
