package route

import (
	"os"

	v1 "github.com/IceWhaleTech/CasaOS-AppManagement/route/v1"
	"github.com/IceWhaleTech/CasaOS-Common/middleware"
	"github.com/IceWhaleTech/CasaOS-Common/utils/jwt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// check if environment variable is set
	ginMode, success := os.LookupEnv(gin.EnvGinMode)
	if !success {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	if ginMode != gin.ReleaseMode {
		r.Use(middleware.WriteLog())
	}

	v1Group := r.Group("/v1")

	v1Group.Use(jwt.ExceptLocalhost())
	{
		v1AppsGroup := v1Group.Group("/apps")
		v1AppsGroup.Use()
		{
			v1AppsGroup.GET("", v1.AppList) // list
			v1AppsGroup.GET("/:id", v1.AppInfo)
		}

		v1ContainerGroup := v1Group.Group("/container")
		v1ContainerGroup.Use()
		{

			v1ContainerGroup.GET("", v1.MyAppList) ///my/list
			v1ContainerGroup.GET("/usage", v1.AppUsageList)
			v1ContainerGroup.GET("/:id", v1.ContainerUpdateInfo)    ///update/:id/info
			v1ContainerGroup.GET("/:id/logs", v1.ContainerLog)      // /app/logs/:id
			v1ContainerGroup.GET("/networks", v1.GetDockerNetworks) // app/install/config

			v1ContainerGroup.GET("/:id/state", v1.GetContainerState) // app/state/:id ?state=install_progress
			// there are problems, temporarily do not deal with
			v1ContainerGroup.GET("/:id/terminal", v1.DockerTerminal) // app/terminal/:id
			v1ContainerGroup.POST("", v1.InstallApp)                 // app/install

			v1ContainerGroup.PUT("/:id", v1.UpdateSetting) ///update/:id/setting

			v1ContainerGroup.PUT("/:id/state", v1.ChangAppState) // /app/state/:id
			v1ContainerGroup.DELETE("/:id", v1.UnInstallApp)     // app/uninstall/:id
			// Not used
			v1ContainerGroup.PUT("/:id/latest", v1.PutAppUpdate)
			// Not used
			v1ContainerGroup.POST("/share", v1.ShareAppFile)
			v1ContainerGroup.GET("/info", v1.GetDockerDaemonConfiguration)
			v1ContainerGroup.PUT("/info", v1.PutDockerDaemonConfiguration)

		}
		v1AppCategoriesGroup := v1Group.Group("/app-categories")
		v1AppCategoriesGroup.Use()
		{
			v1AppCategoriesGroup.GET("", v1.CategoryList)
		}

	}

	return r
}
