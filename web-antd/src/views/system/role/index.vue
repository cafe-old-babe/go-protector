<template>
  <div>
    <a-card :bordered="false" title="角色列表" :style="{height:`calc(${windowHeight}px - 210px)`,overflow:'hidden'}">
      <div class="table-page-search-wrapper">
        <a-form layout="inline">
          <a-row :gutter="48">
            <a-col :md="8" :sm="24">
              <a-form-item label="角色名称">
                <a-input v-model="queryParam.roleName" placeholder="" />
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
        :show-pagination="true"
        :showSizeChanger="true"
        :columns="columns"
        :data="loadData"
        :rowSelection="rowSelection"
      >
        <span slot="action" slot-scope="text, current">
          <a :disabled="current.isInner===1 || current.roleType===0" style="margin-right: 8px" @click="permission(current.id)">
            <a-icon type="link"/>授权
          </a>
          <a :disabled="current.isInner===1" style="margin-right: 8px" @click="editRecord(current)">
            <a-icon type="edit" />编辑
          </a>
          <a :disabled="current.isInner===1" @click="deleteRecord(current.id)">
            <a-icon type="delete" />删除
          </a>
        </span>
        <div slot="status" slot-scope="text,current">
          <a :disabled="current.isInner===1">
            <a-tooltip
              placement="left"
              :title="current.status===0?'点击停用':'点击正常'"
              :get-popup-container="(trigger) => trigger.parentElement">
              <a-tag
                :color="current.status===0?'green':'red'"
                @click="changeStatus(current)"
              >
                {{ current.status===0?"正常":'停用' }}
              </a-tag>
            </a-tooltip>
          </a>
        </div>
      </s-table>
    </a-card>
    <Edit
      :visible="editVisible"
      :record="record"
      @close="editClose"
      @ok="editOk"
    />
    <Permission
      :role-id="roleId"
      :visible="permissionVisible"
      @close="editClose"
      @ok="editOk"
    />
  </div>
</template>

<script>
import STable from '@/components/Table'
import { RoleColumns } from './column'
import request from '@/utils/request'
import Edit from './Edit'
import Permission from './Permission.vue'

export default {
  name: 'Type',
  components: { Edit, Permission, STable },
  data() {
    return {
      loading: false,
      editVisible: false,
      permissionVisible: false,
      columns: RoleColumns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      windowHeight: 0,
      roleId: 0,
      loadData: (parameter) => {
        const promise = request.post('/role/page', Object.assign(parameter, this.queryParam)).then((res) => {
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
            disabled: record.isInner === 1, // Column configuration not to be checked
            name: record.name
          }
        })
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
      console.log(record)
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
        request.post('/role/delete', param).then(res => {
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
    changeStatus: function(record) {
      this.loading = true
      request.post(`/role/setStatus/${record.id }/${(record.status ^ 1)}`).then(res => {
        if (res.code === 200) {
          this.$refs.table.refresh(false)
          this.$message.success(res.message)
        } else {
          this.$message.warning(res.message)
        }
      }).finally(() => {
        this.loading = false
      })
    },
    permission: function (roleId) {
      this.roleId = roleId
      this.permissionVisible = true
    }
  }
}
</script>

<style scoped lang="less">

</style>
