gorm-pagination

### INSTALL
```
go get github.com/xsephiroth/gorm-pagination
```

### TEST
```
make test
```

### USAGE
example/example.go
```
package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	pagination "github.com/xsephiroth/gorm-pagination"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Active   bool   `json:"-"`
}

func main() {
	db, err := gorm.Open("sqlite3", "/tmp/test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.CreateTable(&User{})

	tx := db.Model(&User{}).Where("active = ?", true)
	users := make([]User, 0)
	resp := pagination.PagePagination(tx, 2, 10, &users)
	fmt.Printf("%#v\n", resp)
}
```
