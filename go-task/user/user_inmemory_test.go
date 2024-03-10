package user

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestCreateUser(t *testing.T) {
	app := InitApp()

	// Test valid user creation
	user1 := User{Name: "John", Email: "john@example.com"}
	createdUser, err := app.CreateUser(user1)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	if createdUser.ID == 0 {
		t.Error("Expected user ID to be assigned")
	}

	// Test creating a user with invalid data
	invalidUser := User{Name: "", Email: "invalid@example.com"}
	_, err = app.CreateUser(invalidUser)
	if err == nil {
		t.Error("Expected an error for creating a user with invalid data")
	}
}

func TestReadUser(t *testing.T) {
	app := InitApp()

	// Test reading a non-existent user
	_, err := app.ReadUser(1)
	if err == nil {
		t.Error("Expected an error for reading a non-existent user")
	}

	// Test reading an existing user
	user := User{Name: "Alice", Email: "alice@example.com"}
	createdUser, _ := app.CreateUser(user)

	readUser, err := app.ReadUser(createdUser.ID)
	if err != nil {
		t.Errorf("Error reading user: %v", err)
	}

	if readUser.ID != createdUser.ID {
		t.Error("Expected to read the correct user")
	}
}

func TestUpdateUser(t *testing.T) {
	app := InitApp()

	// Test updating a non-existent user
	invalidUser := User{ID: 1, Name: "Bob", Email: "bob@example.com"}
	_, err := app.UpdateUser(invalidUser)
	if err == nil {
		t.Error("Expected an error for updating a non-existent user")
	}

	// Test updating an existing user
	user := User{Name: "Charlie", Email: "charlie@example.com"}
	createdUser, _ := app.CreateUser(user)

	updatedUser := User{ID: createdUser.ID, Name: "Charlie Updated", Email: "charlie.updated@example.com"}
	_, err = app.UpdateUser(updatedUser)
	if err != nil {
		t.Errorf("Error updating user: %v", err)
	}

	readUser, _ := app.ReadUser(createdUser.ID)
	if readUser.Name != updatedUser.Name || readUser.Email != updatedUser.Email {
		t.Error("Expected to update the user with new data")
	}
}

func TestDeleteUser(t *testing.T) {
	app := InitApp()

	// Test deleting a non-existent user
	err := app.DeleteUser(1)
	if err == nil {
		t.Error("Expected an error for deleting a non-existent user")
	}

	// Test deleting an existing user
	user := User{Name: "David", Email: "david@example.com"}
	createdUser, _ := app.CreateUser(user)

	err = app.DeleteUser(createdUser.ID)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}

	// Verify that the user is deleted
	_, err = app.ReadUser(createdUser.ID)
	if err == nil {
		t.Error("Expected an error for reading a deleted user")
	}
}

func TestConcurrentCreateUser(t *testing.T) {
	app := InitApp()
	var wg sync.WaitGroup
	numRoutines := 100

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(index int) {
			defer wg.Done()
			user := User{Name: fmt.Sprintf("User%d", index), Email: fmt.Sprintf("user%d@example.com", index)}
			_, err := app.CreateUser(user)
			if err != nil {
				t.Errorf("Error creating user in goroutine: %v", err)
			}
		}(i)
	}

	wg.Wait()

	if len(app.userDB) != numRoutines {
		t.Errorf("Expected %d users, got %d", numRoutines, len(app.userDB))
	}
}

func TestMain(m *testing.M) {
	// Run the tests
	exitCode := m.Run()

	// Exit with the test result
	os.Exit(exitCode)
}
