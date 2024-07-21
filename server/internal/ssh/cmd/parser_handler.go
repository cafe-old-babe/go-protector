package cmd

import (
	"context"
	"fmt"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service/iface"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_ascii"
	"go-protector/server/internal/custom/c_structure"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/ssh/notify"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	alphabeticStr = []byte(string(c_ascii.Alphabetic))
	parametersStr = []byte(string(c_ascii.Parameters) + string(c_ascii.Intermediate))
)

type ParserSSHCharHandler struct {
	id              uint64                      // id
	lastPS1         string                      // lastPS1
	ps1Row          []string                    // ps1Row 每行的PS1
	recordPS1       bool                        // recordPS1 是否记录PS1
	quoteStack      c_structure.SafeStack[rune] // quoteStack 引号栈
	recordState     bool                        // recordState true 记录, false 不记录
	currentFirstIn  bool                        // currentFirstIn 第一次输入命令
	cmd             [][]rune                    // cmd 记录的命令行
	escCtl          EscCtl                      // escCtl 控制字符
	operationEntity *entity.SsoOperation        // operationEntity 操作记录
	cmdSort         int                         // cmdSort 操作记录顺序
	sync.Mutex                                  // mutex 互斥锁
	cursor                                      // cursor 光标所在位置
	base.Service
	// 审批
	approveRecordService iface.IApproveRecordService
	approveCmdSlice      []string                     // approveCmdSlice 需要审批的命令
	assetBasic           entity.AssetBasic            // assetBasic 资产信息
	currentUser          *current.User                // currentUser 当前用户
	approveRecord        *entity.ApproveRecord        // approveRecord 审批记录
	writeCliDataChan     chan<- rune                  // writeCliDataChan 写向客户端
	writeServer          func(in []byte) (int, error) // writeServer 写向服务端
}

func NewParserHandler(c context.Context, id uint64, dataChan chan<- rune,
	writeServer func(in []byte) (int, error)) (_self *ParserSSHCharHandler) {

	_self = new(ParserSSHCharHandler)
	_self.id = id
	_self.Make(c)
	_self.ResetCmd()
	// 审批
	_self.writeCliDataChan = dataChan
	_self.writeServer = writeServer
	_self.approveRecordService = iface.ApproveRecordService(c)
	assetId := c.Value("assetId").(uint64)
	_self.approveCmdSlice, _ = iface.ApproveCmdService(c).GetApproveCmdSliceByAssetId(assetId)
	_self.currentUser, _ = current.GetUser(c)
	if len(_self.approveCmdSlice) > 0 {
		_self.assetBasic, _ = dao.AssetBasic.FindAssetBasicByDTO(_self.GetDB(), dto.FindAssetDTO{
			ID: assetId,
		})

	}

	return
}

func (_self *ParserSSHCharHandler) GetIndex() int {
	return DefaultCmdHandler.GetIndex() + 1
}

func (_self *ParserSSHCharHandler) GetId() uint64 {
	return _self.id
}

func (_self *ParserSSHCharHandler) PassToClient(r rune) {
	if _self.recordState {
		//_self.GetLogger().Debug("passToClient: %s\t%v", strconv.QuoteRune(r), r)
	}
	_self.Lock()
	defer _self.Unlock()
	_self.RecordServerWrite(r)
	return
}

