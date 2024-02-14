export const columns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '用户名',
    dataIndex: 'username'
  },
  {
    title: '登录名',
    dataIndex: 'loginName'
  },
  {
    title: '性别',
    dataIndex: 'sex',
    customRender: (text, record, index) => text === 0 ? '女' : '男'
  },
  {
    title: '部门',
    dataIndex: 'deptName'
  },
  {
    title: '状态',
    dataIndex: 'userStatus',
    scopedSlots: { customRender: 'status' }
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }]
