package service

import (
	"errors"
	"fmt"
	"log"
	"net"

	"restapi/internal/httpwire/response"
	"restapi/internal/repository"
)

// UserListService handles GET /users: 200 with the JSON user array.
func UserListService(conn net.Conn, store repository.Store) {
	users, err := store.ListUsers()
	if err != nil {
		log.Printf("ListUsers: %v", err)
		response.WriteError(conn, 500, "Internal Server Error")
		return
	}

	response.WriteJSON(conn, 200, users)
}

// UserGetService handles GET /users/:id: 200 the user, 404 when missing.
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

func User_GET_service(query map[string]string) {
	fmt.Println("USER GET SERVICE", query)
}

func User_POST_service(query map[string]string) {
	fmt.Println("USER POST SERVICE", query)
}

func User_PUT_service(query map[string]string) {
	fmt.Println("USER PUT SERVICE", query)
}

func User_PATCH_service(query map[string]string) {
	fmt.Println("USER PATCH SERVICE", query)
}

func User_DELETE_service(query map[string]string) {
	fmt.Println("USER DELETE SERVICE", query)
}
