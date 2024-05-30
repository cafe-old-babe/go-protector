export const Columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '主帐号',
    dataIndex: 'userAcc'
  },
  {
    title: '资产名称',
    dataIndex: 'assetName'
  },
  {
    title: '资产IP',
    dataIndex: 'assetIp'
  },
  {
    title: '从帐号',
    dataIndex: 'assetAcc'
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }
]
