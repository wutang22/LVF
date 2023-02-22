package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"github.com/RaymondCode/simple-demo/mydb"
	//"fmt"
	//"strconv"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//latest_time := c.Query("latest_time")
	
	// time_int, _ := strconv.ParseInt(latest_time, 10, 64)
	// timeStr := time_int.Format("2006-01-02 15:04:05")
	sqlStr := "select play_id, user_id, play_url, cover_url, title, favorite_count, comment_count from video order by date desc LIMIT 30"
	rows, err := mydb.Db.Query(sqlStr)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Query failed"})
		return
	}
	var videos []Video

	for rows.Next() {
		var v Video
		var u User
		err = rows.Scan(&v.Id, &u.Id, &v.PlayUrl, &v.CoverUrl, &v.Title, &v.FavoriteCount, &v.CommentCount)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video Scan failed"})
			return 
		}
		user_sqlStr := "select username, following_count, follower_count from user where user_id=?"
		err := mydb.Db.QueryRow(user_sqlStr, u.Id).Scan(&u.Name, &u.FollowCount, &u.FollowerCount)

		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Author Scan failed"})
			return 
		}
		
		v.Author = u
		videos = append(videos, v)
	}
	//fmt.Println(videos)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
