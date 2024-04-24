<template>
  <div>
    <a-layout>
      <a-layout-sider width="280">
        <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 210px)`,overflow:'hidden'}" title="部门">
          <c-tree
            :loading="loading"
            ref="deptTree"
            :data="loadDeptTree"
            :load-done="loadUser"
            :check="(keys) => loadUser([],keys)"
            @select="selectNode"
            :show-operate-btn="false"
          >
          </c-tree>

        </a-card>
      </a-layout-sider>
      <a-layout>
        <a-layout-content>
          <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 210px)`,overflow:'hidden'}" title="用户列表">
            <UserList @change="onSelectChange" ref="userList"/>
          </a-card>
        </a-layout-content>
      </a-layout>
    </a-layout>

  </div>
</template>

<script>
import CTree from '@/components/Custom/Tree/'
import UserList from './UserList.vue'
import { loadDept } from '@/api/user'

export default {
  name: 'SelectIndex',
  components: { CTree, UserList },
  data() {
    return {
      loading: false,
      windowHeight: 0,
      deptIds: undefined,
      deptTreeData: [],
      selectedRowKeys: [],
      selectedRows: [],
      loadDeptTree: () => {
        this.loading = true
        return loadDept().then(res => {
          const { code, data } = res

          if (code !== 200) {
            return res
          }
          this.deptTreeData = [data]
          if (data.children && data.children.length > 0) {
            res.data = data.children
          }
          return res
        }).finally(() => {
          this.loading = false
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
    resize: function () {
      this.windowHeight = document.body.clientHeight
    },
    loadUser: function(treeData, checkedKeys) {
      this.$refs.userList.loadUser(checkedKeys)
    },
    selectNode: function(key, e) {
      this.$refs.userList.loadUser(key)
    },
    onSelectChange: function (keys, rows) {
      this.selectedRowKeys = keys
      this.selectedRows = rows
    },
    onOK: function () {
      if (this.selectedRowKeys.length <= 0) {
        this.$message.warn('请选择数据')
        return
      }
      this.$emit('select', this.selectedRowKeys, this.selectedRows)
      return true
    }
  }
}
</script>

<style scoped>
.ant-layout-content {
  margin-left: 20px
}
</style>
