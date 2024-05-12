package consts

import "go-protector/server/internal/custom/c_type"

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

	LockTypeExpire          = 1
	LockTypePasswordFailure = 2

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
	// LoginPolicyShare 共享登录校验
	LoginPolicyShare c_type.LoginPolicyCode = "share"
	// LoginPolicyIntruder 防爆破登录策略
	LoginPolicyIntruder c_type.LoginPolicyCode = "intruder"
)

const (
	// OnlineUserCacheKeyFmt 存放在线用户 -> online:{登录名}:{sessionId}
	OnlineUserCacheKeyFmt = CachePrefix + ":online:%s:%s"
	// OnlineUserCacheLastKeyFmt 换token之前保留当前token 最后一个有效token
	OnlineUserCacheLastKeyFmt = OnlineUserCacheKeyFmt + ":last"
	// LoginPolicyCacheKeyFmt 存放策略 -> login:{登录名}:{sessionId}:policy
	LoginPolicyCacheKeyFmt = CachePrefix + ":login:%s:%s:policy"
	// LoginIntruderCacheKeyFmt 防爆破策略 -> login:intruder:{day}:{登录名}
	LoginIntruderCacheKeyFmt = CachePrefix + ":login:intruder:%d:%s"
)

const (
	// CollFmt 采集脚本
	CollFmt = "cat /etc/passwd | grep ^%s  &&  cat /etc/shadow | grep ^%s"
)