func (_self *ParserSSHCharHandler) PassToServer(r rune) bool {
	if !_self.recordState {
		if r == 0x03 {
			_self.ResetCmd()
		}
		return true
	}
	//_self.GetLogger().Debug("passToServer: %s\t%v", strconv.QuoteRune(r), r)
	_self.Lock()
	defer _self.Unlock()
	if (len(_self.ps1Row) <= 0 ||
		len(_self.ps1Row[0]) <= 0 ||
		len(_self.ps1Row[_self.GetY()]) <= 0) && // 如果没有记录PS1
		// 如果开启记录PS1
		_self.recordPS1 {
		_self.recordPs1(_self.getCmdByLine(_self.GetY()))
	}

	if _self.approveRecord != nil && r != 0x03 {

		sprintf := fmt.Sprintf("\r\n请等待 %s(%s) 审批, ctrl+c 可取消审批",
			_self.approveRecord.ApproveUser.Username,
			_self.approveRecord.ApproveUser.LoginName)
		_self.WriteCli(sprintf)
		return false
	}
	switch r {
	case 0x03: // ctrl+c
		if _self.approveRecord != nil { // 取消

			// 取消审批
			_ = _self.approveRecordService.DoApprove(&dto.DoApproveDTO{
				Id:            _self.approveRecord.ID,
				ApproveStatus: consts.ApproveStatusCancel,
				ApproveUserId: _self.currentUser.ID,
			})
			notify.ApproveManager.UnSubscribe(_self.approveRecord.ID)
			_self.ResetCmd()
			_, _ = _self.writeServer([]byte{0x03})

		} else {
			_self.ResetCmd()
		}

	case 0x0d: // \r
		_self.recordPS1 = true
		if _self.recordState {

			_self.printCmdInfo()
		}

		// \ ' " '
		lastRow := len(_self.cmd) - 1
		last := _self.cmd[lastRow][len(_self.cmd[lastRow])-1]
		if last == '\'' {
			// 跳过记录与清除
			break
		}
		cmd := _self.getCmd()
		defer func() {
			// top 特殊处理
			if strings.HasPrefix(cmd, "top") {
				_self.recordState = false
			}
		}()
		i := strings.IndexAny(cmd, "`'\"")

		if i > -1 {
			// 判断是否成对出现
			_self.quoteStack.Clear()
			_self.forEachCmd(i, 0, func(r rune) {
				if !(r == '\'' || r == '"' || r == '`') {
					return
				}
				if data, exists := _self.quoteStack.Top(); exists && data == r {
					_self.quoteStack.Pop()
				} else {
					_self.quoteStack.Push(r)
				}
			}, func() {

			})

		}
		if !_self.quoteStack.IsEmpty() {
			break
		}
		// _self.doSave(match)
		// _self.ResetCmd()
		// 匹配
		match := _self.matchApproveCmd()
		res := _self.doSave(match)
		if !res.IsSuccess() {

			sprintf := fmt.Sprintf("记录命令失败: %v", res.Message)
			_self.GetLogger().Error(sprintf)
			_self.WriteCli("\r\n" + sprintf + ",请稍后重试")
			_self.ResetCmd()
			_, _ = _self.writeServer([]byte{0x03})
			return false

		}

		if !match {
			_self.ResetCmd()

			return true
		}
		_self.operationEntity = res.Data.(*entity.SsoOperation)
		// 创建审批
		applicantContent := fmt.Sprintf("用户: %s(%s) [%s]在资产: [%s(%s)] 执行了: %s, 触发了审批,请您处理",
			_self.currentUser.UserName, _self.currentUser.LoginName, time.Now().Format(time.DateTime),
			_self.assetBasic.AssetName, _self.assetBasic.IP, cmd)

		res = _self.approveRecordService.Insert(&dto.ApproveRecordInsertDTO{
			ApplicantId:      _self.currentUser.ID,
			ApproveUserId:    _self.assetBasic.ManagerUserId,
			SessionId:        _self.currentUser.SessionId,
			ApplicantContent: applicantContent,
			Timeout:          5 * 60,
			ApproveType:      consts.ApproveTypSsoOperation,
			ApproveBindId:    _self.operationEntity.ID,
		})

		if !res.IsSuccess() {
			sprintf := fmt.Sprintf("创建审批失败: %v", res.Message)
			_self.GetLogger().Error(sprintf)
			_self.WriteCli("\r\n" + sprintf)
			_self.ResetCmd()
			_, _ = _self.writeServer([]byte{0x03})
		} else {
			_self.approveRecord = res.Data.(*entity.ApproveRecord)
			sprintf := fmt.Sprintf("\r\n请等待 %s(%s) 审批, ctrl+c 可取消审批",
				_self.approveRecord.ApproveUser.Username,
				_self.approveRecord.ApproveUser.LoginName)

			notify.ApproveManager.Subscribe(_self.approveRecord.ID, func(record entity.ApproveRecord) {

				_self.WriteCli(fmt.Sprintf("\r\n审批结果: %s, 审批消息: %s",
					record.ApproveStatusText, record.ApproveContent))
				var dataByte byte
				switch record.ApproveStatus {
				case consts.ApproveStatusPass:
					// 通过
					_self.doSave(true)
					_self.ResetCmd()
					dataByte = 0x0d
				default:
					// 非通过
					_self.GetDB().Model(_self.operationEntity).Updates(&entity.SsoOperation{
						ExecStatus: "2",
						CmdExecAt:  c_type.Time{},
					})
					dataByte = 0x03

				}
				_self.ResetCmd()
				_, _ = _self.writeServer([]byte{dataByte})

			})

			_self.WriteCli(sprintf)

		}
		return false
	}

	return true
}

