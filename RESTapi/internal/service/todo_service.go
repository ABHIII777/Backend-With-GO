package service

import (
	"fmt"
	"strconv"
)

func Todo_GET_service(query map[string]string) {
	fmt.Println("TODO GET SERVICE", query)

	var userID *int
	if raw, ok := query["user_id"]; ok && raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			fmt.Println("TODO GET: invalid user_id:", raw)
		} else {
			userID = &n
		}
	}

	var completed *bool
	if raw, ok := query["completed"]; ok && raw != "" {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			fmt.Println("TODO GET: invalid completed:", raw)
		} else {
			completed = &b
		}
	}

	fmt.Println("TODO GET parsed filters userID:", userID, "completed:", completed)
}

func Todo_POST_service(query map[string]string) {
	fmt.Println("TODO POST SERVICE", query)
}

func Todo_PUT_service(query map[string]string) {
	fmt.Println("TODO PUT SERVICE", query)
}

func Todo_PATCH_service(query map[string]string) {
	fmt.Println("TODO PATCH SERVICE", query)
}

func Todo_DELETE_service(query map[string]string) {
	fmt.Println("TODO DELETE SERVICE", query)
}

