export const Columns = [
  {
    title: '名称',
    dataIndex: 'anName'
  },
  {
    title: 'IP',
    dataIndex: 'anIp'
  },
  {
    title: '端口',
    dataIndex: 'anPort'
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }
]
