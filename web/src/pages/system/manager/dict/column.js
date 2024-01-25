
const typeColumn = [
    {
        title: '编号',
        width: 60,
        scopedSlots: { customRender: 'serial' }
    },
    {
        title: '字典类型名称',
        dataIndex: "typeName"
    },
    {
        title: '字典类型编码',
        dataIndex: "typeCode"
    },
    {
        title: '操作',
        width: 200,
        scopedSlots: { customRender: 'action' }
    }
]
const dataColumn= [
    {
        title: '序号',
        width: 60,
        scopedSlots: { customRender: 'serial' }
    },{
        title: '数据名称',
        dataIndex: "dataName"
    },
    {
        title: '数据编码',
        dataIndex: "dataCode"
    },
    {
        title: '状态',
        scopedSlots: { customRender: 'status' }
    },
    {
        title: '操作',
        width: 200,
        scopedSlots: { customRender: 'action' }
    }
]



export  default {typeColumn, dataColumn}