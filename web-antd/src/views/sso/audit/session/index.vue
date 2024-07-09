<template>
  <div>
    <a-spin :spinning="loading">
      <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 150px)`,overflow:'auto'}">
        <div class="table-page-search-wrapper">
          <a-form layout="inline">
            <a-row :gutter="20">
              <a-col :md="5" :sm="5">
                <a-form-item label="资产IP">
                  <a-input v-model="queryParam.assetIp" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="5" :sm="5">
                <a-form-item label="资产名称">
                  <a-input v-model="queryParam.assetName" placeholder="" />
                </a-form-item>
              </a-col>
              <a-col :md="5" :sm="5">
                <a-form-item label="从帐号">
                  <a-input v-model="queryParam.assetAcc" placeholder="" />
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
          <!--          <a-button-->
          <!--            type="primary"-->
          <!--            @click="() => {this.record = {}; this.editVisible = true}">新建</a-button>-->
          <!--          <a-button type="danger" :disabled="disabledDeleteBatch" @click="deleteBatch">批量删除</a-button>-->
          <!--          <a-button @click="dailBatch"> <a-icon type="bell" />批量拨测</a-button>-->
          <!--          <a-button :disabled="disabledCollBatch" @click="collBatch"> <a-icon type="code" />批量采集</a-button>-->
        </div>
        <s-table
          ref="table"
          rowKey="id"
          size="default"
          :show-pagination="true"
          :showSizeChanger="true"
          :columns="columns"
          :data="loadData"
          :scroll="{y:'calc(60vh - 30px)'}"
          :rowSelection="rowSelection"
        >
          <!--          <span slot="action" slot-scope="text, current">-->
          <template v-slot:action="text,current">

            <a :disabled="current.status!=='3'" style="margin-right: 8px" @click="playBack(current)">
              <a-icon type="play-circle" />回放
            </a>
            <a :disabled="current.status!=='2'" style="margin-right: 8px" @click="monitor(current)">
              <a-icon type="video-camera" />云端监控
            </a>
            <!--          </span>-->
          </template>

        </s-table>
      </a-card>
    </a-spin>
    <PlayBack
      :visible="playBackVisible"
      :castData="castData"
      :record="record"
      @close="close"
    >

    </PlayBack>
  </div>
</template>

<script>
import STable from '@/components/Table'
import { Columns } from './column'
import request from '@/utils/request'
import PlayBack from '@/views/sso/audit/session/components/PlayBack.vue'
// import Permission from './Permission.vue'

export default {
  name: 'Session',
  components: { STable, PlayBack },
  data() {
    return {
      loading: false,
      playBackVisible: false,
      permissionVisible: false,
      columns: Columns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      record: {},
      windowHeight: 0,
      castData: '',
      loadData: (parameter) => {
        const promise = request.post('/sso-session/page', Object.assign(parameter, this.queryParam)).then((res) => {
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
    playBack: function (record) {
      request.post(`/sso-session/cast/${record.id}`).then(res => {
        const { code, data, message } = res
        if (code !== 200) {
          this.$message.error(message)
          return
        }
        this.record = record
        this.castData = data
        this.playBackVisible = true
      })
    },
    monitor: function (record) {
      const assign = Object.assign({}, {
        id: record.id,
        uri: '/api/ws/sso-session/monitor/',
        send: false,
        initMsg: '\x1B[1;3;31m正在连接,请稍后\x1B[0m $ ',
        title: '云端监控'
      })
      localStorage.setItem('ssoTerminal', JSON.stringify(assign))
      // window.open(`/sso-terminal?id=${data.id}`)
      window.open(`/sso-terminal`)
    },
    close: function () {
      this.playBackVisible = false
      this.record = {}
      this.castData = ''
    }

  }
}
</script>

<style scoped lang="less">

</style>
