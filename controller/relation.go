package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/mydb"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		to_user_id := c.Query("to_user_id")
		action_type := c.Query("action_type")

		//判断不能关注自己
		userid, _ := strconv.ParseInt(to_user_id, 10, 64)
		if user.Id == userid {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "不能关注自己哦！"})
			return
		}
		err := mydb.Db.QueryRow("select * from follow where user_id = ? and to_user_id = ?", user.Id, to_user_id)
		if action_type == "1" {
			if err == nil {
				_, err := mydb.Db.Exec("update follow set state = 1 where user_id = ? and to_user_id = ?", user.Id, to_user_id)
				if err != nil {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Follow failed 1"})
					return
				}
			} else {
				_, err := mydb.Db.Exec("insert into follow (user_id, to_user_id) value(?, ?)", user.Id, to_user_id)
				if err != nil {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Follow failed 2"})
					return
				}
			}
			mydb.Db.Exec("update user set following_count = following_count + 1 where user_id = ?", user.Id)
			mydb.Db.Exec("update user set follower_count = follower_count + 1 where user_id = ?", to_user_id)

		} else if action_type == "2" {
			_, err := mydb.Db.Exec("update follow set state = 0 where user_id = ? && to_user_id = ?", user.Id, to_user_id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unfollow failed 1"})
				return
			}

			_, err = mydb.Db.Exec("update user set following_count = following_count - 1 where user_id = ?", user.Id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unfollow failed 2"})
				return
			}
			_, err = mydb.Db.Exec("update user set follower_count = follower_count - 1 where user_id = ?", to_user_id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unfollow failed 3"})
				return
			}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 关注的所有用户列表
func FollowList(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		user_id := c.Query("user_id")
		//查找存在关注关系，且state=1（没有取关）
		sqlStr := "select to_user_id from follow where user_id = ? && state = 1"
		rows, err := mydb.Db.Query(sqlStr, user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Query failed"})
			return
		}
		var users []User

		for rows.Next() {
			var u User
			err = rows.Scan(&u.Id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Scan failed"})
				return
			}
			user_sqlStr := "select username, following_count, follower_count from user where user_id=?"
			err := mydb.Db.QueryRow(user_sqlStr, u.Id).Scan(&u.Name, &u.FollowCount, &u.FollowerCount)

			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Scan failed"})
				return
			}

			//查询时已经限制state=1，这里直接赋值true就好
			u.IsFollow = true

			users = append(users, u)
		}

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: users,
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 粉丝列表
func FollowerList(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		user_id := c.Query("user_id")

		sqlStr := "select user_id from follow where to_user_id = ? && state = 1"
		rows, err := mydb.Db.Query(sqlStr, user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Query failed"})
			return
		}
		var users []User

		for rows.Next() {
			var u User
			err = rows.Scan(&u.Id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video Scan failed"})
				return
			}
			user_sqlStr := "select username, following_count, follower_count from user where user_id=?"
			err := mydb.Db.QueryRow(user_sqlStr, u.Id).Scan(&u.Name, &u.FollowCount, &u.FollowerCount)

			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Scan failed"})
				return
			}

			is_Follow := false
			_ = mydb.Db.QueryRow("select state from follow where user_id=? && to_user_id=?", user.Id, u.Id).Scan(&is_Follow)
			if is_Follow {
				u.IsFollow = true
			}

			users = append(users, u)
		}

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: users,
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func FriendList(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		user_id := c.Query("user_id")

		sqlStr := "select user_id from follow where to_user_id = ? and state = 1"
		rows, err := mydb.Db.Query(sqlStr, user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Query failed"})
			return
		}
		var users []User

		for rows.Next() {
			var u User
			var id int64
			err = rows.Scan(&id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video Scan failed"})
				return
			}

			_, err = mydb.Db.Query("select * from follow where user_id = ? and to_user_id = ? and state = 1", user_id, id)
			if err != nil {
				continue
			}
			u.Id = id
			user_sqlStr := "select username, following_count, follower_count from user where user_id=?"
			err := mydb.Db.QueryRow(user_sqlStr, u.Id).Scan(&u.Name, &u.FollowCount, &u.FollowerCount)

			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Scan failed"})
				return
			}

			is_Follow := false
			_ = mydb.Db.QueryRow("select state from follow where user_id=? && to_user_id=?", user.Id, u.Id).Scan(&is_Follow)
			if is_Follow {
				u.IsFollow = true
			}

			users = append(users, u)
		}

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: users,
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}
