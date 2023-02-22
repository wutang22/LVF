package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	//"sync/atomic"
	"fmt"

	"github.com/RaymondCode/simple-demo/mydb"
	"golang.org/x/crypto/bcrypt"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"123123456": {
		Id:            1,
		Name:          "123",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
	"test7123456": {
		Id:            7,
		Name:          "test7",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}
var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	passwordByte := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	password = string(hashedPassword)

	token := username + "loginOK"

	result, err := mydb.Db.Exec("INSERT INTO user (username, password) VALUE(?,?)", username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "somthing wrong"},
			})
		} else {
			u := User{
				Id:   id,
				Name: username,
			}
			usersLoginInfo[token] = u

			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   userIdSequence,
				Token:    username + password,
			})
		}
	}
	// if _, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	// 	})
	// } else {
	// 	atomic.AddInt64(&userIdSequence, 1)
	// 	newUser := User{
	// 		Id:   userIdSequence,
	// 		Name: username,
	// 	}
	// 	usersLoginInfo[token] = newUser
	// 	c.JSON(http.StatusOK, UserLoginResponse{
	// 		Response: Response{StatusCode: 0},
	// 		UserId:   userIdSequence,
	// 		Token:    username + password,
	// 	})
	// }
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + "loginOK"
	hashedPassword := ""
	var u User
	sqlStr := "select user_id, username, password,following_count, follower_count from user where username=?"
	err := mydb.Db.QueryRow(sqlStr, username).Scan(&u.Id, &u.Name, &hashedPassword, &u.FollowCount, &u.FollowerCount)

	usersLoginInfo[token] = u

	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Username error"},
		})
	} else {
		fmt.Println(usersLoginInfo[token])
		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Password error"},
			})
		}
		if user, exist := usersLoginInfo[token]; exist {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   user.Id,
				Token:    token,
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		}

		// if password != u.password {
		// 	c.JSON(http.StatusOK, UserLoginResponse{
		// 		Response: Response{StatusCode: 1, StatusMsg: "用户名或密码错误！"},
		// 	})
		// } else {
		// 	c.JSON(http.StatusOK, UserLoginResponse{
		// 		Response: Response{StatusCode: 0},
		// 		UserId:   u.id,
		// 		Token:    token,
		// 	})
		// }
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	if user, exist := usersLoginInfo[token]; exist {
		user_sqlStr := "select username, following_count, follower_count from user where user_id=?"
		err := mydb.Db.QueryRow(user_sqlStr, user.Id).Scan(&user.Name, &user.FollowCount, &user.FollowerCount)

		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Scan failed"})
			return
		}
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
