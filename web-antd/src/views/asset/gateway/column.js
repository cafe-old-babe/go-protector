export const Columns = [
  {
    title: '名称',
    dataIndex: 'agName'
  },
  {
    title: 'IP',
    dataIndex: 'agIp'
  },
  {
    title: '端口',
    dataIndex: 'agPort'
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }
]
