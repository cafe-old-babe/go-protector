export const Columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '资产名称',
    customRender: (text, record) => record.assetBasic.assetName
  },
  {
    title: '资产IP',
    customRender: (text, record) => record.assetBasic.ip
  },
  {
    title: '从帐号',
    dataIndex: 'account'
  },
  {
    title: '从帐号类型',
    dataIndex: 'accountTypeText'
  },
  {
    title: '从帐号状态',
    dataIndex: 'accountStatusText'
  },
  {
    title: '拨测状态',
    scopedSlots: { customRender: 'dailStatus' }
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }
]
