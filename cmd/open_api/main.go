package main

import (
	"flag"
	"fmt"
	_ "micro_servers/cmd/open_api/docs"
	apiAuth "micro_servers/internal/api/auth"
	"micro_servers/internal/api/manage"
	"micro_servers/internal/api/user"
	"micro_servers/pkg/common/config"
	"micro_servers/pkg/common/log"
	"micro_servers/pkg/utils"

	//_ "github.com/razeencheng/demo-go/swaggo-gin/docs"
	"io"
	"os"
	"strconv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	//"syscall"
	"micro_servers/pkg/common/constant"
)

// @title open-Server API
// @version 1.0
// @description  open-Server 的API服务器文档, 文档中所有请求都有一个operationID字段用于链路追踪

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("../logs/api.log")
	gin.DefaultWriter = io.MultiWriter(f)
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(utils.CorsHandler())

	log.Info("load  config: ", config.Config)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// user routing group, which handles user registration and login services
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/update_user_info", user.UpdateUserInfo) //1
		userRouterGroup.POST("/set_global_msg_recv_opt", user.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_users_info", user.GetUsersInfo)                  //1
		userRouterGroup.POST("/get_self_user_info", user.GetSelfUserInfo)           //1
		userRouterGroup.POST("/get_users_online_status", user.GetUsersOnlineStatus) //1
		userRouterGroup.POST("/get_users_info_from_cache", user.GetUsersInfoFromCache)
		userRouterGroup.POST("/get_user_friend_from_cache", user.GetFriendIDListFromCache)
		userRouterGroup.POST("/get_black_list_from_cache", user.GetBlackIDListFromCache)
		userRouterGroup.POST("/get_all_users_uid", manage.GetAllUsersUid) //1
		userRouterGroup.POST("/account_check", manage.AccountCheck)       //1
		//	userRouterGroup.POST("/get_users_online_status", manage.GetUsersOnlineStatus) //1
	}
	//certificate
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_register", apiAuth.UserRegister) //1
		authRouterGroup.POST("/user_token", apiAuth.UserToken)       //1
		authRouterGroup.POST("/parse_token", apiAuth.ParseToken)     //1
		authRouterGroup.POST("/force_logout", apiAuth.ForceLogout)   //1
	}
	// //Third service
	// thirdGroup := r.Group("/third")
	// {
	// 	thirdGroup.POST("/tencent_cloud_storage_credential", apiThird.TencentCloudStorageCredential)
	// 	thirdGroup.POST("/ali_oss_credential", apiThird.AliOSSCredential)
	// 	thirdGroup.POST("/minio_storage_credential", apiThird.MinioStorageCredential)
	// 	thirdGroup.POST("/minio_upload", apiThird.MinioUploadFile)
	// 	thirdGroup.POST("/upload_update_app", apiThird.UploadUpdateApp)
	// 	thirdGroup.POST("/get_download_url", apiThird.GetDownloadURL)
	// 	thirdGroup.POST("/get_rtc_invitation_info", apiThird.GetRTCInvitationInfo)
	// 	thirdGroup.POST("/get_rtc_invitation_start_app", apiThird.GetRTCInvitationInfoStartApp)
	// 	thirdGroup.POST("/fcm_update_token", apiThird.FcmUpdateToken)
	// 	thirdGroup.POST("/aws_storage_credential", apiThird.AwsStorageCredential)
	// }
	// go apiThird.MinioInit()
	defaultPorts := config.Config.Api.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10002 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start api server, address: ", address)
	err := r.Run(address)
	if err != nil {
		log.Error("", "run failed ", *ginPort, err.Error())
	}
}