func (_self *ParserSSHCharHandler) Close() {
	_self.ResetCmd()

}

func (_self *ParserSSHCharHandler) RecordServerWrite(r rune) {
	if _self.escCtl.isEsc {
		_self.escCtl.doEsc(_self, r)
		return
	}

	switch r {
	case c_ascii.ESCKey:
		_self.escCtl.isEsc = true
	case 0x08: // \b 退格

		_self.MoveX(-1)
	default:
		if _self.recordState {
			_self.RecordCmd(r)
			_self.printCmdInfo()
		}
	}

}

func (_self *ParserSSHCharHandler) ResetCmd() {
	_self.ps1Row = make([]string, 1)
	_self.cmd = make([][]rune, 1)
	_self.recordPS1 = true
	_self.recordState = true
	_self.currentFirstIn = true
	_self.ResetCursor()
	_self.quoteStack.Clear()
	_self.operationEntity = new(entity.SsoOperation)
	_self.approveRecord = nil

}

func (_self *ParserSSHCharHandler) printCmdInfo() {
	if len(_self.ps1Row) <= 0 {
		return
	}
	//_self.GetLogger().Debug("-------------current cmd------------------")
	//
	//_self.GetLogger().Debug("firstPs1: %s,x: %d,y %d\n%s\n↓↓↓↓↓↓↓↓onlyCmd↓↓↓↓↓↓↓↓\n%s\n↑↑↑↑↑↑↑onlyCmd↑↑↑↑↑↑↑",
	//	_self.ps1Row[0], _self.GetX(), _self.GetY(), _self.getPs1AndCmd(), _self.getCmd())
	//
	//_self.GetLogger().Debug("-------------current cmd------------------")
}

