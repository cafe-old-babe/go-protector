package service

import (
	"errors"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/database/condition"
	"go-protector/server/internal/ssh/cmd"
	"go-protector/server/internal/ssh/monitor"
	"go-protector/server/internal/ssh/ssh_con"
	"go-protector/server/internal/ssh/ssh_term"
	"go-protector/server/internal/utils"
	"gorm.io/gorm"
	"io"
	"os"
	"path"
)

type SsoSession struct {
	base.Service
}

// CreateSession 创建会话
func (_self *SsoSession) CreateSession(authId uint64) (res *base.Result) {
	if authId <= 0 {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	// 查询授权信息
	var auth entity.AssetAuth
	var err error
	if auth, err = dao.AssetAuth.FindById(_self.DB, authId); err != nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	if auth.StartDate.Valid && auth.EndDate.Valid {
		now := c_type.NowTime()
		if now.Time.Before(auth.StartDate.Time) || now.Time.After(auth.EndDate.Time) {
			res = base.ResultFailureErr(errors.New("授权无效或已过期"))
			return
		}
	}
	// 校验授权信息
	if auth.UserId != current.GetUserId(_self.Context) {
		res = base.ResultFailureErr(errors.New("授权无效或已过期"))
		return
	}

	assetBasic, err := dao.AssetBasic.FindAssetBasicByDTO(_self.DB, dto.FindAssetDTO{
		ID: auth.AssetId,
	})
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	assetAccount, err := dao.AssetAccount.FindAssetAccountByDTO(_self.DB, dto.FindAssetAccountDTO{
		AssetId: auth.AssetId,
		Account: auth.AssetAcc,
	})
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}

	// 创建
	var session entity.SsoSession
	session.AuthId = auth.ID

	session.AssetId = auth.AssetId
	session.AssetName = auth.AssetName
	session.AssetIp = auth.AssetIp
	session.AssetPort = assetBasic.Port
	session.AssetGatewayId = assetBasic.AssetGatewayId

	session.AssetAccId = auth.AssetAccId
	session.AssetAcc = auth.AssetAcc
	session.AssetAccPwd = assetAccount.Password

	session.UserId = auth.UserId
	session.UserAcc = auth.UserAcc

	session.Status = consts.SessionWaiting

	if err = _self.DB.Create(&session).Error; err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	return base.ResultSuccess(map[string]uint64{"id": session.ID}, "创建成功")
}

// ConnectBySession
func (_self *SsoSession) ConnectBySession(req *dto.ConnectBySessionReq) (err error) {
	if req == nil || req.Id <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	var ws *base.WsContext
	// 校验会话信息
	var model entity.SsoSession

	var term *ssh_term.Terminal
	var forward *ssh_term.TermForward

	// websocket
	if ws, err = base.Upgrade(&_self.Service); err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = ws.Write(base.NewWsMsg(consts.MsgClose, err.Error()))
		}
		if model.ID > 0 {
			model.Status = consts.SessionClose
			_self.DB.Updates(&model)
		}
		if term != nil {
			_ = term.Close()
		}

		if forward != nil {
			forward.Stop()
		} else {
			if ws != nil {
				_ = ws.Close()
			}
		}
	}()
	if err = _self.DB.First(&model, req.Id).Error; err != nil {
		return
	}
	assetAccPwd := model.AssetAccPwd
	// test
	if model.Status != consts.SessionWaiting {
		err = c_error.ErrIllegalAccess
		return
	}
	if model.UserId != current.GetUserId(_self.Context) {
		err = c_error.ErrIllegalAccess
		return
	}
	model.Status = consts.SessionConnecting
	if err = _self.DB.Updates(&model).Error; err != nil {
		return err
	}
	// 启动shell
	if term, err = ssh_term.NewTerminal(req, &ssh_con.ConnectParam{
		ID:        model.AssetId,
		IP:        model.AssetIp,
		Port:      model.AssetPort,
		Username:  model.AssetAcc,
		Password:  assetAccPwd,
		GatewayId: model.AssetGatewayId,
	}); err != nil {
		return
	}
	// 更新连接状态
	model.Status = consts.SessionConnected
	model.ConnectAt = c_type.NowTime()
	term.ConnectAt = model.ConnectAt.Time
	// 启动转发
	if forward, err = ssh_term.NewTermForward(ws, term,
		cmd.NewParserHandler(_self.GetContext(), req.Id),
		cmd.NewObserveHandler(_self.GetContext(), req.Id)); err != nil {
		return
	}

	var castPath string
	if castPath, err = forward.GetCastPath(); err != nil {
		return
	}
	model.CastPath = castPath
	if err = _self.DB.Updates(&model).Error; err != nil {
		return err
	}
	forward.Start()

	_ = forward.ReadWsToWriteTerm()

	return
}

func (_self *SsoSession) Page(req *dto.SsoSessionPageReq) (res *base.Result) {

	if req == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var slice []entity.SsoSession
	count, err := _self.Count(
		_self.GetDB().Scopes(
			condition.Paginate(req),
			condition.Like("user_acc ", req.UserAcc),
			condition.Like("asset_ip ", req.AssetIp),
			condition.Like("asset_name ", req.AssetName),
			condition.Like("asset_acc ", req.AssetAcc),
		).Order("created_at desc").Find(&slice))
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	res = base.ResultPage(slice, req, count)
	return

}

func (_self *SsoSession) GetCast(ssoSessionId uint64) (res *base.Result) {
	var err error

	var castFile string
	var ssoSession entity.SsoSession
	if err = _self.DB.First(&ssoSession, ssoSessionId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = c_error.ErrIllegalAccess
		}
		res = base.ResultFailureErr(err)
		return
	}
	if ssoSession.Status != consts.SessionClose {
		res = base.ResultFailureErr(errors.New("会话未结束,无法查看录像"))
		return
	}

	if castFile = ssoSession.CastPath; len(castFile) <= 0 {
		err = errors.New("无法查看回放,请联系管理员")
	}

	if path.Ext(castFile) != ".cast" {
		err = errors.New("回放文件格式不正确,请联系管理员")
	}
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	file, err := os.Open(castFile)
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	defer file.Close()
	all, err := io.ReadAll(file)
	if err != nil {
		res = base.ResultFailureErr(err)
	}
	res = base.ResultSuccess(string(all))

	return

}

func (_self *SsoSession) MonitorBySession(req *dto.ConnectBySessionReq) (err error) {
	if req == nil || req.Id <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	var ws *base.WsContext
	// websocket
	if ws, err = base.Upgrade(&_self.Service); err != nil {
		return
	}
	observer := &monitor.Observer{
		ObId:      utils.GetNextId(),
		SsoId:     req.Id,
		WsContext: ws,
	}
	defer monitor.Subject.RemoveObserver(observer)
	if err = monitor.Subject.RegisterObserver(observer); err != nil {
		return
	}
	if err = ws.Write(base.NewWsMsg(consts.MsgConnected, "")); err != nil {
		return
	}
	for {
		if _, err = observer.ReadMsg(); err != nil {

			break
		}

	}
	return nil
}
