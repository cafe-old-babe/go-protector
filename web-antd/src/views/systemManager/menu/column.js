export const Columns = [
  {
    title: '名称',
    dataIndex: 'name'
  },
  {
    title: '类型',
    dataIndex: 'menuTypeName'
  },
  {
    title: '权限标识',
    dataIndex: 'permission'
  },
  {
    title: '目录/菜单显示状态',
    dataIndex: 'hidden',
    scopedSlots: { customRender: 'status' }
  },
  {
    title: '目录/菜单组件名称',
    dataIndex: 'component'
  },
  {
    title: '操作',
    width: 200,
    scopedSlots: { customRender: 'action' }
  }
]
