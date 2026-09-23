package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"restapi/internal/httpwire/request"
	"restapi/internal/httpwire/response"
	"restapi/internal/repository"
)

type UserCreateStruct struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UserListService(conn net.Conn, store repository.Store) {
	users, err := store.ListUsers()
	if err != nil {
		log.Printf("ListUsers: %v", err)
		response.WriteError(conn, 500, "Internal Server Error")
		return
	}

	response.WriteJSON(conn, 200, users)
}

func UserGetService(conn net.Conn, store repository.Store, id int) {
	user, err := store.GetUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.WriteError(conn, 404, "User not found")
			return
		}
		log.Printf("GetUser(%d): %v", id, err)
		response.WriteError(conn, 500, "Internal Server Error")
		return
	}

	response.WriteJSON(conn, 200, user)
}

func UserCreateService(conn net.Conn, store repository.Store, req *request.HTTPRequestStruct) {
	var input UserCreateStruct

	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		response.WriteError(conn, 400, "Invalid JSON body")
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if input.Name == "" {
		response.WriteError(conn, 400, "Name field should not be left empty")
		return
	}

	if input.Email == "" {
		response.WriteError(conn, 400, "Email field should not be left empty")
		return
	}

	if !strings.Contains(input.Email, "@") {
		response.WriteError(conn, 400, "The email format is not valid. Please enter a valid email.")
		return
	}

	user, err := store.CreateUser(input.Name, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			response.WriteError(conn, 409, "email already exists")
			return
		}
		log.Printf("CreateUser: %v", err)
		response.WriteError(conn, 500, "Internal Server Error")
		return
	}

	response.WriteJSON(conn, 201, user)
}

func UserPutService(conn net.Conn, store repository.Store, req *request.HTTPRequestStruct) {
	var input UserCreateStruct

	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		response.WriteError(conn, 400, "Invalid JSON body")
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if input.Name == "" {
		response.WriteError(conn, 400, "Name field cannot be left empty")
		return
	}

	if input.Email == "" {
		response.WriteError(conn, 400, "Email field cannot be left empty")
		return
	}

	if !strings.Contains(input.Email, "@") {
		response.WriteError(conn, 400, "The email format is not valid. Please enter a valid email")
		return
	}

	parts := strings.Split(strings.Trim(req.Path, "/"), "/")

	if len(parts) != 2 {
		response.WriteError(conn, 404, "Not found")
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		response.WriteError(conn, 400, "Invalid user ID")
		return
	}

	user, err := store.UpdateUser(id, input.Name, input.Email)

	if err != nil {
		response.WriteError(conn, 500, "Internal Server Error")
		return
	}

	response.WriteJSON(conn, 200, user)
}

func User_PATCH_service(query map[string]string) {
	fmt.Println("USER PATCH SERVICE", query)
}

func User_DELETE_service(query map[string]string) {
	fmt.Println("USER DELETE SERVICE", query)
}
