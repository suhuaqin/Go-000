package main

import (
	"fmt"

	"github.com/pkg/errors"
)

/*
	操作数据库时，如果遇到一个 <not found> 错误
	可以和其他错误一样，Wrap抛给上层，让交给上层决定是否对<not found>错误做处理
*/
func main() {
	if err := biz(); err != nil {
		fmt.Printf("%+v", err)
		if errors.Is(err, ErrNotFound) {
			// 如果要根据特定的错误做特定的处理，使用errors.Is判断错误
			// ...
		} else {
			return
		}
	}

	//...
}

func biz() error {
	userId := 10086
	if _, err := Dao(userId); err != nil {
		return err
	}

	// ...
	return nil
}

type User struct {
	ID int
	// ...
}

var ErrNotFound = errors.New("Not Found")

func Dao(id int) (*User, error) {
	var user *User
	// ... 根据id查询数据库
	err := ErrNotFound
	if err != nil {
		return user, errors.Wrapf(err, "id: %v", id)
	}
	return user, nil
}
