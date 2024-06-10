export const assetColumns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '资产名称',
    dataIndex: 'assetName'
  },
  {
    title: '资源组',
    customRender: (text, record) => record.assetGroup.name
  },
  {
    title: 'IP',
    dataIndex: 'ip'
  },
  {
    title: '端口',
    dataIndex: 'port'
  },
  {
    title: '资产管理员',
    customRender: (text, record) => record.managerUser.username
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }]

export const accountColumns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '从帐号',
    customRender: (text, record, index) => record.assetAccount.account

  },
  {
    title: '从帐号类型',
    customRender: (text, record, index) => record.assetAccount.accountTypeText
  },
  {
    title: '从帐号状态',
    customRender: (text, record, index) => record.assetAccount.accountStatusText

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
