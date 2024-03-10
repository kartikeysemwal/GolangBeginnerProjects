package main

import (
	"fmt"
	"net/http"
	"strconv"

	"goproj.com/user"
)

type CreateUserReqPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserReqPayload struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (server *Server) Ping(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Your connection is good",
	}

	_ = server.writeJSON(w, http.StatusAccepted, payload)
}

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload CreateUserReqPayload

	err := server.readJSON(w, r, &requestPayload)

	if err != nil {
		server.errorJSON(w, err)
		return
	}

	user, err := server.userManager.CreateUser(user.User{
		Name:  requestPayload.Name,
		Email: requestPayload.Email,
	})

	if err != nil {
		server.errorJSON(w, err)
		return
	}

	var response CreateUserResponse
	response.Id = strconv.Itoa(user.ID)
	response.Name = user.Name
	response.Email = user.Email

	server.writeJSON(w, http.StatusAccepted, response)
}

func (server *Server) ReadUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request parameters or request body
	userID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		server.errorJSON(w, fmt.Errorf("invalid user ID"))
		return
	}

	// Call the ReadUser method from userManager
	foundUser, err := server.userManager.ReadUser(userID)
	if err != nil {
		server.errorJSON(w, err)
		return
	}

	// Create a response payload
	response := CreateUserResponse{
		Id:    strconv.Itoa(foundUser.ID),
		Name:  foundUser.Name,
		Email: foundUser.Email,
	}

	// Send the response
	server.writeJSON(w, http.StatusOK, response)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload UpdateUserReqPayload

	err := server.readJSON(w, r, &requestPayload)
	if err != nil {
		server.errorJSON(w, err)
		return
	}

	// Call the UpdateUser method from userManager
	updatedUser, err := server.userManager.UpdateUser(user.User{
		ID:    requestPayload.ID,
		Name:  requestPayload.Name,
		Email: requestPayload.Email,
	})

	if err != nil {
		server.errorJSON(w, err)
		return
	}

	// Create a response payload
	response := UpdateUserResponse{
		ID:    strconv.Itoa(updatedUser.ID),
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}

	// Send the response
	server.writeJSON(w, http.StatusOK, response)
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request parameters or request body
	userID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		server.errorJSON(w, fmt.Errorf("invalid user ID"))
		return
	}

	// Call the DeleteUser method from userManager
	err = server.userManager.DeleteUser(userID)
	if err != nil {
		server.errorJSON(w, err)
		return
	}

	// Send a success response
	server.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "User deleted successfully"})
}
