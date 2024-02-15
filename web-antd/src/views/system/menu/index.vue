<template>
  <div>
    <a-card :bordered="false" >
      <div class="table-operator">
        <a-button type="primary" @click="loadData">刷新</a-button>
        <a-button type="primary" @click="handleAdd">新建</a-button>
        <a-button type="danger" @click="deleteBatch">批量删除</a-button>
      </div>
      <a-table
        ref="table"
        rowKey="id"
        :loading="loading"
        :columns="columns"
        :data-source="data"
        :rowSelection="rowSelection"
        :pagination="false"
      >
        <span slot="serial" slot-scope="{index}">
          {{ index + 1 }}
        </span>
        <div slot="status" slot-scope="text,current">
          <a v-if="current.menuType!==2">
            <a-tooltip
              placement="left"
              :title="current.hidden===0?'点击隐藏':'点击显示'"
              :get-popup-container="(trigger) => trigger.parentElement">
              <a-tag
                :color="current.hidden===0?'green':'red'"
                @click="changeStatus(current)"
              >
                {{ current.hidden===0?'显示':'隐藏' }}
              </a-tag>
            </a-tooltip>
          </a>
        </div>
        <span slot="action" slot-scope="text, current">
          <a style="margin-right: 8px" @click="editRecord(current)"> <a-icon type="edit" />编辑 </a>
          <a v-action:delete @click="deleteRecord(current.id)"> <a-icon type="delete" />删除 </a>
        </span>
      </a-table>
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
      columns: Columns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      data: [],
      loadData: () => {
        this.loading = true
        return request.post('/menu/list', { }).then((res) => {
          const { code, data, message } = res
          if (code === 200) {
            this.data = data.children
          } else {
            this.$message.error(message)
          }
        }).finally(() => {
          this.loading = false
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
  created() {
    this.loadData()
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
        request.post('/menu/delete', param).then(res => {
          const { code, message } = res
          if (code === 200) {
            this.$message.success(message ?? '删除成功')
            this.loadData()
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
    }
  }
}
</script>

<style scoped></style>
