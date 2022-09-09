package token_verify

import (
	"micro_servers/pkg/common/config"
	"micro_servers/pkg/common/constant"
	"micro_servers/pkg/common/log"
	"micro_servers/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//var (
//	TokenExpired     = errors.New("token is timed out, please log in again")
//	TokenInvalid     = errors.New("token has been invalidated")
//	TokenNotValidYet = errors.New("token not active yet")
//	TokenMalformed   = errors.New("that's not even a token")
//	TokenUnknown     = errors.New("couldn't handle this token")
//)

type Claims struct {
	UID string
	jwt.RegisteredClaims
}

func BuildClaims(uid string, ttl int64) Claims {
	now := time.Now()
	return Claims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*24) * time.Hour)), //Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                        //Issuing time
			NotBefore: jwt.NewNumericDate(now),                                        //Begin Effective time
		}}
}

func CreateToken(userID string) (string, int64, error) {
	claims := BuildClaims(userID, config.Config.TokenPolicy.AccessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.TokenPolicy.AccessSecret))
	if err != nil {
		return "", 0, err
	}
	return tokenString, claims.ExpiresAt.Time.Unix(), err
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenPolicy.AccessSecret), nil
	}
}

func GetClaimFromToken(tokensString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, constant.ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, constant.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, constant.ErrTokenNotValidYet
			} else {
				return nil, constant.ErrTokenUnknown
			}
		} else {
			return nil, constant.ErrTokenNotValidYet
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			//log.NewDebug("", claims.UID, claims.Platform)
			return claims, nil
		}
		return nil, constant.ErrTokenNotValidYet
	}
}

func GetUserIDFromToken(token string, operationID string) (bool, string, string) {
	claims, err := ParseToken(token, operationID)
	if err != nil {
		log.Error(operationID, "ParseToken failed, ", err.Error(), token)
		return false, "", err.Error()
	}
	log.Debug(operationID, "token claims.ExpiresAt.Second() ", claims.ExpiresAt.Unix())
	return true, claims.UID, ""
}

func GetUserIDFromTokenExpireTime(token string, operationID string) (bool, string, string, int64) {
	claims, err := ParseToken(token, operationID)
	if err != nil {
		log.Error(operationID, "ParseToken failed, ", err.Error(), token)
		return false, "", err.Error(), 0
	}
	return true, claims.UID, "", claims.ExpiresAt.Unix()
}

func ParseTokenGetUserID(token string, operationID string) (string, error) {
	claims, err := ParseToken(token, operationID)
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return claims.UID, nil
}

func ParseToken(tokensString, _ string) (claims *Claims, err error) {
	claims, err = GetClaimFromToken(tokensString)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return claims, nil
}

//func MakeTheTokenInvalid(currentClaims *Claims, platformClass string) (bool, error) {
//	storedRedisTokenInterface, err := db.DB.GetPlatformToken(currentClaims.UID, platformClass)
//	if err != nil {
//		return false, err
//	}
//	storedRedisPlatformClaims, err := ParseRedisInterfaceToken(storedRedisTokenInterface)
//	if err != nil {
//		return false, err
//	}
//	//if issue time less than redis token then make this token invalid
//	if currentClaims.IssuedAt.Time.Unix() < storedRedisPlatformClaims.IssuedAt.Time.Unix() {
//		return true, constant.TokenInvalid
//	}
//	return false, nil
//}

func ParseRedisInterfaceToken(redisToken interface{}) (*Claims, error) {
	return GetClaimFromToken(string(redisToken.([]uint8)))
}

//Validation token, false means failure, true means successful verification
func VerifyToken(token, uid string) (bool, error) {
	claims, err := ParseToken(token, "")
	if err != nil {
		return false, utils.Wrap(err, "ParseToken failed")
	}
	if claims.UID != uid {
		return false, &constant.ErrTokenUnknown
	}

	log.NewDebug("VerifyToken", claims.UID)
	return true, nil
}

func WsVerifyToken(token, uid string, operationID string) (bool, error, string) {
	argMsg := "args: token: " + token + " operationID: " + operationID + " userID: " + uid
	claims, err := ParseToken(token, operationID)
	if err != nil {
		//if errors.Is(err, constant.ErrTokenUnknown) {
		//	errMsg := "ParseToken failed ErrTokenUnknown " + err.Error()
		//	log.Error(operationID, errMsg)
		//}
		//e := errors.Unwrap(err)
		//if errors.Is(e, constant.ErrTokenUnknown) {
		//	errMsg := "ParseToken failed ErrTokenUnknown " + e.Error()
		//	log.Error(operationID, errMsg)
		//}

		errMsg := "parse token err " + err.Error() + argMsg
		return false, utils.Wrap(err, errMsg), errMsg
	}
	// if claims.UID != uid {
	// 	errMsg := " uid is not same to token uid " + argMsg + " claims.UID: " + claims.UID
	// 	return false, utils.Wrap(constant.ErrTokenDifferentUserID, errMsg), errMsg
	// }
	// if claims.Platform != constant.PlatformIDToName(utils.StringToInt(platformID)) {
	// 	errMsg := " platform is not same to token platform " + argMsg + " claims platformID: " + claims.Platform
	// 	return false, utils.Wrap(constant.ErrTokenDifferentPlatformID, errMsg), errMsg
	// }
	log.NewDebug(operationID, utils.GetSelfFuncName(), " check ok ", claims.UID, uid)
	return true, nil, ""
}
