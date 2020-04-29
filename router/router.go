package router
import (
	"FileLake/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)
func InitRouter() *gin.Engine {
	r := gin.Default()
	//swagger 配置相关
	url := ginSwagger.URL("http://localhost:9000/swagger/doc.json")
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//用户模块路由
	r.GET("/getuserfile/:servicelink",service.GetFileByLink)
	r.POST("/add/:accountaddress",service.Add)
	r.PUT("/user",service.Update)
	r.DELETE("/user",service.Delete)
	r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK,"hello")
	})
	return r
}
