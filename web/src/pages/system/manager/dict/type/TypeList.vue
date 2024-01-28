<template>
  <div>
    <a-card :bordered="false" style="margin-bottom: 24px" size="small"  title="字典类型列表">
      <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 14 }" ref="typeForm" :form="searchForm" @submit="handleSearch">
        <a-row :gutter="[8,1]">
          <a-col :span="12">

            <a-form-item label="类型名称">
              <a-input
                v-decorator="[ 'typeName', {
                                  rules: [
                                    {
                                      required: false,
                                      message: 'Input something!',
                                    },
                                  ],
                                },
                              ]" placeholder="请输类型名称"
              />
            </a-form-item>

          </a-col>
          <a-col :span="12">
            <a-form-item label="类型编码">
              <a-input
                v-decorator="[ 'typeCode', {
                                  rules: [
                                    {
                                      required: false,
                                      message: 'Input something!',
                                    },
                                  ],
                                },
                              ]" placeholder="请输入类型编码"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="[8,8]">
          <a-col :span="12" :style="{ textAlign: 'left' }">
            <a-space :size="8">
              <a-button @click="addNew" type="primary">新建</a-button>
              <a-button type="danger" @click="deleteBatch">批量删除</a-button>
<!--              <a-button type="delete">批量删除</a-button>
              <a-dropdown>
                &lt;!&ndash;                    <a-menu @click="handleMenuClick" slot="overlay">&ndash;&gt;
                &lt;!&ndash;                        <a-menu-item key="delete">删除</a-menu-item>&ndash;&gt;
                &lt;!&ndash;                        <a-menu-item key="audit">审批</a-menu-item>&ndash;&gt;
                &lt;!&ndash;                    </a-menu>&ndash;&gt;
                <a-button>
                  更多操作 <a-icon type="down" />
                </a-button>
              </a-dropdown>-->
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
      >
        <div slot="action" slot-scope="{text, record}">
          <a style="margin-right: 8px"
             @click="eyeRecord(record.typeCode)"
          >
            <a-icon type="eye"/>查看
          </a>
          <a style="margin-right: 8px"
             @click="editRecord(record)"
          >
            <a-icon type="edit"/>编辑
          </a>
          <a @click="deleteRecord(record.id)">
            <a-icon type="delete" />删除
          </a>
        </div>

      </StandardTable>

    </a-card>
    <Edit
      :visible="editVisible"
      :record="record"
      @close="editClose"
      @ok="editOk"
    />
  </div>
</template>

<script>
import StandardTable from '@/components/table/StandardTable'
import {request} from '@/utils/request'
import column from "../column";
import Edit from "@/pages/system/manager/dict/type/Edit.vue";

export default {
  name: 'TypeList',
  components: {Edit, StandardTable},
  data() {
    return {
      editVisible: false,
      record: {},
      searchForm: this.$form.createForm(this),
      searchFormParam: {},
      loading: false,
      columns: column.typeColumn,
      dataSource: [],
      selectedRows: [],
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0
      }
    }
  },
  mounted() {
    this.getData()
  },
  computed: {

  },
  methods: {
    editOk() {
      this.editClose()
      this.getData()
    },
    editClose() {
      this.editVisible = false
      this.record = {}
    },
    editRecord(record) {
      this.record = record;
      this.editVisible = true
    },
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
        this.searchFormParam = values
        this.getData()
      });
    },
    getData() {
      this.loading = true

      request('/api/dict/type',  {...this.pagination,... this.searchFormParam }).then(res => {
          const {list, current, pageSize, total} = res?.data?.data ?? {}
          this.dataSource = list
          this.pagination.current = current
          this.pagination.pageSize = pageSize
          this.pagination.total = total
      }).finally(() => this.loading = false)
    },
    handleReset() {
      this.searchForm.resetFields()
    },
    deleteBatch() {
      this.doDelete({ids: this.selectedRows.map(item => item.id)});
    },
    deleteRecord(id) {
      console.log(id, "delete")
      this.doDelete({"ids": [id]})
    },
    doDelete: function (param, confirm = true) {
      console.log((param?.ids ?? [] ).length);
      if ((param?.ids??[]).length <= 0) {
        this.$message.warning("请选择要删除的数据");
        return;
      }

      if (!confirm) {
        this.loading = true;
        request("/api/dict/type/delete", param).then(res => {
          const {code, message} = res?.data ?? {};
          if (code === 200) {
            this.$message.success(message ?? "删除成功");
            this.getData();
          } else {
            this.$message.error(message);
            this.loading = false;
          }
        })
        return;
      }

      this.$confirm({
        title: '确认删除',
        content: '是否删除选择数据?',
        okText: '确认',
        okType: 'danger',
        cancelText: '取消',
        confirmLoading: this.loading,
        onOk: () => this.doDelete(param, false),
      })
    },
    onSelectChange() {
      console.log(this.selectedRows)
      this.$message.info('选中行改变了')
    },
    addNew () {
      this.record = {}
      this.editVisible = true
    },
    eyeRecord(typeCode) {
      this.$emit('updateTypeCode',typeCode)
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
