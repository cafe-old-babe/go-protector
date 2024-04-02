<template>
  <div>
    <a-card :bordered="false" >
      <div class="table-page-search-wrapper">
        <a-form layout="inline">
          <a-row :gutter="48">
            <a-col :md="8" :sm="24">
              <a-form-item label="名称">
                <a-input v-model="queryParam.anName" placeholder="" />
              </a-form-item>
            </a-col>
            <a-col :md="8" :sm="24">
              <a-form-item label="IP">
                <a-input v-model="queryParam.anIp" placeholder="" />
              </a-form-item>
            </a-col>

            <a-col :md="8" :sm="24">
              <a-button type="primary" @click="$refs.table.refresh(true)">查询</a-button>
              <a-button style="margin-left: 8px" @click="() => (this.queryParam = {})">重置</a-button>
            </a-col>
          </a-row>
        </a-form>
      </div>
      <div class="table-operator">
        <a-button
          type="primary"
          @click="() => {this.record = {}; this.editVisible = true}">新建</a-button>
        <a-button type="danger" @click="deleteBatch">批量删除</a-button>
      </div>
      <s-table
        ref="table"
        rowKey="id"
        :loading="loading"
        :columns="columns"
        :data="loadData"
        :rowSelection="rowSelection"
        :showSizeChanger="true"
      >
        <span slot="serial" slot-scope="{index}">
          {{ index + 1 }}
        </span>
        <span slot="action" slot-scope="text, current">
          <a style="margin-right: 8px" @click="editRecord(current)"> <a-icon type="edit" />编辑 </a>
          <a @click="deleteRecord(current.id)"> <a-icon type="delete" />删除 </a>
        </span>
      </s-table>
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
import STable from '@/components/Table'
import { Columns } from './column'
import request from '@/utils/request'
import Edit from './edit'
export default {
  name: 'Type',
  components: { STable, Edit },
  data() {
    return {
      loading: false,
      editVisible: false,
      columns: Columns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      loadData: (parameter) => {
        return request.post('/asset-network/page', parameter).then((res) => {
          const { code, data, message } = res
          if (code === 200) {
            return data
          }
          return new Error(message)
        })
      }
    }
  },
  computed: {
    rowSelection() {
      return {
        selectedRowKeys: this.selectedRowKeys,
        onChange: this.onSelectChange
      }
    }
  },
  methods: {
    onSelectChange(selectedRowKeys, selectedRows) {
      this.selectedRowKeys = selectedRowKeys
      this.selectedRows = selectedRows
    },
    editRecord: function (record) {
      this.record = Object.assign({}, record)
      this.editVisible = true
    },
    deleteRecord: function (id) {
      this.doDelete({ ids: [id] })
    },
    deleteBatch: function() {
      // 将 this.selectRows的id提取出来
      this.doDelete({ ids: this.selectedRows.map(item => item.id) })
    },
    doDelete: function (param, confirm = true) {
      if (param.ids.length <= 0) {
        this.$message.warning('请选择要删除的数据')
        return
      }

      if (!confirm) {
        this.loading = true
        request.post('/asset-network/delete', param).then(res => {
          const { code, message } = res
          if (code === 200) {
            this.$message.success(message ?? '删除成功')
            this.$refs.table.refresh(false)
          } else {
            this.$message.error(message)
            this.loading = false
          }
        })
        return
      }
      this.$confirm({
        title: '确认删除',
        content: '已使用的网络域无法删除,是否删除选择数据?',
        okText: '确认',
        okType: 'danger',
        cancelText: '取消',
        confirmLoading: this.loading,
        onOk: () => this.doDelete(param, false)
      })
    },
    editOk: function () {
      this.editClose()
      this.$refs.table.refresh(false)
    },
    editClose: function () {
      this.editVisible = false
      this.record = {}
    }

  }
}
</script>

<style scoped></style>
