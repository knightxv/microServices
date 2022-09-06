package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

var Config config

type callBackConfig struct {
	Enable                 bool `yaml:"enable"`
	CallbackTimeOut        int  `yaml:"callbackTimeOut"`
	CallbackFailedContinue bool `yaml:"callbackFailedContinue"`
}

type config struct {
	ServerIP string `yaml:"serverip"`

	RpcRegisterIP string `yaml:"rpcRegisterIP"`
	ListenIP      string `yaml:"listenIP"`

	ServerVersion string `yaml:"serverversion"`
	Api           struct {
		GinPort  []int  `yaml:"openImApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	CmsApi struct {
		GinPort  []int  `yaml:"openImCmsApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	Mysql struct {
		DBAddress      []string `yaml:"dbMysqlAddress"`
		DBUserName     string   `yaml:"dbMysqlUserName"`
		DBPassword     string   `yaml:"dbMysqlPassword"`
		DBDatabaseName string   `yaml:"dbMysqlDatabaseName"`
		DBTableName    string   `yaml:"DBTableName"`
		DBMsgTableNum  int      `yaml:"dbMsgTableNum"`
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
	}
	Mongo struct {
		DBUri               string `yaml:"dbUri"`
		DBAddress           string `yaml:"dbAddress"`
		DBDirect            bool   `yaml:"dbDirect"`
		DBTimeout           int    `yaml:"dbTimeout"`
		DBDatabase          string `yaml:"dbDatabase"`
		DBSource            string `yaml:"dbSource"`
		DBUserName          string `yaml:"dbUserName"`
		DBPassword          string `yaml:"dbPassword"`
		DBMaxPoolSize       int    `yaml:"dbMaxPoolSize"`
		DBRetainChatRecords int    `yaml:"dbRetainChatRecords"`
	}
	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBUserName    string   `yaml:"dbUserName"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	}
	RpcPort struct {
		OpenImUserPort           []int `yaml:"openImUserPort"`
		OpenImFriendPort         []int `yaml:"openImFriendPort"`
		OpenImMessagePort        []int `yaml:"openImMessagePort"`
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"`
		OpenImGroupPort          []int `yaml:"openImGroupPort"`
		OpenImAuthPort           []int `yaml:"openImAuthPort"`
		OpenImPushPort           []int `yaml:"openImPushPort"`
		OpenImStatisticsPort     []int `yaml:"openImStatisticsPort"`
		OpenImMessageCmsPort     []int `yaml:"openImMessageCmsPort"`
		OpenImAdminCmsPort       []int `yaml:"openImAdminCmsPort"`
		OpenImOfficePort         []int `yaml:"openImOfficePort"`
		OpenImOrganizationPort   []int `yaml:"openImOrganizationPort"`
		OpenImConversationPort   []int `yaml:"openImConversationPort"`
		OpenImCachePort          []int `yaml:"openImCachePort"`
	}
	RpcRegisterName struct {
		OpenImStatisticsName string `yaml:"openImStatisticsName"`
		OpenImUserName       string `yaml:"openImUserName"`
		OpenImFriendName     string `yaml:"openImFriendName"`
		//	OpenImOfflineMessageName     string `yaml:"openImOfflineMessageName"`
		OpenImMsgName          string `yaml:"openImMsgName"`
		OpenImPushName         string `yaml:"openImPushName"`
		OpenImRelayName        string `yaml:"openImRelayName"`
		OpenImGroupName        string `yaml:"openImGroupName"`
		OpenImAuthName         string `yaml:"openImAuthName"`
		OpenImMessageCMSName   string `yaml:"openImMessageCMSName"`
		OpenImAdminCMSName     string `yaml:"openImAdminCMSName"`
		OpenImOfficeName       string `yaml:"openImOfficeName"`
		OpenImOrganizationName string `yaml:"openImOrganizationName"`
		OpenImConversationName string `yaml:"openImConversationName"`
		OpenImCacheName        string `yaml:"openImCacheName"`
		OpenImRealTimeCommName string `yaml:"openImRealTimeCommName"`
	}
	Etcd struct {
		EtcdSchema string   `yaml:"etcdSchema"`
		EtcdAddr   []string `yaml:"etcdAddr"`
	}
	Log struct {
		StorageLocation       string   `yaml:"storageLocation"`
		RotationTime          int      `yaml:"rotationTime"`
		RemainRotationCount   uint     `yaml:"remainRotationCount"`
		RemainLogLevel        uint     `yaml:"remainLogLevel"`
		ElasticSearchSwitch   bool     `yaml:"elasticSearchSwitch"`
		ElasticSearchAddr     []string `yaml:"elasticSearchAddr"`
		ElasticSearchUser     string   `yaml:"elasticSearchUser"`
		ElasticSearchPassword string   `yaml:"elasticSearchPassword"`
	}
	ModuleName struct {
		LongConnSvrName string `yaml:"longConnSvrName"`
		MsgTransferName string `yaml:"msgTransferName"`
		PushName        string `yaml:"pushName"`
	}
	LongConnSvr struct {
		WebsocketPort       []int `yaml:"openImWsPort"`
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`
		WebsocketTimeOut    int   `yaml:"websocketTimeOut"`
	}

	Kafka struct {
		Ws2mschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		//Ws2mschatOffline struct {
		//	Addr  []string `yaml:"addr"`
		//	Topic string   `yaml:"topic"`
		//}
		MsgToMongo struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		Ms2pschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		ConsumerGroupID struct {
			MsgToRedis string `yaml:"msgToTransfer"`
			MsgToMongo string `yaml:"msgToMongo"`
			MsgToMySql string `yaml:"msgToMySql"`
			MsgToPush  string `yaml:"msgToPush"`
		}
	}

	TokenPolicy struct {
		AccessSecret string `yaml:"accessSecret"`
		AccessExpire int64  `yaml:"accessExpire"`
	}
}

func init() {
	cfgName := os.Getenv("CONFIG_NAME")
	fmt.Println("get config path is:", Root, cfgName)

	if len(cfgName) != 0 {
		Root = cfgName
	}

	bytes, err := ioutil.ReadFile(filepath.Join(Root, "config", "config.yaml"))
	if err != nil {
		panic(err.Error())
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		panic(err.Error())
	}
}
