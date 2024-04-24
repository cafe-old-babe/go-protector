<template>
  <div>
    <div class="table-page-search-wrapper">
      <a-form layout="inline">
        <a-row :gutter="48">
          <a-col :md="8" :sm="24">
            <a-form-item label="用户名">
              <a-input v-model="queryParam.username" placeholder="" />
            </a-form-item>
          </a-col>
          <a-col :md="8" :sm="24">
            <a-form-item label="登录名">
              <a-input v-model="queryParam.loginName" placeholder="" />
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
      <div slot="postNames" slot-scope="text,current">
        <template v-for="(name, index) in current.postNames">
          <a-tag :key="index" :title="name">
            {{ name }}
          </a-tag>
        </template>
      </div>
      <div slot="roleNames" slot-scope="text,current">
        <template v-for="(name, index) in current.roleNames">
          <a-tag :key="index" :title="name">
            {{ name }}
          </a-tag>
        </template>
      </div>
    </s-table>

  </div>
</template>

<script>
import STable from '@/components/Table'
import { loadUser } from '@/api/user'
import { columns } from './UserList.js'

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
      selectedRowKeys: [],
      selectedRows: [],
      queryParam: {},
      record: {},
      loadData: (parameter) => {
        return loadUser(Object.assign(parameter, this.queryParam)).then((res) => {
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
    loadUser: function (deptIds = [], bool = true) {
      this.queryParam.deptIds = deptIds
      this.$refs.table.refresh(bool)
    }
  }
}
</script>

<style scoped></style>
