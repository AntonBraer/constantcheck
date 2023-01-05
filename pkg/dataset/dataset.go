package dataset

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(http.MethodGet)
	fmt.Println("GET") // want `GET literal contains in constant with name MethodGet`
}
