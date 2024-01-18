
const typeColumn = [
    {
        title: '编号',
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
]
const dataColumn= [
    {
        title: '序号',
        scopedSlots: { customRender: 'serial' }
    },{
        title: '字典名称',
        dataIndex: "typeName"
    },
    {
        title: '字典编码',
        dataIndex: "typeCode"
    },
]



export  default {typeColumn, dataColumn}