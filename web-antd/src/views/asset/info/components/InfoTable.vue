<template>
  <div>
    <div class="table-page-search-wrapper">
      <a-form layout="inline">
        <a-row :gutter="48">
          <a-col :md="8" :sm="24">
            <a-form-item label="资源名称">
              <a-input v-model="queryParam.assetName" placeholder="" />
            </a-form-item>
          </a-col>
          <a-col :md="8" :sm="24">
            <a-form-item label="资源IP">
              <a-input v-model="queryParam.IP" placeholder="" />
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
      size="default"
      :showSizeChanger="true"
      :columns="columns"
      :autoLoad="false"
      :data="loadData"
      :rowSelection="rowSelection"
    >
      <span slot="serial" slot-scope="{index}">
        {{ index + 1 }}
      </span>
      <span slot="action" slot-scope="text, current">
        <a style="margin-right: 8px" @click="editRecord(current)"> <a-icon type="edit" />编辑 </a>
        <a @click="deleteRecord(current.id)"> <a-icon type="delete" />删除 </a>
      </span>
    </s-table>

    <EditInfo
      :visible="editVisible"
      :record="record"
      @close="editClose"
      @ok="editOk"
    />

  </div></template>

<script>
import STable from '@/components/Table'
import request from '@/utils/request'
import { loadAsset } from '@/api/asset'
import EditInfo from './EditInfo.vue'
import { columns } from './InfoTable.js'

export default {
  name: 'UserList',
  components: { EditInfo, STable },
  data() {
    return {
      windowHeight: 0,
      loading: false,
      editVisible: false,
      editStatusVisible: false,
      columns: columns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      loadData: (parameter) => {
        return loadAsset(Object.assign(parameter, this.queryParam)).then((res) => {
          const { code, data, message } = res
          if (code === 200) {
            return data
          }
          throw new Error(message)
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
  mounted() {
    this.resize()
    window.addEventListener('resize', this.resize, false)
  },
  destroyed () {
    window.removeEventListener('resize', this.resize, false)
  },
  methods: {
    resize () {
      this.windowHeight = document.body.clientHeight
    },
    onSelectChange(selectedRowKeys, selectedRows) {
      this.selectedRowKeys = selectedRowKeys
      this.selectedRows = selectedRows
    },
    editRecord: function (record) {
      this.record = Object.assign({}, {
        'id': record.id,
        'assetName': record.assetName,
        'groupId': record.groupId,
        'ip': record.ip,
        'port': record.port,
        'managerUsername': record.managerUser.username,
        'managerUserId': record.managerUser.id,
        'account': record.rootAcc.account,
        'password': record.rootAcc.password
      })
      this.editVisible = true
    },
    deleteRecord: function (id) {
      this.doDelete({ 'ids': [id] })
    },
    deleteBatch: function() {
      this.doDelete({ ids: this.selectedRows.map(item => item.id) })
    },
    doDelete: function (param, confirm = true) {
      if ((param?.ids ?? []).length <= 0) {
        this.$message.warning('请选择要删除的数据')
        return
      }

      if (!confirm) {
        this.loading = true
        request.post('/asset-info/delete', param).then(res => {
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
    editOk: function () {
      this.editClose()
      this.$refs.table.refresh(false)
    },
    editClose: function () {
      this.record = {}
      this.editVisible = false
      this.editStatusVisible = false
    },
    loadInfo: function (groupIds = [], bool = true) {
      this.queryParam.groupIds = groupIds
      this.$refs.table.refresh(bool)
    }
  }
}
</script>

<style scoped></style>
