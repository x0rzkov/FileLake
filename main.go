package main
import (
	"FileLake/common/config"
	//_ "FileLake/docs"
	"FileLake/router"
	"log"
)
// @title Gin swagger
// @version 1.0
// @description Gin swagger 示例项目
// @contact.name qihang666
// @contact.url https://github.com/qihang666
// @contact.email 977767937@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9000
//CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
func main() {
	log.Println("init gin server")
	rou:= router.InitRouter()
	rou.Run(":" + config.Config.GetString("port")) // listen and serve on 0.0.0.0:8080
}