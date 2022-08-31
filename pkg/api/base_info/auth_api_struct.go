package base_info

type UserRegisterReq struct {
	Secret   string `json:"secret" binding:"required,max=32"`
	Platform int32  `json:"platform" binding:"required,min=1,max=7"`
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenInfo struct {
	UserID      string `json:"userID"`
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}
type UserRegisterResp struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}

type ApiUserInfo struct {
	UserID      string `json:"userID" binding:"required,min=1,max=64" swaggo:"true,用户ID,"`
	Nickname    string `json:"nickname" binding:"omitempty,min=1,max=64" swaggo:"true,my id,19"`
	FaceURL     string `json:"faceURL" binding:"omitempty,max=1024"`
	Gender      int32  `json:"gender" binding:"omitempty,oneof=0 1 2"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=32"`
	Birth       uint32 `json:"birth" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,max=64"`
	Ex          string `json:"ex" binding:"omitempty,max=1024"`
}
