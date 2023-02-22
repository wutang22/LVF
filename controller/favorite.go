package controller

import (
	"fmt"
	"net/http"

	"github.com/RaymondCode/simple-demo/mydb"
	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		video_id := c.Query("video_id")
		action_type := c.Query("action_type")
		err := mydb.Db.QueryRow("select user_id from favorite where user_id = ? and video_id = ?", user.Id, video_id)
		if action_type == "1" {
			if err != nil {
				fmt.Println(err)
				_, err := mydb.Db.Exec("insert into favorite (user_id, video_id) value(?, ?)", user.Id, video_id)
				if err != nil {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Favorite failed 2"})
					return
				}
			} else {
				_, err := mydb.Db.Exec("update favorite set state = 1 where user_id = ? and video_id = ?", user.Id, video_id)
				if err != nil {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "favorite failed 1"})
					return
				}
			}
			mydb.Db.Exec("update video set favorite_count = favorite_count  + 1 where play_id = ?", video_id)

		} else if action_type == "2" {
			_, err := mydb.Db.Exec("update favorite set state = 0 where user_id = ? and video_id = ?", user.Id, video_id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unfavorite failed 1"})
				return
			}

			_, err = mydb.Db.Exec("update video set favorite_count = favorite_count - 1 where play_id = ?", video_id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unfavorite failed 2"})
				return
			}
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Favorite success!"})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		user_id := c.Query("user_id")
		//查到所有点赞过的video id
		sqlStr := "select video_id from favorite where user_id = ? && state = 1"
		rows, err := mydb.Db.Query(sqlStr, user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Query failed"})
			return
		}

		var videos []Video
		//查询每个video信息
		for rows.Next() {
			var v Video
			var u User
			err = rows.Scan(&v.Id)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video Scan failed"})
				return
			}
			video_sqlStr := "select user_id, play_url, cover_url,title,favorite_count,comment_count from video where play_id=? && play_status = 0"
			err := mydb.Db.QueryRow(video_sqlStr, v.Id).Scan(&u.Id, &v.PlayUrl, &v.CoverUrl, &v.Title, &v.FavoriteCount, &v.CommentCount)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video Scan failed"})
				return
			}
			// fmt.Println(v.Author.Id)
			user_sqlStr := "select username, following_count, follower_count from user where user_id=?"
			err = mydb.Db.QueryRow(user_sqlStr, u.Id).Scan(&u.Name, &u.FollowCount, &u.FollowerCount)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Scan failed"})
				return
			}
			//喜欢列表返回的都是点赞过的视频，直接true就好
			v.IsFavorite = true
			v.Author = u

			videos = append(videos, v)
		}

		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: videos,
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}
