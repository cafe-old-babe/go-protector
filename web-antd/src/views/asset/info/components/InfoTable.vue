<template>
  <div>
    <a-spin :spinning="loading">
      <div
        class="table-page-search-wrapper">
        <a-form layout="inline">
          <a-row :gutter="48">
            <a-col :md="8" :sm="24">
              <a-form-item label="资源名称">
                <a-input v-model="queryParam.assetName" placeholder="" />
              </a-form-item>
            </a-col>
            <a-col :md="8" :sm="24">
              <a-form-item label="资源IP">
                <a-input v-model="queryParam.ip" placeholder="" />
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
        <a-button @click="dailBatch"> <a-icon type="bell" />批量拨测</a-button>
        <a-button :disabled="disabledCollBatch" @click="collBatch"> <a-icon type="code"/>批量采集</a-button>
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
          <a @click="editInstruction(current)"> <a-icon type="code" />审批指令管理 </a>

        </span>
        <span slot="dailStatus" slot-scope="text,current">
          <template>
            <a-tooltip
              placement="left"
              :title="current.rootAcc.dailMsg"
              :get-popup-container="(trigger) => trigger.parentElement">
              <a-tag
                :color="dailColor(current.rootAcc.dailStatus)" >
                {{
                  current.rootAcc.dailStatusText
                }}
              </a-tag>
            </a-tooltip>
          </template>
        </span>
      </s-table>

      <EditInfo
        :visible="editVisible"
        :record="record"
        @close="editClose"
        @ok="editOk"
      />
      <EditInstruction
        :visible="instructionVisible"
        :asset-id="record.id"
        @close="editClose"
      />
    </a-spin>
  </div>
</template>

<script>
import STable from '@/components/Table'
import request from '@/utils/request'
import { loadAsset } from '@/api/asset'
import EditInfo from './EditInfo.vue'
import { columns } from './InfoTable.js'
import EditInstruction from '@/views/asset/info/components/EditInstruction.vue'

export default {
  name: 'UserList',
  components: { EditInstruction, EditInfo, STable },
  data() {
    return {
      windowHeight: 0,
      loading: false,
      editVisible: false,
      instructionVisible: false,
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
    },
    disabledCollBatch() {
      return this.selectedRows.some(item => item.rootAcc.dailStatus !== '1')
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
        'password': record.rootAcc.password,
        'gatewayId': record.gatewayId <= 0 ? null : record.gatewayId
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
      this.instructionVisible = false
      this.editStatusVisible = false
    },
    loadInfo: function (groupIds = [], bool = true) {
      this.queryParam.groupIds = groupIds
      this.$refs.table.refresh(bool)
    },
    dailBatch: function () {
      const param = { ids: this.selectedRows.map(item => item.id) }
      this.loading = true
      request.post('/asset-info/dial/asset', param).then(res => {
        const { code, message } = res
        if (code === 200) {
          this.$message.success(message)
          this.$refs.table.refresh(false)
          this.selectedRowKeys = []
          this.selectedRows = []
        } else {
          this.$message.error(message)
        }
      }).finally(() => {
        this.loading = false
      })
    },
    dailColor(dailStatus) {
      switch (dailStatus) {
        case '0':
          return 'red'
        case '1':
          return 'green'
        default:

          return 'grey'
      }
    },
    collBatch: function () {
      const param = { ids: this.selectedRows.map(item => item.id) }
      this.loading = true
      request.post('/asset-info/collectors/asset', param).then(res => {
        const { code, message } = res
        if (code === 200) {
          this.$message.success(message)
          this.$refs.table.refresh(false)
          this.selectedRowKeys = []
          this.selectedRows = []
        } else {
          this.$message.error(message)
        }
      }).finally(() => {
        this.loading = false
      })
    },
    editInstruction: function (param) {
      this.record = param
      this.instructionVisible = true
    }

  }
}
</script>

<style scoped></style>
