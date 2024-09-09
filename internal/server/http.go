package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	apiV1 "go-chat/api/v1"
	"go-chat/docs"
	"go-chat/internal/handler"
	"go-chat/internal/middleware"
	"go-chat/pkg/jwt"
	"go-chat/pkg/log"
	"go-chat/pkg/server/http"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	uploadHandler *handler.UploadHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	v1 := s.Group("/v1")

	user := v1.Group("/user")
	{
		user.POST("/register", userHandler.Register)
		user.POST("/register_email_code_check", userHandler.CheckRegisterEmailCode)
		user.POST("/login", userHandler.Login)
		user.POST("email_login_code_check", userHandler.EmailLoginCodeCheck)
		user.POST("/email_login", userHandler.EmailLogin)
		auth := user.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			auth.POST("/user_info_update", userHandler.UserInfoUpdate)
		}
	}

	upload := v1.Group("/upload").Use(middleware.StrictAuth(jwt, logger))
	{
		upload.POST("/file", uploadHandler.Upload)
	}

	return s
}
