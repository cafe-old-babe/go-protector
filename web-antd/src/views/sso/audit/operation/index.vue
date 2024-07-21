<template>
  <div>
    <a-spin :spinning="loading">
      <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 140px)`,overflow:'auto'}">
        <div class="table-page-search-wrapper">
          <a-form layout="inline">
            <a-row :gutter="20">
              <a-col :md="4" :sm="5">
                <a-form-item label="资产IP">
                  <a-input v-model="queryParam.assetIp" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="4" :sm="5">
                <a-form-item label="资产名称">
                  <a-input v-model="queryParam.assetName" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="4" :sm="5">
                <a-form-item label="从帐号">
                  <a-input v-model="queryParam.account" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="4" :sm="5">
                <a-form-item label="执行指令">
                  <a-input v-model="queryParam.cmd" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="4" :sm="5">
                <a-button type="primary" @click="$refs.table.refresh(true)">查询</a-button>
                <a-button style="margin-left: 8px" @click="() => (this.queryParam = {})">重置</a-button>
              </a-col>
            </a-row>
          </a-form>
        </div>

        <s-table
          ref="table"
          rowKey="id"
          size="default"
          :show-pagination="true"
          :showSizeChanger="true"
          :columns="columns"
          :data="loadData"
          :rowSelection="rowSelection"
          :scroll="{y:'calc(60vh - 20px)'}"
        >

        </s-table>
      </a-card>
    </a-spin>
  </div>
</template>

<script>
import STable from '@/components/Table'
import { Columns } from './column'
import request from '@/utils/request'
// import Permission from './Permission.vue'

export default {
  name: 'AuditOperation',
  components: { STable },
  data() {
    return {
      loading: false,
      editVisible: false,
      permissionVisible: false,
      columns: Columns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      windowHeight: 0,
      roleId: 0,
      loadData: (parameter) => {
        const promise = request.post('/sso-operation/page', Object.assign(parameter, this.queryParam)).then((res) => {
          const { code, data, message } = res
          if (code === 200) {
            return data
          }
          this.$message.error(message)
        })
        return promise.catch((error) => {
          this.$message.error(error.message)
          return {
            data: [],
            pageNo: 1,
            pageSize: 10,
            totalCount: 0,
            totalPage: 0
          }
        })
      }
    }
  },
  computed: {
    rowSelection() {
      return {
        selectedRowKeys: this.selectedRowKeys,
        onChange: this.onSelectChange,
        getCheckboxProps: record => ({
          props: {
            // disabled: record.accountType === '0', // Column configuration not to be checked
            name: record.name
          }
        })
      }
    },
    disabledDeleteBatch() {
       return this.selectedRows.some(item => item.accountType === '0')
    },
    disabledCollBatch() {
      return this.selectedRows.some(item => item.dailStatus !== '1')
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
        'assetId': record.assetBasic.id,
        'assetInfoName': `${record.assetBasic.assetName}(${record.assetBasic.ip})`,
        'accountType': record.accountType,
        'account': record.account,
        'password': record.password
      })
      this.editVisible = true
    },
    deleteRecord: function (id) {
      this.doDelete({ 'ids': [id] })
    },
    deleteBatch: function() {
      // 将 this.selectRows的id提取出来
      this.doDelete({ ids: this.selectedRows.map(item => item.id) })
    },
    doDelete: function (param, confirm = true) {
      if ((param?.ids ?? []).length <= 0) {
        this.$message.warning('请选择要删除的数据')
        return
      }

      if (!confirm) {
        this.loading = true
        request.post('/asset-account/delete', param).then(res => {
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
      this.permissionVisible = false
      this.roleId = 0
    },
    dailBatch: function () {
      const param = { ids: this.selectedRows.map(item => item.id) }
      this.loading = true
      request.post('/asset-info/dial/account', param).then(res => {
        const { code, message } = res
        if (code === 200) {
          this.$message.success(message)
          this.selectedRowKeys = []
          this.selectedRows = []
          this.$refs.table.refresh(false)
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
      request.post('/asset-info/collectors/account', param).then(res => {
        const { code, message } = res
        if (code === 200) {
          this.$message.success(message)
          this.selectedRowKeys = []
          this.selectedRows = []
          this.$refs.table.refresh(false)
        } else {
          this.$message.error(message)
        }
      }).finally(() => {
        this.loading = false
      })
    }
  }
}
</script>

<style scoped lang="less">

</style>
