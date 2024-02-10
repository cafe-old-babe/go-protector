<template>
  <div>
    <a-card :bordered="false" :title="this.typeCode ? `[${typeCode}]字典数据列表` : '字典数据列表'">
      <div class="table-page-search-wrapper">
        <a-form layout="inline">
          <a-row :gutter="48">
            <a-col :md="8" :sm="24">
              <a-form-item label="数据名称">
                <a-input v-model="queryParam.typeName" placeholder="" />
              </a-form-item>
            </a-col>
            <a-col :md="8" :sm="24">
              <a-form-item label="数据编码">
                <a-input v-model="queryParam.typeCode" placeholder="" />
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
        <a-button type="primary" @click="handleAdd">新建</a-button>
        <a-button type="danger" @click="deleteBatch">批量删除</a-button>
      </div>
      <s-table
        ref="table"
        rowKey="id"
        size="default"
        :showPagination="true"
        :showSizeChanger="true"
        :columns="columns"
        :data="loadData"
        :rowSelection="rowSelection"
        :autoLoad="false"
      >
        <span slot="serial" slot-scope="{index}">
          {{ index + 1 }}
        </span>
        <div slot="status" slot-scope="text,current">
          <a >
            <a-tooltip
              placement="left"
              :title="current.status==='0'?'点击停用':'点击锁定'"
              :get-popup-container="getPopupContainer">
              <a-tag
                :color="current.status===0?'green':'red'"
                @click="changeStatus(current)"
              >
                {{ current.statusText }}
              </a-tag>
            </a-tooltip>
          </a>
        </div>
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
import column from '../column'
import request from '@/utils/request'
import Edit from './edit'

export default {
  name: 'Type',
  components: { STable, Edit },
  props: {
    typeCode: {
      type: String,
      default: ''
    }
  },
  watch: {
    typeCode() {
      this.$refs.table.refresh(true)
    }
  },
  data() {
    return {
      loading: false,
      editVisible: false,
      columns: column.dataColumn,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      loadData: (parameter) => {
        let promise
        if (!this.typeCode) {
          this.$message.warning('请选择字典类型')
          promise = Promise.resolve({
            data: [],
            pageNo: 1,
            pageSize: 10,
            totalCount: 0,
            totalPage: 0
          })
        } else {
          promise = request.post('/dict/data',
            Object.assign(parameter, this.queryParam, { typeCode: this.typeCode })).then((res) => {
            const { code, data, message } = res
            if (code === 200) {
              return data
            }
            this.$message.error(message)
          })
        }

        return promise
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
      console.log('editRecord', record)
      this.record = record
      this.editVisible = true
    },
    deleteRecord: function (id) {
      console.log(id, 'delete')
      this.doDelete({ 'ids': [id] })
    },
    deleteBatch: function() {
      // 将 this.selectRows的id提取出来
      this.doDelete({ ids: this.selectedRows.map(item => item.id) })
    },
    doDelete: function (param, confirm = true) {
      if (param?.ids ?? [].length <= 0) {
        this.$message.warning('请选择要删除的数据')
        return
      }

      if (!confirm) {
        this.loading = true
        request.post('/dict/data/delete', param).then(res => {
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
        content: '是否删除选择数据?',
        okText: '确认',
        okType: 'danger',
        cancelText: '取消',
        confirmLoading: this.loading,
        onOk: () => this.doDelete(param, false)
      })
    },
    handleAdd: function () {
      if (!this.typeCode) {
        this.$message.warning('请选择字典类型')
        return
      }
      this.record = { 'typeCode': this.typeCode }
      this.editVisible = true
    },
    editOk: function () {
      this.editClose()
      this.$refs.table.refresh(false)
    },
    editClose: function () {
      this.editVisible = false
      this.record = {}
    },
    changeStatus: function(record) {
      const url = '/dict/data/' + record.id + '/' + (record.status ^ 1)
      this.loading = true
      request.post(url).then(res => {
        // console.log(res)
        const resData = res?.data ?? {}
        if (resData.code === 200) {
          this.getData()
        } else {
          this.$message.warning(resData.message)
        }
      }).finally(() => {
        this.loading = false
      })
      // console.log(record, 'changeStatus')
    },
    getPopupContainer(trigger) {
      return trigger.parentElement
    }
  }
}
</script>

<style scoped></style>
