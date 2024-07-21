export const Columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '审批工单号',
    customRender: (text, record) => record.workNum
  },
  {
    title: '审批类型',
    customRender: (text, record) => record.approveTypeName
  },
  {
    title: '工单状态',
    customRender: (text, record) => record.approveStatusText
  },
  {
    title: '申请人',
    customRender: (text, record) => record.applicantUser.username
  },
  {
    title: '审批人',
    customRender: (text, record) => record.approveUser.username
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }
]
