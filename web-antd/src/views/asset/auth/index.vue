<template>
  <div>
    <a-card :bordered="false" title="授权列表" :style="{height:`calc(${windowHeight}px - 150px)`,overflow:'auto'}">
      <div class="table-page-search-wrapper">
        <a-form layout="inline">
          <a-row :gutter="20">
            <a-col :md="5" :sm="5">
              <a-form-item label="主帐号">
                <a-input v-model="queryParam.userAcc" placeholder="" />
              </a-form-item>
            </a-col>
            <a-col :md="5" :sm="5">
              <a-form-item label="从帐号">
                <a-input v-model="queryParam.assetAcc" placeholder="" />
              </a-form-item>
            </a-col>
            <a-col :md="5" :sm="5">
              <a-form-item label="资产IP">
                <a-input v-model="queryParam.assetIp" placeholder="" />
              </a-form-item>
            </a-col>

            <a-col :md="5" :sm="5">
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
        <a-button type="danger" :disabled="disabledDeleteBatch" @click="deleteBatch">批量删除</a-button>
        <a-dropdown>
          <a-menu slot="overlay" @click="handleMenuClick">
            <a-menu-item key="template"> <a-icon type="download" />导出模板 </a-menu-item>
            <a-menu-item key="import">
              <a-upload
                name="file"
                :multiple="true"
                :customRequest="importData"
                accept=".xlsx"
                :showUploadList="false"
              >
                <a-icon
                  type="cloud-upload" /> 导入授权
              </a-upload>
            </a-menu-item>
            <a-menu-item key="data"> <a-icon type="cloud-download" />导出授权 </a-menu-item>
          </a-menu>
          <a-button type="primary" style="margin-left: 8px"> 导入导出 <a-icon type="down" /> </a-button>
        </a-dropdown>
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
      >
        <span slot="action" slot-scope="text, current">
          <a style="margin-right: 8px" @click="editRecord(current)">
            <a-icon type="edit" />编辑
          </a>
          <a @click="deleteRecord(current.id)">
            <a-icon type="delete" />删除
          </a>
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
import Edit from './Edit'
// import Permission from './Permission.vue'

export default {
  name: 'Account',
  components: { Edit, STable },
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
        const promise = request.post('/asset-auth/page', Object.assign(parameter, this.queryParam)).then((res) => {
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
      this.record = record
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
        request.post('/asset-auth/delete', param).then(res => {
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
    handleMenuClick: function (e) {
      if (e.key === 'import') {
        // 上传
        return
      }
      // 下载
      this.loading = true
      request.post('/asset-auth/excel/' + e.key, this.queryParam, { responseType: 'blob' })
      .then(this.processRes).finally(() => {
        this.loading = false
      })
    },
    importData: function (data) {
      this.loading = true
      // console.log(data)
      const formData = new FormData()
      // const fileList = {}
      formData.append('file', data.file)
      request.post('/asset-auth/excel/import', formData, {
        responseType: 'blob',
        timeout: 1000 * 60 * 5,
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      }).then(this.processRes).catch(() => {
        this.$message.error('导入失败！')
      }).finally(() => {
        this.$refs.table.refresh(false)
        this.loading = false
      })
    },
    processRes: function (res) {
      if (res?.type === 'application/json') {
        const reader = new FileReader()
        const self = this
        reader.onload = function (event) {
          const content = JSON.parse(reader.result)
          if (content && content.code !== 200) {
            self.$message.error(content.message)
          }
        }
        reader.readAsText(res)

        return
      }
      const element = res.headers['content-disposition'].split('=')[1]
        .replaceAll('"', '')
        .replaceAll("UTF-8''", '')
      const fileName = decodeURI(element)
      const blob = res.data

      // 创建新的URL并指向File对象或者Blob对象的地址
      const url = window.URL.createObjectURL(new Blob([blob]))
      // 创建a标签，用于跳转至下载链接
      const link = document.createElement('a')
      link.style.display = 'none'
      link.href = url
      link.setAttribute('download', fileName)
      // 挂载a标签
      document.body.appendChild(link)
      // 点击
      link.click()
      // 删除
      document.body.removeChild(link)
      // 释放blob URL地址
      window.URL.revokeObjectURL(url)
    }
  }
}
</script>

<style scoped lang="less">

</style>
