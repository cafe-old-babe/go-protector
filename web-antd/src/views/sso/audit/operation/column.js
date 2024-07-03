export const Columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '资产名称',
    customRender: (text, record) => record.ssoSession.assetName
  },
  {
    title: '资产IP',
    customRender: (text, record) => record.ssoSession.assetIp
  },
  {
    title: '资产端口',
    customRender: (text, record) => record.ssoSession.assetPort
  },
  {
    title: '从帐号',
    customRender: (text, record) => record.ssoSession.assetAcc
  },
  {
    title: '主帐号',
    customRender: (text, record) => record.ssoSession.userAcc
  },
  {
    title: '连接状态',
    customRender: (text, record) => record.ssoSession.statusText

  },
  {
    title: 'PS1',
    dataIndex: 'PS1'
  },

  {
    title: '执行的命令',
    dataIndex: 'cmd'
  }
]
