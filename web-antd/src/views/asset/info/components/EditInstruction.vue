<script>
import STable from '@/components/Table'
import { columns } from './Instruction'
import request from '@/utils/request'

export default {
  name: 'EditInstruction',
  components: { STable },
  data() {
    return {
      loading: false,
      columns: columns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      loadData: (parameter) => {
        this.queryParam.assetId = this.assetId
        return request.post(`/approve-cmd/page`, Object.assign(parameter, this.queryParam)).then((res) => {
          const { code, data, message } = res
          if (code === 200) {
            return data
          }
          return new Error(message)
        })
      }
    }
  },
  watch: {
    assetId(val) {
      if (!val) {
        return
      }
      this.$nextTick(() => {
        this.$refs.instructionTable.refresh(true)
      })
    }
  },
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    assetId: {
      type: Number,
      require: true,
      default: null
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
        request.post('/approve-cmd/delete', param).then(res => {
          const { code, message } = res
          if (code === 200) {
            this.$message.success(message ?? '删除成功')
            this.$refs.instructionTable.refresh(false)
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
    save: function () {
      request.post('/approve-cmd/save', Object.assign(this.queryParam, {
        'assetId': this.assetId
      })).then(res => {
        const { code, message } = res
        if (code === 200) {
          this.$message.success(message)
          this.queryParam.cmd = ''
          this.$refs.instructionTable.refresh(true)
          return
        }
        this.$message.error(message)
      })
    },
    handleCancel: function () {
      this.queryParam.cmd = ''
      this.$emit('close')
    }

  }
}
</script>

<template>
  <a-modal
    :visible="visible"
    title="管理审批指令"
    :closable="true"
    :maskClosable="false"
    @cancel="handleCancel"
    :dialog-style="{ top: '20px' }"
    :width="800">
    <template v-slot:footer>
      <a-button key="back" @click="handleCancel">
        完成
      </a-button>
    </template>
    <a-card :bordered="false" >
      <div class="table-page-search-wrapper">
        <a-form layout="inline">
          <a-row :gutter="48">
            <a-col :md="8" :sm="24">
              <a-form-item label="指令">
                <a-input v-model="queryParam.cmd" placeholder="" />
              </a-form-item>
            </a-col>

            <a-col :md="16" :sm="24">
              <a-button type="primary" @click="$refs.instructionTable.refresh(true)">查询</a-button>
              <a-button style="margin-left: 8px" @click="save">添加</a-button>
              <a-button style="margin-left: 8px"type="danger" @click="deleteBatch">批量删除</a-button>
            </a-col>
          </a-row>
        </a-form>
      </div>
      <div class="table-operator">

      </div>
      <s-table
        ref="instructionTable"
        rowKey="id"
        :loading="loading"
        :columns="columns"
        :data="loadData"
        :auto-load="false"
        :rowSelection="rowSelection"
        :showSizeChanger="true"
        :scroll="{y:'calc(50vh - 30px)'}"
      >
        <span slot="serial" slot-scope="{index}">
          {{ index + 1 }}
        </span>
        <span slot="action" slot-scope="text, current">
          <a @click="deleteRecord(current.id)"> <a-icon type="delete" />删除 </a>
        </span>
      </s-table>
    </a-card>

  </a-modal>
</template>

<style scoped lang="less">

</style>
