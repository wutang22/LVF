package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	//"time"
	"github.com/RaymondCode/simple-demo/mydb"
	"net/url"
)


type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	title := c.PostForm("title")
	token := c.PostForm("token")
		
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	//filename = time.Now().Format("2023-01-01-00:00:00") + filename
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	play_url, _ := url.JoinPath(public_url, finalName)
	default_coverurl := "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
	fmt.Println(saveFile)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	_, err = mydb.Db.Exec("INSERT INTO video (user_id, play_url, cover_url, title) VALUE(?,?,?,?)", user.Id, play_url, default_coverurl, title)

	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Update failed"},
		})
		return
	}
	_, err = mydb.Db.Exec("update user set play_count = play_count + 1 where user_id = ?", user.Id)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Update failed 1"})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if user, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		fmt.Println(usersLoginInfo[token])
		return
	} else {
	
		sqlStr := "select play_id, play_url, cover_url, title, favorite_count, comment_count from video where user_id = ?"
		rows, err := mydb.Db.Query(sqlStr, user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Something wrong"})
			return
		} 

		var videos []Video

		for rows.Next() {
			var v Video
			var u User
			err = rows.Scan(&v.Id, &v.PlayUrl, &v.CoverUrl, &v.Title, &v.FavoriteCount, &v.CommentCount)
			u = usersLoginInfo[token]
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Scan failed"})
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
			v.Author = u
			videos = append(videos, v)
		}

		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: videos,
		})
	}
}
