export const columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '审批指令',
    dataIndex: 'cmd'
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }]
