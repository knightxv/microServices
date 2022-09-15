/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/15 10:28).
 */
package manage

import (
	"context"
	api "micro_servers/pkg/base_info"
	"micro_servers/pkg/common/config"
	"micro_servers/pkg/common/log"
	"micro_servers/pkg/common/token_verify"
	"micro_servers/pkg/grpc-etcdv3/getcdv3"
	rpc "micro_servers/pkg/proto/user"
	"micro_servers/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	params := api.DeleteUsersReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.DeleteUsersReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(params.OperationID, "DeleteUser args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)

	RpcResp, err := client.DeleteUsers(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call delete users rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call delete users rpc server failed"})
		return
	}
	resp := api.DeleteUsersResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, FailedUserIDList: RpcResp.FailedUserIDList}
	if len(RpcResp.FailedUserIDList) == 0 {
		resp.FailedUserIDList = []string{}
	}
	log.NewInfo(req.OperationID, "DeleteUser api return", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取所有用户uid列表
// @Description 获取所有用户uid列表
// @Tags 用户相关
// @ID GetAllUsersUid
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetAllUsersUidReq true "请求体"
// @Produce json
// @Success 0 {object} api.GetAllUsersUidResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/get_all_users_uid [post]
func GetAllUsersUid(c *gin.Context) {
	params := api.GetAllUsersUidReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetAllUserIDReq{}
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

	log.NewInfo(params.OperationID, "GetAllUsersUid args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetAllUserID(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call GetAllUsersUid users rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call GetAllUsersUid users rpc server failed"})
		return
	}
	resp := api.GetAllUsersUidResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserIDList: RpcResp.UserIDList}
	if len(RpcResp.UserIDList) == 0 {
		resp.UserIDList = []string{}
	}
	log.NewInfo(req.OperationID, "GetAllUsersUid api return", resp)
	c.JSON(http.StatusOK, resp)

}

// @Summary 检查列表账户注册状态，并且返回结果
// @Description 传入UserIDList检查列表账户注册状态，并且返回结果
// @Tags 用户相关
// @ID AccountCheck
// @Accept json
// @Param token header string true "im token"
// @Param req body api.AccountCheckReq true "请求体"
// @Produce json
// @Success 0 {object} api.AccountCheckResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/account_check [post]
func AccountCheck(c *gin.Context) {
	params := api.AccountCheckReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AccountCheckReq{}
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

	log.NewInfo(params.OperationID, "AccountCheck args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)

	RpcResp, err := client.AccountCheck(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call AccountCheck users rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call AccountCheck users rpc server failed"})
		return
	}
	resp := api.AccountCheckResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, ResultList: RpcResp.ResultList}
	if len(RpcResp.ResultList) == 0 {
		resp.ResultList = []*rpc.AccountCheckResp_SingleUserStatus{}
	}
	log.NewInfo(req.OperationID, "AccountCheck api return", resp)
	c.JSON(http.StatusOK, resp)
}
