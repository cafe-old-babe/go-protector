package consts

const (
	// EnvConfig 环境变量
	EnvConfig = "config"

	CfgEnvTest    = "test"
	CfgEnvDebug   = "debug"
	CfgEnvRelease = "release"

	// CtxKeyLog ctx 中logger的key
	CtxKeyLog = "local-logger"

	CtxKeyTraceId = "traceId"
	CtxKeyDB      = "db"

	CtxKeyCurrentUser = "currentUser"
	CtxKeyUserId      = "userId"
	CtxKeyLoginName   = "loginName"
	CtxKeyUserName    = "userName"

	LockTypeExpire = 1

	// CachePrefix 缓存前缀
	CachePrefix = "go-protector"

	// ServerUrlPrefix 服务前缀
	ServerUrlPrefix = "api"

	AuthHeaderKey = "Authorization"
	AuthUrlKey    = "token"
)

var EmptyVal any
