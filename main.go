package main

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"github.com/RaymondCode/simple-demo/mydb"
	"net/http"
)

func main() {
	go service.RunMessageServer()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.StaticFS("/public", http.Dir("./public"))

	mydb.Init_db()
	
	initRouter(r)

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
