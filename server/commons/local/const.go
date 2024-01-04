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

	LockTypeExpire = 1
)
