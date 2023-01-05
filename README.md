<div align="center">

# constantcheck

</div>

---

constantcheck is a golang analyzer for checking 
the possibility of replacing the literals used in the code 
with constants from already imported packages.

### Installation

```shell
go get -u github.com/AntonBraer/constantcheck@latest
```

### Usage

```
./constantcheck ./...
```

### Examples

Let's say we have a project structure like this:
<table>
<tr>
<td> Structure </td> <td> main.go </td> <td> helps.go </td>
</tr>
<tr>
<td> 

```bash   
├── helps
│   ├── helps.go
├──main.go
``` 
</td>
<td>

```go
package main

import (
	"fmt"

	"example/helps"
)

func main() {
	if err := helps.Ping(); err != nil{
		fmt.Println("myImportantError")	
    }
}
```

</td>
<td>

```go
package helps

import (
	"fmt"
)

const CustomError = "myImportantError"

func Ping() error{
	fmt.Println("Pong")
	return nil
}
```
</td>
</tr>
</table>

After checking the linter, we should get an error that we can use a constant from 
`helps` package with name `CustomError`

---
I couldn't show it in tests because there is some problem with imports in testing. But the linter parses ALL imported packages of the file, not just standard ones
