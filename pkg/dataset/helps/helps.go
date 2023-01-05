package helps

import (
	"fmt"
)

const CustomError = "internalError"

func Ping() error {
	fmt.Println("Pong")
	return nil
}
