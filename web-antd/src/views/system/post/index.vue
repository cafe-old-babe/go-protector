
<template>
  <a-layout>
    <a-layout-sider width="280">
      <a-card
        :bordered="false"
        :style="{height:`calc(${windowHeight}px - 210px)`,overflow:'hidden'}"
        title="部门">
        <c-tree
          :loading="loading"
          ref="deptTree"
          :data="loadDeptTree"
          :load-done="loadPost"
          :check="(keys) => loadPost([],keys)"
          @updateTreeNode="updateTreeNode"
          @addTreeNode="addTreeNode"
          @deleteTreeNode="deleteTreeNode"
          @select="selectNode"
        >
        </c-tree>

      </a-card>
    </a-layout-sider>
    <a-layout>
      <a-layout-content>
        <a-card
          :bordered="false"
          :style="{height:`calc(${windowHeight}px - 210px)`,overflow:'hidden'}"
          :title="`[${currentDeptName??'全部'}]岗位列表`">
          <PostTable ref="postTable"/>
        </a-card>
      </a-layout-content>
    </a-layout>
    <EditDept
      :visible="editVisible"
      :record="record"
      @close="editClose"
      @ok="editOk"
      :dept-tree-data="deptTreeData"/>
  </a-layout>
</template>
<script>
import CTree from '@/components/Custom/Tree/index.vue'
import { loadDept } from '@/api/user'
import request from '@/utils/request'
import PostTable from './components/PostTable.vue'
import EditDept from './components/EditDept.vue'
import EditPost from './components/EditPost.vue'

export default {
  name: 'Index',
  components: { EditDept, CTree, PostTable, EditPost },
  data() {
    return {
      loading: false,
      windowHeight: 0,
      deptIds: undefined,
      deptTreeData: [],
      currentDeptName: null,
      loadDeptTree: () => {
        this.loading = true
        return loadDept().then(res => {
          const { code, data } = res

          if (code !== 200) {
            return res
          }
          if (data.children && data.children.length > 0) {
            res.data = data.children
            this.deptTreeData = res.data
          } else {
            res.data = []
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
  destroyed() {
    window.removeEventListener('resize', this.resize, false)
  },
  methods: {
    resize: function () {
      this.windowHeight = document.body.clientHeight
    },
    loadPost: function (treeData, checkedKeys) {
      this.$refs.postTable.loadPost(checkedKeys)
    },
    deleteTreeNode: function (node) {
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
    addTreeNode: function (node) {
      this.operateTreeNode = true
      this.record = {
        pid: node.id
      }
      this.editVisible = true
    },
    updateTreeNode: function (node) {
      this.operateTreeNode = true
      this.record = {
        pid: node.pid,
        id: node.id,
        deptName: node.name
      }
      this.editVisible = true
    },
    selectNode: function (key, e) {
      console.log(e)
      if (!this.operateTreeNode) {
        this.currentDeptName = e.selectedNodes.length > 0 ? e.selectedNodes[0].data.props.name : null
        this.$refs.postTable.loadPost(key)
      }
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
<style scoped lang="less">
.ant-layout-content {
  margin-left: 20px
}
.ant-card {
  border-radius: 0
}
</style>
