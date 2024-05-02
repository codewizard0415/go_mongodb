package microservicecontroller

import (
	"fmt"
	"net/http"
)

func SendPostToService1(response http.ResponseWriter, request *http.Request) {
	fmt.Println("send Post To Service 1 function called")
}
