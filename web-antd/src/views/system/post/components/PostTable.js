export const columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '岗位名称',
    dataIndex: 'name'
  },
  {
    title: '岗位代码',
    dataIndex: 'code'
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }]
