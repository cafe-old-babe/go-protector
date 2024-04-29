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
      :rowSelection="rowSelection"
    >
      <span slot="serial" slot-scope="{index}">
        {{ index + 1 }}
      </span>
    </s-table>

  </div></template>

<script>
import STable from '@/components/Table'
import { loadAsset } from '@/api/asset'
import { columns } from './InfoTable.js'

export default {
  name: 'UserList',
  components: { STable },
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
        onChange: this.onSelectChange,
        type: 'radio'
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
      this.$emit('change', selectedRowKeys, selectedRows)
    },
    loadInfo: function (groupIds = [], bool = true) {
      this.queryParam.groupIds = groupIds
      this.$refs.table.refresh(bool)
    }
  }
}
</script>

<style scoped></style>
