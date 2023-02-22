package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/mydb"
	"github.com/gin-gonic/gin"
)

var tempChat = map[string][]Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	if user, exist := usersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB))

		result, err := mydb.Db.Exec("insert into message (from_user_id, to_user_id, content, chat_key) value (?, ?, ?, ?)", user.Id, toUserId, content, chatKey)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "消息发送失败！"},
			})
			return
		}
		message_id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Something Wrong"},
			})
			return
		}

		curMessage := Message{
			Id:         message_id,
			ToUserId:   int64(userIdB),
			FromUserId: user.Id,
			Content:    content,
			//CreateTime: time.Now().Format(time.Kitchen),
		}

		if messages, exist := tempChat[chatKey]; exist {
			tempChat[chatKey] = append(messages, curMessage)
		} else {
			tempChat[chatKey] = []Message{curMessage}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	//pre_msg_time := c.Query("pre_msg_time")

	if user, exist := usersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB))
		rows, err := mydb.Db.Query("select message_id, from_user_id, to_user_id, content from message where chat_key = ?", chatKey)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Query failed"})
			return
		}
		var messages []Message
		for rows.Next() {
			var m Message
			err = rows.Scan(&m.Id, &m.ToUserId, &m.FromUserId, &m.Content)

			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Author Scan failed"})
				return
			}
			messages = append(messages, m)
		}

		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: messages})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
