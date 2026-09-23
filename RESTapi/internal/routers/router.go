package routers

import (
	"net"
	"strconv"
	"strings"

	"restapi/internal/httpwire/request"
	"restapi/internal/httpwire/response"
	"restapi/internal/repository"
	"restapi/internal/service"
)

func Dispatch(conn net.Conn, store repository.Store, req *request.HTTPRequestStruct) {
	parts := strings.Split(strings.Trim(req.Path, "/"), "/")

	if parts[0] == "" {
		response.WriteError(conn, 404, "Query not found from the client request...")
		return
	}

	switch parts[0] {
	case "users":
		// GET is fully wired to the store. Other methods still hit
		// the stubs below until their phase lands.
		if req.Method == "GET" {
			dispatchUserGET(conn, store, parts)
			return
		}
		switch req.Method {
		case "POST":
			service.User_POST_service(req.Query)
		case "PUT":
			service.User_PUT_service(req.Query)
		case "PATCH":
			service.User_PATCH_service(req.Query)
		case "DELETE":
			service.User_DELETE_service(req.Query)
		default:
			response.WriteError(conn, 405, "Method Not Allowed")
		}

	case "todos":
		switch req.Method {
		case "GET":
			service.Todo_GET_service(req.Query)
		case "POST":
			service.Todo_POST_service(req.Query)
		case "PUT":
			service.Todo_PUT_service(req.Query)
		case "PATCH":
			service.Todo_PATCH_service(req.Query)
		case "DELETE":
			service.Todo_DELETE_service(req.Query)
		default:
			response.WriteError(conn, 405, "Method Not Allowed")
		}

	default:
		response.WriteError(conn, 404, "Not Found")
	}
}

// dispatchUserGET routes GET /users (collection) vs GET /users/:id (item).
// A trailing slash ("/users/") trims to the collection.
func dispatchUserGET(conn net.Conn, store repository.Store, parts []string) {
	if len(parts) == 1 {
		service.UserListService(conn, store)
		return
	}

	if len(parts) == 2 {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			response.WriteError(conn, 400, "Invalid user id")
			return
		}
		service.UserGetService(conn, store, id)
		return
	}

	response.WriteError(conn, 404, "Not Found")
}
