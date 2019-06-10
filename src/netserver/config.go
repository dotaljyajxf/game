package netserver

type Config struct {
	PublicAddr string

	MaxClientReq   int
	MaxClientResp  int
	UserIdleTimeMs int
	FrontPingMs    int
	ServerId       string //游戏服
	ConfigFile     string

	OnlineMaxNum int

	MaxPackageLen    int
	CommonPackageLen int
	MaxConnNum       int
	Weight           int

	LogPath  string
	LogName  string
	LogLevel int

	LanIp   string
	LanPort int
	WanIp   string
	WanPort int

	EtcdEndPoints []string
}

var GlobalConfig Config

//
//func getValue(needpanic bool, node ini.Node, key string, value interface{}) bool {
//	n := node.Child(key)
//	if n == nil {
//		if !needpanic {
//			return false
//		} else {
//			panic("fail to get key " + key)
//		}
//	}
//	err := n.Value(value)
//	if err != nil {
//		if !needpanic {
//			return false
//		} else {
//			panic("fail to get value " + key)
//		}
//	}
//	return true
//}
//
//func ParseConfig(content string, needpanic bool) (ret bool) {
//	ret = false
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Printf("parse config failed. recover:%v, stack:\n%s\n", r, debug.Stack())
//		}
//	}()
//
//	node, err := ini.ParseFromString(content)
//	if err != nil {
//		if needpanic {
//			fmt.Printf("parse config err:%s\n", err.Error())
//		}
//		ret = false
//		return
//	}
//
//	getValue(needpanic, node, "server.max_client_req", &GlobalConfig.MaxClientReq)
//	getValue(needpanic, node, "server.max_client_resp", &GlobalConfig.MaxClientResp)
//	getValue(needpanic, node, "server.user_idle_time_ms", &GlobalConfig.UserIdleTimeMs)
//	getValue(needpanic, node, "server.net_codec_type", &GlobalConfig.FrontCodecType)
//	getValue(needpanic, node, "server.client_ping_ms", &GlobalConfig.FrontPingMs)
//	getValue(needpanic, node, "server.user_logoff_method", &GlobalConfig.UserLogoffMethodName)
//	getValue(needpanic, node, "server.login_max_num", &GlobalConfig.LoginMaxNum)
//	getValue(needpanic, node, "server.online_max_num", &GlobalConfig.OnlineMaxNum)
//	getValue(needpanic, node, "server.max_package_len", &GlobalConfig.MaxPackageLen)
//	getValue(needpanic, node, "server.common_package_len", &GlobalConfig.CommonPackageLen)
//	getValue(needpanic, node, "server.private_idle_time_ms", &GlobalConfig.PrivateIdleTimeMs)
//	getValue(needpanic, node, "server.max_client_broadcast", &GlobalConfig.MaxClientBroadcast)
//	getValue(needpanic, node, "server.max_web_package_len", &GlobalConfig.MaxWebPackageLen)
//	getValue(needpanic, node, "server.max_conn_num", &GlobalConfig.MaxConnNum)
//	getValue(needpanic, node, "server.weight", &GlobalConfig.Weight)
//
//	if getValue(needpanic, node, "server.broadcast_package_len", &GlobalConfig.BroadcastPackageLen) {
//		GlobalConfig.BroadcastPackageLen = GlobalConfig.BroadcastPackageLen * 1024
//	}
//
//	var messCode string
//	if getValue(needpanic, node, "server.mess_code", &messCode) {
//		GlobalConfig.ArrMessCode = [][]byte{[]byte(messCode)}
//	}
//	ret = true
//	return
//}
