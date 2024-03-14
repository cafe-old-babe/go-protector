package consts

import "go-protector/server/core/custom/c_type"

const (
	// EnvConfig 环境变量
	EnvConfig    = "config"
	EnvMigration = "migration"

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
	ServerUrlPrefix = "/api"

	AuthHeaderKey = "Authorization"
	AuthUrlKey    = "token"
)

var EmptyVal any

const (
	User c_type.RelationType = "user"
	Menu c_type.RelationType = "menu"
	Dept c_type.RelationType = "dept"
)

const (
	// LoginPolicyGlobal 全局配置
	LoginPolicyGlobal c_type.LoginPolicyCode = "global"
	// LoginPolicyEmail 邮件配置
	LoginPolicyEmail c_type.LoginPolicyCode = "email"
	// LoginPolicyOtp OTP认证
	LoginPolicyOtp c_type.LoginPolicyCode = "otp"
)
