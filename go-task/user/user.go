package user

import (
	"errors"
	"fmt"
	"sync"
)

type Config struct {
	userDB map[int]User // key is UserID, value is User
	mutex  sync.Mutex   // mutex for synchronization
}

type User struct {
	ID    int
	Name  string
	Email string
}

var incrID int

func InitApp() *Config {
	app := &Config{
		userDB: make(map[int]User),
	}

	incrID = 1

	return app
}

func (app *Config) CreateUser(user User) (User, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	if len(user.Name) == 0 || len(user.Email) == 0 {
		return User{}, errors.New("User data is invalid. Either name or email is empty")
	}

	for _, existingUser := range app.userDB {
		if existingUser.Name == user.Name || existingUser.Email == user.Email {
			return User{}, errors.New("User with the same name/email already exists")
		}
	}

	toAddUserID := incrID
	incrID++

	user.ID = toAddUserID

	app.userDB[toAddUserID] = user

	return user, nil
}

func (app *Config) ReadUser(id int) (User, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	if user, ok := app.userDB[id]; ok {
		return user, nil
	}

	return User{}, errors.New(fmt.Sprintf("No user exists with ID: [%d]", id))
}

func (app *Config) UpdateUser(toUpdateUser User) (User, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	curUser, err := app.ReadUser(toUpdateUser.ID)

	if err != nil {
		return User{}, err
	}

	shouldUpdate := false

	if len(toUpdateUser.Name) > 0 {
		curUser.Name = toUpdateUser.Name
		shouldUpdate = true
	}

	if len(toUpdateUser.Email) > 0 {
		curUser.Email = toUpdateUser.Email
		shouldUpdate = true
	}

	if !shouldUpdate {
		return User{}, errors.New("No new field value is provided for update")
	}

	app.userDB[curUser.ID] = curUser

	return curUser, nil
}

func (app *Config) DeleteUser(id int) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	_, err := app.ReadUser(id)

	if err != nil {
		return err
	}

	delete(app.userDB, id)

	return nil
}
