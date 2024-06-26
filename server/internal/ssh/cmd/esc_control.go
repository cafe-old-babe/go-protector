package cmd

import (
	"bytes"
	"go-protector/server/internal/custom/c_ascii"
	"strconv"
	"strings"
)

type doFunc func(*EscCtl, *ParserSSHCharHandler, rune)

var doCSIFuncMap map[rune]doFunc

func init() {

	doCSIFuncMap = make(map[rune]doFunc)
	doCSIFuncMap['h'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		// ESC [ ? 1 0 3 4 h
		var lastLine []rune
		switch string(ctl.param) {
		case "?1034":
			if len(handler.ps1Row) > 0 && len(handler.ps1Row[0]) > 0 {
				ctl.reset()
				return
			}

			if len(ctl.temp) > 0 {
				handler.RecordCmd(r)
				// 记录直到空格为止
				if r != ' ' {
					return
				}
				lastLine = handler.cmd[handler.GetY()]
				if len(lastLine) < 1 {
					return
				}
				last := lastLine[handler.GetX()]
				//if last == '#' || last == '$' {
				if last == 0x23 || last == 0x24 {
					ps1 := strings.TrimPrefix(string(lastLine), string(ctl.temp))
					handler.recordPs1(ps1)

					handler.GetLogger().Debug("ps1: %v", ps1)
					ctl.reset()
				}
				return
			}
			lastLine = handler.cmd[handler.GetY()]
			// 获取最后一行,并获取当前光标到左侧的数据
			if ctl.temp = lastLine; len(ctl.temp) <= 0 {
				ctl.temp = append(ctl.temp, 0)
			}
		case "?47", "?1047", "?1048", "?1049": // 开启缓存区
			handler.recordState = false
			fallthrough
		default:
			ctl.reset()
			return
		}
		return
	}
	doCSIFuncMap['l'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		// ESC [ ? 1 0 3 4 l
		var lastLine []rune
		switch string(ctl.param) {
		case "?1034":
			if len(handler.ps1Row) > 0 && len(handler.ps1Row[0]) > 0 {
				ctl.reset()
				return
			}

			if len(ctl.temp) > 0 {
				handler.RecordCmd(r)
				// 记录直到空格为止
				if r != ' ' {
					return
				}
				lastLine = handler.cmd[handler.GetY()]
				if len(lastLine) < 1 {
					return
				}
				last := lastLine[handler.GetX()]
				//if last == '#' || last == '$' {
				if last == 0x23 || last == 0x24 {
					ps1 := strings.TrimPrefix(string(lastLine), string(ctl.temp))
					handler.recordPs1(ps1)

					handler.GetLogger().Debug("ps1: %v", ps1)
					ctl.reset()
				}
				return
			}
			lastLine = handler.cmd[handler.GetY()]
			// 获取最后一行,并获取当前光标到左侧的数据
			if ctl.temp = lastLine; len(ctl.temp) <= 0 {
				ctl.temp = append(ctl.temp, 0)
			}
		case "?47", "?1047", "?1048", "?1049": // 关闭缓存区
			handler.recordState = true
			fallthrough
		default:
			ctl.reset()
			return
		}
		return
	}

	// CSI Ps A  Cursor Up Ps Times (default = 1) (CUU). 上移
	// CSI Ps SP(SPACE) A Shift right Ps columns(s) (default = 1) (SR), ECMA-48. 右移
	doCSIFuncMap['A'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		defer ctl.reset()
		p := 1
		up := true
	paramLabel:
		if len(ctl.param) > 0 {

			param := ctl.param[len(ctl.param)-1]
			if param == 0x20 {
				ctl.param = ctl.param[:len(ctl.param)-1]
				up = false
				goto paramLabel
			}

			if i, err := strconv.ParseInt(string(param), 10, 64); err != nil {
				return
			} else {
				p = int(i)
			}
		}

		if up {
			// 上移
			handler.MoveY(-p)
		} else {
			// 右移
			handler.MoveX(p)

		}

	}
	//  CSI Ps B  Cursor Down Ps Times (default = 1) (CUD). 下移
	doCSIFuncMap['B'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {

		defer ctl.reset()
		ps := 1
		if len(ctl.param) > 0 {
			if i, err := strconv.ParseInt(string(ctl.param), 10, 32); err == nil {
				ps = int(i)
			} else {
				// 转换失败
				return
			}
		}
		handler.MoveY(ps)
	}
	// CSI Ps C  Cursor Forward Ps Times (default = 1) (CUF). 右移 左移为 \b
	doCSIFuncMap['C'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		defer ctl.reset()
		ps := 1
		if len(ctl.param) > 0 {
			if i, err := strconv.ParseInt(string(ctl.param), 10, 32); err == nil {
				ps = int(i)
			} else {
				// 转换失败
				return
			}

		}
		// 如果右移超过实际长度
		if handler.MoveX(ps) > len(handler.cmd[handler.GetY()]) {
			handler.SetX(len(handler.cmd[handler.GetY()]))
		}

	}
	// CSI Ps D  Cursor Backward Ps Times (default = 1) (CUB). 左移
	doCSIFuncMap['D'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		defer ctl.reset()
		ps := 1
		if len(ctl.param) > 0 {
			if i, err := strconv.ParseInt(string(ctl.param), 10, 32); err == nil {
				ps = int(i)
			} else {
				// 转换失败
				return
			}
		}
		// 左移
		handler.MoveX(-ps)
	}
	//CSI Ps E  Cursor Next Line Ps Times (default = 1) (CNL). 下移
	doCSIFuncMap['E'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		defer ctl.reset()
		ps := 1
		if len(ctl.param) > 0 {
			if i, err := strconv.ParseInt(string(ctl.param), 10, 32); err == nil {
				ps = int(i)
			} else {
				// 转换失败
				return
			}
		}
		handler.MoveY(ps)
	}
	//CSI Ps F  Cursor Preceding Line Ps Times (default = 1) (CPL). 上移
	doCSIFuncMap['F'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		defer ctl.reset()
		ps := 1
		if len(ctl.param) > 0 {
			if i, err := strconv.ParseInt(string(ctl.param), 10, 32); err == nil {
				ps = int(i)
			} else {
				// 转换失败
				return
			}
		}
		handler.MoveY(-ps)
	}
	//CSI Ps F  Cursor Preceding Line Ps Times (default = 1) (CPL). 上移
	doCSIFuncMap['G'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		defer ctl.reset()
		ps := 1
		if len(ctl.param) > 0 {
			if i, err := strconv.ParseInt(string(ctl.param), 10, 32); err == nil {
				ps = int(i)
			} else {
				// 转换失败
				return
			}
		}
		handler.MoveY(-ps)
	}
	// 删除上下
	//CSI Ps J  Erase in Display (ED), VT100.
	//            Ps = 0  ⇒  Erase Below (default).   删除光标下方
	//            Ps = 1  ⇒  Erase Above. 删除光标上方
	//            Ps = 2  ⇒  Erase All. 删除所有
	//            Ps = 3  ⇒  Erase Saved Lines, xterm.
	//
	//CSI ? Ps J
	//          Erase in Display (DECSED), VT220.
	//            Ps = 0  ⇒  Selective Erase Below (default).
	//            Ps = 1  ⇒  Selective Erase Above.
	//            Ps = 2  ⇒  Selective Erase All.
	//            Ps = 3  ⇒  Selective Erase Saved Lines, xterm.
	doCSIFuncMap['J'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		defer ctl.reset()
		p := "0" //default
		if len(ctl.param) > 0 {
			if ctl.param[0] == '?' && len(ctl.param) > 1 {
				ctl.param = ctl.param[1:]
			}
			p = string(ctl.param)
		}
		switch p {
		case "0":
			// 删除光标后面及下面的字符
			handler.cmd[handler.GetY()] = handler.cmd[handler.GetY()][:handler.GetX()]
		case "1":
			// 删除光标前面及上面的字符
			handler.cmd[handler.GetY()] = handler.cmd[handler.GetY()][handler.GetX():]
		default:
			handler.ResetCmd()

		}

	}
	//删除 左右
	//CSI Ps K  Erase in Line (EL), VT100.
	//            Ps = 0  ⇒  Erase to Right (default). 删除光标右边
	//            Ps = 1  ⇒  Erase to Left.  删除光标左边
	//            Ps = 2  ⇒  Erase All.  删除整行
	//
	//CSI ? Ps K
	//          Erase in Line (DECSEL), VT220.
	//            Ps = 0  ⇒  Selective Erase to Right (default). 删除光标右边
	//            Ps = 1  ⇒  Selective Erase to Left. 删除光标左边
	//            Ps = 2  ⇒  Selective Erase All. 删除整行
	doCSIFuncMap['K'] = func(ctl *EscCtl, handler *ParserSSHCharHandler, r rune) {
		p := "0" //default
		if len(ctl.param) > 0 {
			if ctl.param[0] == '?' && len(ctl.param) > 1 {
				ctl.param = ctl.param[1:]
			}
			p = string(ctl.param)
		}
		switch p {
		case "0": // 删除右侧
			handler.cmd[handler.GetY()] = handler.cmd[handler.GetY()][:handler.GetX()]
		case "1": // 删除左侧
			handler.cmd[handler.GetY()] = handler.cmd[handler.GetY()][handler.GetX():]
		case "2": // 删除整行
			handler.cmd[handler.GetY()] = handler.cmd[handler.GetY()][:0]

		}
		ctl.reset()
	}
}

// \x1b[6n 服务端发送 要客户端发送光标位置
// \x1b[11;5R
// \x1b[11;5H 以ESC[n;mR(就像在键盘上输入)向应用程序设置光标位置(CPR)，其中n是行，m是列
// ESC ] 0 ; <string><ST> 	将控制台窗口的标题设置为
// ESC ] 2 ; <string><ST> 	将控制台窗口的标题设置为
// ESC [ ? 1 0 4 9 h 	切换到新的备用屏幕缓冲区。 禁止记录
// ESC [ ? 1 0 3 4 h 	?
// ESC [ ? 1 0 4 9 l    切换到主缓冲区。 开启记录
// ESC [ ? 25 l
// ESC [ ? 1 l
// ESC > (62) Normal Keypad (DECKPNM), VT100.
// ESC [ ? 12 l
// ESC [ ? 25 h
// ESC [ ? 1 0 4 9 l

type EscCtl struct {
	isEsc   bool
	ctlCode rune // CSI or OSC Control Characters
	// param Ps A single (usually optional) numeric parameter, composed of one of more digits.
	// param Pm A multiple numeric parameter composed of any number of single numeric parameters, separated by ; character(s). Individual values for the parameters are listed with P s .
	// param Pt A text parameter composed of printable characters.
	param []rune
	// endCode 结束code
	endCode rune
	// temp // 临时存储
	temp []rune
}

func (_self *EscCtl) reset() {
	_self.isEsc = false
	_self.ctlCode = 0
	_self.param = _self.param[:0]
	_self.endCode = 0
	_self.temp = _self.temp[:0]
}

func (_self *EscCtl) doEsc(handler *ParserSSHCharHandler, r rune) {
	if r == 0x33 {
		_self.isEsc = true
		return
	}

	if _self.ctlCode <= 0 {
		// 补全 ctlCode
		switch r {
		// Controls beginning with ESC 没有 endCode
		case '6', '7', '8', '9', '=', '>', 'F', 'c', 'l', 'm', 'n', 'o', '|', '}', '~':
			_self.reset()
		default:
			_self.ctlCode = r
		}

		return
	} else {
		switch _self.ctlCode {
		case ' ', '#', '%', '(', ')', '*', '+', '-', '.', '\\':
			// Controls beginning with ESC 有 endCode 不处理
			_self.reset()
			return
		case ']', // Operating System Command(OSC): ]
			'_', //  Application Program Command(APC):_
			'p', // Device Control String(DCS): p
			'^': // Privacy Message(PM): ^

			if r == c_ascii.ST || r == c_ascii.BEL {
				_self.reset()
			}
			return
		}

	}

	if bytes.IndexRune(parametersStr, r) >= 0 {
		_self.param = append(_self.param, r)
		return
	} else if _self.endCode == 0 && bytes.IndexRune(alphabeticStr, r) >= 0 {
		_self.endCode = r
	}
	// C1 [ CSI
	switch _self.ctlCode {
	case 0x5b: // C1 [ CSI
		if f, ok := doCSIFuncMap[_self.endCode]; ok {
			f(_self, handler, r)
		} else {
			_self.reset()
		}
	default:
		_self.reset()
		return
	}

	return

	//if bytes.IndexRune(parametersStr, r) >= 0 {
	//	_self.param = append(_self.param, r)
	//
	//} else if bytes.IndexRune(alphabeticStr, r) >= 0 && _self.endCode == 0 {
	//	_self.endCode = r
	//
	//} else {
	//	_self.processEndCode(handler, r)
	//}

}
