export const RoleColumns = [
  {
    title: '序号',
    width: 60,
    // scopedSlots: { customRender: 'serial' }
    customRender: (text, record, index) => index + 1
  },
  {
    title: '角色名称',
    dataIndex: 'roleName'
  },
  {
    title: '角色类型',
    dataIndex: 'roleTypeName',
    customRender: (text) => text === 0 ? '管理员角色' : '普通角色'
  },
  {
    title: '是否内置',
    dataIndex: 'isInner',
    customRender: (text) => text === 0 ? '否' : '是'
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
export const MenuColumns = [
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
  }
]
