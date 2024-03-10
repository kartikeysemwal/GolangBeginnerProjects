package user

import (
	"errors"
	"fmt"
	"sync"
)

type InMemoryUserConfig struct {
	userDB map[int]User // key is UserID, value is User
	mutex  sync.Mutex   // mutex for synchronization
	incrID int
}

func InitInMemoryUserApp() *InMemoryUserConfig {
	app := &InMemoryUserConfig{
		userDB: make(map[int]User),
	}

	app.incrID = 1

	return app
}

func (app *InMemoryUserConfig) CreateUser(user User) (User, error) {
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

	toAddUserID := app.incrID
	app.incrID++

	user.ID = toAddUserID

	app.userDB[toAddUserID] = user

	return user, nil
}

func (app *InMemoryUserConfig) ReadUser(id int) (User, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	if user, ok := app.userDB[id]; ok {
		return user, nil
	}

	return User{}, errors.New(fmt.Sprintf("No user exists with ID: [%d]", id))
}

func (app *InMemoryUserConfig) UpdateUser(toUpdateUser User) (User, error) {
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

func (app *InMemoryUserConfig) DeleteUser(id int) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	_, err := app.ReadUser(id)

	if err != nil {
		return err
	}

	delete(app.userDB, id)

	return nil
}
