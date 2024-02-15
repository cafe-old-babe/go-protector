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
            @addTreeNode="addTreeNode"
            @updateTreeNode="updateTreeNode"
            @deleteTreeNode="deleteTreeNode"
          >
          </c-tree>

        </a-card>
      </a-layout-sider>
      <a-layout>
        <a-layout-content>
          <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 210px)`,overflow:'hidden'}" title="用户列表">
            <UserList ref="userList"/>
          </a-card>
        </a-layout-content>
      </a-layout>
    </a-layout>

    <EditDept
      :visible="editVisible"
      :record="record"
      @close="editClose"
      @ok="editOk"
      :dept-tree-data="deptTreeData"/>
  </div>
</template>

<script>
import CTree from '@/components/Custom/Tree/'
import UserList from './components/UserList.vue'
import { loadDept } from '@/api/user'
import EditDept from './components/EditDept.vue'
import request from '@/utils/request'

export default {
  name: 'User',
  components: { EditDept, CTree, UserList },
  data() {
    return {
      loading: false,
      windowHeight: 0,
      deptIds: undefined,
      deptTreeData: [],
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
      },
      operateTreeNode: false,
      editVisible: false,
      record: {}
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
    deleteTreeNode: function(node) {
      this.operateTreeNode = true
      this.$confirm({
        title: '确认删除',
        content: '是否删除选择数据?',
        okText: '确认',
        okType: 'danger',
        cancelText: '取消',
        confirmLoading: this.loading,
        onOk: () => {
          request.post('/user/dept/delete', { ids: [node.id] }).then(res => {
            console.log(res)
            if (res.code === 200) {
              this.$refs.deptTree.loadData()
              this.$message.info(res.message ?? '删除成功')
              return
            }
            this.$message.warn(res.message)
          }).finally(() => {
            this.operateTreeNode = false
          })
        }
      })
    },
    addTreeNode: function(node) {
      this.operateTreeNode = true
      this.record = {
        pid: node.id
      }
      this.editVisible = true
    },
    updateTreeNode: function(node) {
      console.log(node)
      this.operateTreeNode = true
      this.record = {
        pid: node.pid,
        id: node.id,
        deptName: node.name
      }
      this.editVisible = true
    },
    selectNode: function(key, e) {
      if (this.operateTreeNode) {
        return
      }
      this.$refs.userList.loadUser(key)
    },
    editOk: function () {
      this.editClose()
      this.$refs.deptTree.loadData()
    },
    editClose: function () {
      this.record = {}
      this.operateTreeNode = false
      this.editVisible = false
    }
  }
}
</script>

<style scoped>
.ant-layout-content {
  margin-left: 20px
}
</style>
