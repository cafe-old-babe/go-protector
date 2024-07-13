const WsMsg = class WsMsg {
  constructor (msgNum, body) {
    this.msgNum = msgNum
    this.body = body
  }

  toString () {
    return this.msgNum + ' ' + this.body
  }

  static MsgConnected = 0;
  static MsgClose = 1;
  static MsgData = 2;
  static MsgAlarm = 3;

  static parse (s) {
    const indexOf = s.indexOf(' ')

    const msgNum = parseInt(s.substring(0, indexOf))
    const body = s.substring(indexOf + 1, s.length)
    return new WsMsg(msgNum, body)
  }
}
export default WsMsg
