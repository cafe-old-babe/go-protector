export const columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '策略名称',
    dataIndex: 'name'
  },
  {
    title: '角色类型',
    dataIndex: 'roleType',
    customRender: (text) => text === 0 ? '管理员角色' : '普通角色'
  },
  {
    title: '状态',
    dataIndex: 'status',
    scopedSlots: { customRender: 'status' }
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }
]
