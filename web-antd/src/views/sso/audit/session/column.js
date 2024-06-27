export const Columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '资产名称',
    customRender: (text, record) => record.assetName
  },
  {
    title: '资产IP',
    customRender: (text, record) => record.assetIp
  },
  {
    title: '资产端口',
    customRender: (text, record) => record.assetPort
  },
  {
    title: '登录从帐号',
    dataIndex: 'assetAcc'
  },
  {
    title: '登录主帐号',
    dataIndex: 'userAcc'
  },
  {
    title: '连接状态',
    dataIndex: 'statusText'
  },
  {
    title: '连接时间',
    dataIndex: 'connectAt'
  },
  {
    title: '会话结束时间',
    customRender: (text, record) => record.updatedAt.Valid ? new Date(record.updatedAt.Time).toLocaleString() : ''
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }
]
