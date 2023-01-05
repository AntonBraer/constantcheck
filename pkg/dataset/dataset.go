package dataset

import (
	"fmt"
	"net/http"

	"github.com/AntonBraer/constantcheck/pkg/dataset/helps"
)

func CheckDataset() {
	fmt.Println(http.MethodGet)
	fmt.Println("GET") // want `GET literal contains in constant with name MethodGet`
	if err := helps.Ping(); err != nil {
		fmt.Println("internalError") // want `internalError literal contains in constant with name CustomError`
	}
}
