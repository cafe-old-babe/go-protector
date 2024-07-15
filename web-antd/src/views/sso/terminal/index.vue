<script>
import 'xterm/css/xterm.css'
// npm install xterm
// yarn add xterm
import { Terminal } from 'xterm'

import { FitAddon } from 'xterm-addon-fit'
import { setDocumentTitle } from '@/utils/domUtil'
import WsMsg from './lib/WsMsg'
import request from '@/utils/request'
export default {
  name: 'Terminal',
  data() {
    return {
      xterm: null,
      connected: false,
      ssoTerminal: {},
       inputMsg: ''
    }
  },
  // https://github.com/xtermjs/xterm.js/blob/3.14.2/README.md#addons
  // https://xtermjs.org/
  mounted() {
    const ssoTerminalStr = localStorage.getItem('ssoTerminal')
    console.log(ssoTerminalStr)
   this.ssoTerminal = JSON.parse(ssoTerminalStr)
    // localStorage.removeItem('ssoTerminal')
    // setDocumentTitle('单点登录')
    setDocumentTitle(this.ssoTerminal.title)
    this.initTerm()
    this.initWebsocket()
  },
  methods: {
    initTerm: function () {
      this.xterm = new Terminal({
        theme: {
          background: '#1b1b1b'
        }
      })
      this.xterm.open(document.getElementById('terminal'))
      const fitAddon = new FitAddon()
      this.xterm.loadAddon(fitAddon)
      fitAddon.fit()
      this.xterm.focus()
      // xterm.writeln('trying to connect to the server ...')
      if (this.ssoTerminal.initMsg) {
        // this.xterm.writeln('\x1B[1;3;31m正在连接,请稍后\x1B[0m $ ')
        this.xterm.writeln(this.ssoTerminal.initMsg)
      }

      console.log(this.xterm.cols, this.xterm.rows)
    },
    initWebsocket: function () {
      // 获取当前路由信息
      // const id = this.$route.query?.id
      const id = this.ssoTerminal.id
      if (!id) {
        this.xterm.clear()
        this.xterm.writeln('\x1B[1;3;31m非法访问\x1B[0m ')
        return
      }
      const paramStr = 'h=' + this.xterm.rows + '&w=' + this.xterm.cols + '&token=' + this.$store.getters.token
      const origin = window.location.origin
      const wsOrigin = origin.replace('http', 'ws')

      const wsUrl = `${wsOrigin}${this.ssoTerminal.uri}${id}?${paramStr}`
      console.log(wsUrl)
      const ws = new WebSocket(wsUrl)
      this.xterm.onData(data => {
        if (!this.connected || this.websocket === null) {
          return
        }
        if (this.ssoTerminal?.send) { ws.send(new WsMsg(WsMsg.MsgData, data).toString()) }
      })
      ws.onclose = (e) => {
        this.connected = false
        this.xterm.writeln('已断开连接')
      }
      ws.onmessage = (e) => {
        const msg = WsMsg.parse(e.data)
        switch (msg.msgNum) {
          case WsMsg.MsgConnected:
            this.connected = true
            this.xterm.clear()
            break
          case WsMsg.MsgData:
            this.xterm.write(msg.body)
            break
          case WsMsg.MsgAlarm:
            this.$notification.warn({
              message: '监控告警消息',
              description: msg.body,
              duration: 0
            })
            break
          case WsMsg.MsgClose:

            this.xterm.writeln(msg.body)
            this.xterm.writeln('连接已关闭')
            ws.close()
            break
        }
      }
    },
    alarm: function () {
      if (!this.inputMsg) {
        return
      }
      request.post(`/sso-session/operationForMonitor`, {
        'ssoSessionId': this.ssoTerminal.id,
        'type': '0',
        'message': this.inputMsg
      }).then(res => {
        const { code, message } = res
        if (code !== 200) {
          this.$message.error(message)
          return
        }
        this.inputMsg = ''
        this.$message.success('告警成功')
      })
    },
    off: function () {
      if (!this.inputMsg) {
        return
      }
      request.post(`/sso-session/operationForMonitor`, {
        'ssoSessionId': this.ssoTerminal.id,
        'type': '1',
        'message': this.inputMsg
      }).then(res => {
        const { code, message } = res
        if (code !== 200) {
          this.$message.error(message)
          return
        }
        this.inputMsg = ''
        this.$message.success('阻断成功')
      })
    }
  }
}
</script>

<template>
  <div
    :style="{
      position: 'absolute',
      top: 0,
      left: 0,
      width: '100%',
      height: '100%',
      overflow: 'hidden',
      backgroundColor: '#1b1b1b'
    }"
  >
    <div
      id="terminal"
      :style="{
        position: 'absolute',
        top: 0,
        left: 0,
        width: 'calc(100% + 15px)',
        height: 'calc(100% - 70px)',
        overflow: 'hidden',
        backgroundColor: '#1b1b1b'
      }"
    >

      <div
        :style="{
          position: 'fixed',
          bottom: 0,
          left: '10%',
          width: '100%',
          height: '50px',
          backgroundColor: '#1b1b1b'
        }"
      >
        <a-row v-if="!this.ssoTerminal.send" justify="center" gutter="16">
          <a-col :span="16">
            <a-input-search
              ref="alarm"
              placeholder="请输入告警内容"
              v-model="inputMsg"
              enter-button="发送"
              size="large"
              @search="alarm"
            />
          </a-col>
          <a-col :span="4">
            <a-button type="danger" size="large" icon="poweroff" @click="off">
              实时阻断
            </a-button>
          </a-col>
        </a-row>
      </div>
    </div>

  </div>
</template>

<style scoped lang="less">

</style>
