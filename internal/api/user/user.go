package user

import (
	"context"
	jsonData "micro_servers/internal/utils"
	api "micro_servers/pkg/base_info"
	"micro_servers/pkg/common/config"
	"micro_servers/pkg/common/log"
	"micro_servers/pkg/common/token_verify"
	"micro_servers/pkg/grpc-etcdv3/getcdv3"
	open_im_sdk "micro_servers/pkg/proto/sdk_ws"
	rpc "micro_servers/pkg/proto/user"
	"micro_servers/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUsersInfoFromCache(c *gin.Context) {
	params := api.GetUsersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "GetUsersInfoFromCache req: ", params)
	req := &rpc.GetUserInfoReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	var publicUserInfoList []*open_im_sdk.PublicUserInfo
	for _, v := range RpcResp.UserInfoList {
		publicUserInfoList = append(publicUserInfoList,
			&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex})
	}
	resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
	resp.Data = jsonData.JsonDataList(resp.UserInfoList)
	log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取用户信息
// @Description 根据用户列表批量获取用户信息
// @Tags 用户相关
// @ID GetUsersInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUsersInfoReq true "请求体"
// @Produce json
// @Success 0 {object} api.GetUsersInfoResp{Data=[]open_im_sdk.PublicUserInfo}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/get_users_info [post]
func GetUsersInfo(c *gin.Context) {
	params := api.GetUsersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetUserInfoReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	var publicUserInfoList []*open_im_sdk.PublicUserInfo
	for _, v := range RpcResp.UserInfoList {
		publicUserInfoList = append(publicUserInfoList,
			&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex})
	}

	resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
	resp.Data = jsonData.JsonDataList(resp.UserInfoList)
	log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 修改用户信息
// @Description 修改用户信息 userID faceURL等
// @Tags 用户相关
// @ID UpdateUserInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.UpdateSelfUserInfoReq true "请求体"
// @Produce json
// @Success 0 {object} api.UpdateUserInfoResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/update_user_info [post]
func UpdateUserInfo(c *gin.Context) {
	params := api.UpdateSelfUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.UpdateUserInfoReq{UserInfo: &open_im_sdk.UserInfo{}}
	utils.CopyStructFields(req.UserInfo, &params)
	req.OperationID = params.OperationID
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "UpdateUserInfo args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.UpdateUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UpdateUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "UpdateUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取自己的信息
// @Description 传入ID获取自己的信息
// @Tags 用户相关
// @ID GetSelfUserInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetSelfUserInfoReq true "请求体"
// @Produce json
// @Success 0 {object} api.GetSelfUserInfoResp{data=open_im_sdk.UserInfo}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/get_self_user_info [post]
func GetSelfUserInfo(c *gin.Context) {
	params := api.GetSelfUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := "BindJSON failed " + err.Error()
		log.NewError("0", "BindJSON failed ", errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	req := &rpc.GetUserInfoReq{}

	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	req.UserIDList = append(req.UserIDList, req.OpUserID)
	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	if len(RpcResp.UserInfoList) == 1 {
		resp := api.GetSelfUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfo: RpcResp.UserInfoList[0]}
		resp.Data = jsonData.JsonDataOne(resp.UserInfo)
		log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
		c.JSON(http.StatusOK, resp)
	} else {
		resp := api.GetSelfUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
		log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
		c.JSON(http.StatusOK, resp)
	}

}
