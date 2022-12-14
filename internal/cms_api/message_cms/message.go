package messageCMS

import (
	"context"
	"micro_servers/pkg/cms_api_struct"
	"micro_servers/pkg/common/config"
	openIMHttp "micro_servers/pkg/common/http"
	"micro_servers/pkg/common/log"
	"micro_servers/pkg/grpc-etcdv3/getcdv3"
	pbMessage "micro_servers/pkg/proto/message_cms"
	pbCommon "micro_servers/pkg/proto/sdk_ws"
	"micro_servers/pkg/utils"
	"net/http"
	"strings"

	"micro_servers/pkg/common/constant"

	"github.com/gin-gonic/gin"
)

func BroadcastMessage(c *gin.Context) {
	var (
		reqPb pbMessage.BoradcastMessageReq
	)
	reqPb.OperationID = utils.OperationIDGenerator()
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMessageCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbMessage.NewMessageCMSClient(etcdConn)
	_, err := client.BoradcastMessage(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetChatLogs rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func MassSendMassage(c *gin.Context) {
	var (
		reqPb pbMessage.MassSendMessageReq
	)
	reqPb.OperationID = utils.OperationIDGenerator()
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMessageCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbMessage.NewMessageCMSClient(etcdConn)
	_, err := client.MassSendMessage(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetChatLogs rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func WithdrawMessage(c *gin.Context) {
	var (
		reqPb pbMessage.WithdrawMessageReq
	)
	reqPb.OperationID = utils.OperationIDGenerator()
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMessageCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbMessage.NewMessageCMSClient(etcdConn)
	_, err := client.WithdrawMessage(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetChatLogs rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func GetChatLogs(c *gin.Context) {
	var (
		req   cms_api_struct.GetChatLogsRequest
		resp  cms_api_struct.GetChatLogsResponse
		reqPb pbMessage.GetChatLogsReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	reqPb.Pagination = &pbCommon.RequestPagination{
		PageNumber: int32(req.PageNumber),
		ShowNumber: int32(req.ShowNumber),
	}
	utils.CopyStructFields(&reqPb, &req)
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMessageCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbMessage.NewMessageCMSClient(etcdConn)
	respPb, err := client.GetChatLogs(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetChatLogs rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, resp)
		return
	}
	//utils.CopyStructFields(&resp, &respPb)
	for _, chatLog := range respPb.ChatLogs {
		resp.ChatLogs = append(resp.ChatLogs, cms_api_struct.ChatLog{
			SessionType:      int(chatLog.SessionType),
			ContentType:      int(chatLog.ContentType),
			SenderNickName:   chatLog.SenderNickName,
			SenderId:         chatLog.SenderId,
			SearchContent:    chatLog.SearchContent,
			WholeContent:     chatLog.WholeContent,
			ReceiverNickName: chatLog.ReciverNickName,
			ReceiverID:       chatLog.ReciverId,
			GroupName:        chatLog.GroupName,
			GroupId:          chatLog.GroupId,
			Date:             chatLog.Date,
		})
	}
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.ChatLogsNum = int(respPb.ChatLogsNum)
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "resp", resp)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}
