package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/RaymondCode/simple-demo/mydb"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			video_id := c.Query("video_id")
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			result, err := mydb.Db.Exec("INSERT INTO comment (user_id, play_id, content) VALUE(?,?,?)", user.Id, video_id, text)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Comment failed"})
				return
			}
			_, err = mydb.Db.Exec("update video set comment_count = comment_count + 1 where play_id = ?", video_id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Something Wrong"})
				return
			}

			id, err := result.LastInsertId()
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Something Wrong"})
				return
			}

			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id:         id,
					User:       user,
					Content:    text,
					CreateDate: timeStr,
				}})
			return
		} else if actionType == "2" {
			comment_id := c.Query("comment_id")
			video_id := c.Query("video_id")
			_, err := mydb.Db.Exec("update comment set state = 0 where comment_id = ?", comment_id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Delete failed"})
				return
			}
			_, err = mydb.Db.Exec("update video set comment_count = comment_count - 1 where play_id = ?", video_id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Something Wrong"})
				return
			}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	var comments []Comment
	video_id := c.Query("video_id")
	sqlStr := "select comment_id, user_id, content, date from comment where play_id = ? && state = 1 order by date desc"
	rows, err := mydb.Db.Query(sqlStr, video_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Query failed"})
		return
	}
	
	for rows.Next() {
		var comment Comment
		var u User
		err = rows.Scan(&comment.Id, &u.Id, &comment.Content, &comment.CreateDate)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Comment Scan failed"})
			return 
		}
		user_sqlStr := "select username, following_count, follower_count from user where user_id = ?"
		err := mydb.Db.QueryRow(user_sqlStr, u.Id).Scan(&u.Name, &u.FollowCount, &u.FollowerCount)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Scan failed"})
			return 
		}
		comment.User = u
		comments = append(comments, comment)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
