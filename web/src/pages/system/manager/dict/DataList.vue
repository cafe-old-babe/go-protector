<template>
    <div>
        <a-card :bordered="false" style="margin-bottom: 24px" size="small"  title="字典数据列表">
            <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 14 }" :form="searchForm" @submit="handleSearch">
                <a-row :gutter="[8,1]">
                    <a-col :span="12">

                        <a-form-item label="数据名称">
                            <a-input
                                v-decorator="[ 'dataName', {
                                  rules: [
                                    {
                                      required: false,
                                      message: 'Input something!',
                                    },
                                  ],
                                },
                              ]" placeholder="请输入数据名称"
                            />
                        </a-form-item>

                    </a-col>
                    <a-col :span="12">
                        <a-form-item label="数据编码">
                            <a-input
                                v-decorator="[ 'dataCode', {
                                  rules: [
                                    {
                                      required: false,
                                      message: 'Input something!',
                                    },
                                  ],
                                },
                              ]" placeholder="请输入数据编码"
                            />
                        </a-form-item>
                    </a-col>
                </a-row>

                <a-row :gutter="[8,8]">
                    <a-col :span="12" :style="{ textAlign: 'left' }">
                        <a-space :size="8">
                            <a-button @click="addNew" type="primary">新建</a-button>
                            <a-button >批量操作</a-button>
                            <a-dropdown>
                                <!--                    <a-menu @click="handleMenuClick" slot="overlay">-->
                                <!--                        <a-menu-item key="delete">删除</a-menu-item>-->
                                <!--                        <a-menu-item key="audit">审批</a-menu-item>-->
                                <!--                    </a-menu>-->
                                <a-button>
                                    更多操作 <a-icon type="down" />
                                </a-button>
                            </a-dropdown>
                        </a-space>
                    </a-col>
                    <a-col :span="12" :style="{ textAlign: 'right' }">
                        <a-space :size="8">
                            <a-button type="primary" html-type="submit">
                                查询
                            </a-button>
                            <a-button @click="handleReset">
                                Clear
                            </a-button>

                        </a-space>
                    </a-col>
                </a-row>
            </a-form>
        </a-card>

        <a-card :bordered="false" style="">
            <StandardTable
                :columns="columns"
                :loading="loading"
                rowKey="id"
                :data-source=dataSource
                :selectedRows.sync="selectedRows"
                @selectedRowChange="onSelectChange"
                :pagination="{...pagination, onChange: onPageChange}"
            />

        </a-card>
    </div>
</template>

<script>
import StandardTable from '@/components/table/StandardTable'
// import {request} from '@/utils/request'
import column from "./column";

export default {
    name: 'DataList',
    components: {StandardTable},
    props: {
        typeCode: String
    },
    data() {
        return {
            targetCount: 12,
            expand: false,
            searchForm: this.$form.createForm(this),
            searchFormParam: {},
            loading: false,
            columns: column.dataColumn,
            dataSource: [],
            selectedRows: [],
            pagination: {
                current: 1,
                pageSize: 10,
                total: 0
            }
        }
    },
    watch: {
        typeCode(typeCode) {
            alert(typeCode)
            console.log(typeCode)
        }
    },
    mounted() {
        this.getData()
    },
    computed: {

    },
    methods: {
        onPageChange(page, pageSize) {
            this.pagination.current = page
            this.pagination.pageSize = pageSize
            this.getData()
        },

        handleSearch(e) {
            e.preventDefault();
            this.searchForm.validateFields((error, values) => {
                if (error) {
                    console.log('error', error)
                    return
                }
                console.log('Received values of form: ', values);
                this.getData()
            });
        },
        getData() {
            this.loading = true
            setTimeout(() => {
                this.dataSource = [
                ];
              // this.dataSource = [
              //   {"id":'1',"d1": "data1"},
              //   {"id":'2',"d1": "data1"},
              //   {"id":'3',"d1": "data2"},
              //   {"id":'4',"d1": "data2"},
              //   {"id":'5',"d1": "data2"},
              //   {"id":'6',"d1": "data2"},
              //   {"id":'7',"d1": "data2"},
              //   {"id":'8',"d1": "data2"},
              //   {"id":'9',"d1": "data2"},
              //   {"id":'0',"d1": "data2"},
              // ];
                this.pagination.current = 1
                this.pagination.pageSize = 10
                this.pagination.total = 0
                this.loading = false
            }, 500)

            // request(process.env.VUE_APP_API_BASE_URL + '/list', 'get', {page: this.pagination.current,
            //     pageSize: this.pagination.pageSize}).then(res => {
            //     const {list, page, pageSize, total} = res?.data?.data ?? {}
            //     this.dataSource = list
            //     this.pagination.current = page
            //     this.pagination.pageSize = pageSize
            //     this.pagination.total = total
            // })
        },
        handleReset() {
            this.searchForm.resetFields()
        },
        deleteRecord(key) {
            this.dataSource = this.dataSource.filter(item => item.key !== key)
            this.selectedRows = this.selectedRows.filter(item => item.key !== key)
        },
        remove () {
            this.dataSource = this.dataSource.filter(item => this.selectedRows.findIndex(row => row.key === item.key) === -1)
            this.selectedRows = []
        },
        onClear() {
            this.$message.info('您清空了勾选的所有行')
        },


        onSelectChange() {
            console.log(this.selectedRows)
            this.$message.info('选中行改变了')
        },
        addNew () {
           this.$message.info("addNew")
        },
    }
}
</script>

<style lang="less" scoped>
.search{
    //margin-bottom: 54px;
}
.fold{
    width: calc(100% - 216px);
    display: inline-block
}
.operator{
    margin: 10px;
    //margin-bottom: 18px;
}
@media screen and (max-width: 900px) {
    .fold {
        width: 100%;
    }
}

.ant-advanced-search-form {
    //padding: 12px;
    //background: #fbfbfb;
    //border: 1px solid #d9d9d9;
    //border-radius: 6px;
}

.ant-advanced-search-form .ant-form-item {
    display: flex;
    //display: flow;
}

.ant-advanced-search-form .ant-form-item-control-wrapper {
    flex: 1;
}

#components-form-demo-advanced-search .ant-form {
    max-width: none;
}

</style>
