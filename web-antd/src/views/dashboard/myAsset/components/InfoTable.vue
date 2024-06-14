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
      </div>
      <s-table
        ref="table"
        rowKey="id"
        size="default"
        :showSizeChanger="true"
        :columns="columns"
        :autoLoad="false"
        :data="loadData"

      >
        <span slot="serial" slot-scope="{index}">
          {{ index + 1 }}
        </span>
        <span slot="action" slot-scope="text, current">
          <a style="margin-right: 8px" @click="connect(current)"> <a-icon type="link" />连接 </a>
        </span>

      </s-table>
      <SelectAccount
        :visible="this.visible"
        :accounts="this.accounts"
        @close="selectClose"
        @login="createSsoSession"
      >

      </SelectAccount>
    </a-spin>
  </div>
</template>

<script>
import STable from '@/components/Table'
import { loadAuthAsset } from '@/api/asset'
import { assetColumns } from './component.js'
import SelectAccount from './SelectAccount.vue'
import request from '@/utils/request'

export default {
  name: 'MyAsset',
  components: { SelectAccount, STable },
  data() {
    return {
      windowHeight: 0,
      loading: false,
      visible: false,
      columns: assetColumns,
      queryParam: {},
      selectedRowKeys: [],
      selectedRows: [],
      accounts: [],
      record: {},
      loadData: (parameter) => {
        return loadAuthAsset(Object.assign(parameter, this.queryParam)).then((res) => {
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
    connect: function (record) {
      this.loading = true
      request.post('/asset-account/auth/' + record.id).then((res) => {
        const { code, data, message } = res
        if (code !== 200) {
          this.$message.error(message)
          return
        }
        if (data.length <= 0) {
          this.$message.error('无可用从帐号')
          return
        }
        if (data.length === 1) {
          this.createSsoSession(data[0])
          return
        }
        this.accounts = data
        this.visible = true
      }).finally(() => {
        this.loading = false
      })
    },
    selectClose: function () {
      this.record = {}
      this.visible = false
    },
    loadInfo: function (groupIds = [], bool = true) {
      this.queryParam.groupIds = groupIds
      this.$refs.table.refresh(bool)
    },
    createSsoSession(data) {
      this.selectClose()
      console.log('createSsoSession', data)
      request.post('/sso-session/create/' + data.id).then(res => {
        const { code, message, data } = res
        if (code === 200) {
          console.log(data)
          window.open(`/sso-terminal?id=${data.id}`)
          return
        }
        this.$message.error(message)
      })
      // todo createSsoSession
    }
  }
}
</script>

<style scoped></style>