func (_self *ParserSSHCharHandler) getPs1AndCmd() string {

	var sb strings.Builder
	for i := range _self.cmd {
		line := string(_self.cmd[i])
		sb.WriteString(line)

		if i < len(_self.cmd)-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func (_self *ParserSSHCharHandler) getCmd() string {

	var sb strings.Builder
	for i := range _self.cmd {
		line := string(_self.cmd[i])
		if len(_self.ps1Row) > i {
			line = strings.TrimPrefix(line, _self.ps1Row[i])
		}

		sb.WriteString(line)

		if i < len(_self.cmd)-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func (_self *ParserSSHCharHandler) getCmdByLine(i int) string {

	if i < 0 {
		abs := int(math.Abs(float64(i)))
		if abs > len(_self.cmd) {
			return ""
		}
		return string(_self.cmd[len(_self.cmd)+i])
	}
	if i > len(_self.cmd)-1 {
		return ""
	}
	return string(_self.cmd[i])
}

func (_self *ParserSSHCharHandler) RecordCmd(r rune) {
	switch r {
	case 0x0d: //\r
		_self.ResetX()
	case 0x0a: // \n
		y := _self.MoveY(1)
		newLine := make([]rune, 0)
		if y > len(_self.cmd)-1 {
			_self.cmd = append(_self.cmd, newLine)
			_self.ps1Row = append(_self.ps1Row, "")
		} else {
			newCmd := make([][]rune, len(_self.cmd)+1)
			copy(newCmd, _self.cmd[:y])
			newCmd[y] = newLine
			copy(newCmd[y+1:], _self.cmd[y:])
			_self.cmd = newCmd
			_self.ps1Row = append(_self.ps1Row[:y], append(make([]string, 1), _self.ps1Row[y:]...)...)

		}
		_self.ResetX()
	case 0x07: // \a
		return
	default:
		y := _self.GetY()
		x := _self.GetX()
		if x < len(_self.cmd[y]) {
			_self.cmd[y][x] = r
		} else {
			_self.cmd[y] = append(_self.cmd[y][:x],
				append([]rune{r}, _self.cmd[y][x:]...)...)
		}
		_self.MoveX(1)
	}

}

func (_self *ParserSSHCharHandler) recordPs1(ps1 string) {
	if _self.currentFirstIn {
		// 清空上面的数据
		_self.ps1Row = _self.ps1Row[_self.GetY():]
		_self.cmd = _self.cmd[_self.GetY():]
		_self.ResetY()
		_self.currentFirstIn = false
		_self.operationEntity.PS1 = ps1
		_self.operationEntity.CmdStartAt = c_type.NowTime()
	}
	//if len(_self.ps1Row) <= 0 {
	//	_self.cmd = _self.cmd[len(_self.cmd)-1:]
	//	_self.ResetY()
	//}

	_self.ps1Row[_self.GetY()] = ps1
	//_self.ps1Row = append(_self.ps1Row, ps1)
	_self.recordPS1 = false
}

func (_self *ParserSSHCharHandler) forEachCmd(x, y int, runeFunc func(r rune), newLine func()) {
	if y < 0 || y > len(_self.cmd) {
		return
	}
	rowIndex := y
	colIndex := x
	var rows []rune
	for rowIndex < len(_self.cmd) {
		rows = []rune(strings.TrimPrefix(string(_self.cmd[rowIndex]), _self.ps1Row[rowIndex]))
		for ; colIndex < len(rows); colIndex++ {
			runeFunc(rows[colIndex])
		}
		if rowIndex != len(_self.cmd)-1 {
			newLine()
		}

		colIndex = 0
		rowIndex++
	}

}

// doSave 保存
func (_self *ParserSSHCharHandler) doSave(match bool) (res *base.Result) {
	cmd := _self.getCmd()
	if len(cmd) <= 0 {
		return
	}
	_self.operationEntity.Cmd = cmd
	_self.operationEntity.CmdExecAt = c_type.NowTime()
	_self.operationEntity.Sort = _self.cmdSort

	_self.operationEntity.SsoSessionId = _self.GetId()
	if match {
		// 待审批
		_self.operationEntity.ExecStatus = "1"
	} else {
		// 执行成功
		_self.operationEntity.ExecStatus = "0"
		_self.cmdSort++
	}
	return _self.SimpleSave(_self.operationEntity)

}

// matchApproveCmd 匹配命令
func (_self *ParserSSHCharHandler) matchApproveCmd() (match bool) {
	cmd := _self.getCmd()
	if len(cmd) <= 0 {
		return false
	}
	for _, approveCmd := range _self.approveCmdSlice {
		if match, _ = regexp.MatchString(approveCmd, cmd); !match {
			continue
		} else {
			break
		}
	}
	return

}

// WriteCli 写向客户端
func (_self *ParserSSHCharHandler) WriteCli(str string) {
	if len(str) <= 0 {
		return
	}
	go func() {

		for _, strRune := range []rune(str) {
			_self.writeCliDataChan <- strRune
		}
	}()

}
