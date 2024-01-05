package local

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

	CtxKeyUserId    = "userId"
	CtxKeyLoginName = "loginName"
	CtxKeyUserName  = "userName"

	LockTypeExpire = 1
)
