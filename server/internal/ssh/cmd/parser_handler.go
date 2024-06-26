package cmd

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_ascii"
	"go-protector/server/internal/custom/c_structure"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/utils/async"
	"math"
	"strconv"
	"strings"
	"sync"
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
	mutex           sync.Mutex                  // mutex 互斥锁
	escCtl          EscCtl                      // escCtl 控制字符
	operationEntity *entity.SsoOperation        // operationEntity 操作记录
	cmdSort         int                         // cmdSort 操作记录顺序
	cursor                                      // cursor 光标所在位置
	base.Service
}

func NewParserHandler(c *gin.Context, id uint64) (_self *ParserSSHCharHandler) {
	_self = new(ParserSSHCharHandler)
	_self.id = id
	_self.Make(c)
	_self.ResetCmd()
	return
}

func (_self *ParserSSHCharHandler) GetIndex() int {
	return DefaultCmdHandler.GetIndex() + 1
}

func (_self *ParserSSHCharHandler) GetId() uint64 {
	return _self.id
}

func (_self *ParserSSHCharHandler) PassToClient(r rune) bool {
	_self.GetLogger().Debug("passToClient: %s\t%v", strconv.QuoteRune(r), r)
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	_self.RecordServerWrite(r)
	return true
}

func (_self *ParserSSHCharHandler) PassToServer(r rune) bool {
	_self.GetLogger().Debug("passToServer: %s\t%v", strconv.QuoteRune(r), r)
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	if (len(_self.ps1Row) <= 0 ||
		len(_self.ps1Row[0]) <= 0 ||
		len(_self.ps1Row[_self.GetY()]) <= 0) && // 如果没有记录PS1
		// 如果开启记录PS1
		_self.recordPS1 {
		_self.recordPs1(_self.getCmdByLine(_self.GetY()))

	}
	switch r {
	case 0x03: // ctrl+c
		_self.ResetCmd()
	case 0x0d: // \r
		_self.recordPS1 = true
		_self.printCmdInfo()

		// \ ' " '
		lastRow := len(_self.cmd) - 1
		last := _self.cmd[lastRow][len(_self.cmd[lastRow])-1]
		if last == '\'' {
			// 跳过记录与清除
			break
		}
		cmd := _self.getCmd()

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
		_self.doSave()
		_self.ResetCmd()

	}

	return true
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
}

func (_self *ParserSSHCharHandler) printCmdInfo() {
	if len(_self.ps1Row) <= 0 {
		return
	}
	_self.GetLogger().Debug("-------------current cmd------------------")

	_self.GetLogger().Debug("firstPs1: %s,x: %d,y %d\n%s\n↓↓↓↓↓↓↓↓onlyCmd↓↓↓↓↓↓↓↓\n%s\n↑↑↑↑↑↑↑onlyCmd↑↑↑↑↑↑↑",
		_self.ps1Row[0], _self.GetX(), _self.GetY(), _self.getPs1AndCmd(), _self.getCmd())

	_self.GetLogger().Debug("-------------current cmd------------------")
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
func (_self *ParserSSHCharHandler) doSave() {
	_self.operationEntity.Cmd = _self.getCmd()
	_self.operationEntity.CmdExecAt = c_type.NowTime()
	_self.operationEntity.Sort = _self.cmdSort
	_self.cmdSort++
	_self.operationEntity.SsoSessionId = _self.GetId()
	ssoOperation := _self.operationEntity
	async.CommonWorkPool.Submit(func() {
		_self.SimpleSave(ssoOperation)
	})
}
