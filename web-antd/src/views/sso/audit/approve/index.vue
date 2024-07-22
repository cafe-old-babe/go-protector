<template>
  <div>
    <a-spin :spinning="loading">
      <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 150px)`,overflow:'auto'}">
        <div
          class="table-page-search-wrapper">
          <a-form layout="inline">
            <a-row :gutter="48">
              <a-col :md="6" :sm="24">
                <a-form-item label="工单号">
                  <a-input v-model="queryParam.workNum" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="6" :sm="24">
                <a-form-item label="申请人">
                  <a-input v-model="queryParam.applicantUsername" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="6" :sm="24">
                <a-form-item label="审批人">
                  <a-input v-model="queryParam.approveUsername" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="6" :sm="24">
                <a-button type="primary" @click="$refs.table.refresh(true)">查询</a-button>
                <a-button style="margin-left: 8px" @click="() => (this.queryParam = {})">重置</a-button>
              </a-col>
            </a-row>
          </a-form>
        </div>
        <div class="table-operator">
        </div>
        <s-table
          ref="table"
          rowKey="id"
          size="default"
          :showSizeChanger="true"
          :columns="columns"
          :data="loadData"
          :scroll="{y:'calc(60vh - 30px)'}"
        >
          <template v-slot:action="text, current">
            <span >
              <a :disabled="current.approveStatus!==0" style="margin-right: 8px" @click="doApprove(current)">
                <a-icon type="solution" />处理
              </a>
            </span>
          </template>

        </s-table>
      </a-card>
      <do-approve
        :visible="visible"
        :record="record"
        @close="approveClose"
        @ok="approveOk"
      />

    </a-spin>
  </div>
</template>

<script>
import STable from '@/components/Table'
import { Columns } from './column.js'
import request from '@/utils/request'
import DoApprove from '@/views/sso/audit/approve/DoApprove.vue'
import Edit from '@/views/asset/account/Edit.vue'

export default {
  name: 'Approve',
  components: { Edit, DoApprove, STable },

  data() {
    return {
      windowHeight: 0,
      loading: false,
      visible: false,
      columns: Columns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      loadData: (parameter) => {
        return request.post('/approve-record/page', Object.assign(parameter, this.queryParam)).then((res) => {
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
    approveClose: function () {
      this.record = {}
      this.visible = false
    },
    approveOk: function () {
      this.approveClose()
      this.$refs.table.refresh(false)
    },
    doApprove: function (record) {
      this.record = Object.assign(record, {
        applicantUsername: record.applicantUser.username
      })
      this.visible = true
    },
    refresh: function() {
      this.$refs.table.refresh(true)
    }
  }
}
</script>

<style scoped></style>
