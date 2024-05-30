<template>
  <div>
    <a-layout>
      <a-layout-sider width="280">
        <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 150px)`,overflow:'auto'}" title="资源组">
          <c-tree
            :loading="loading"
            ref="groupTree"
            :data="loadGroupTree"
            :load-done="loadUser"
            :check="(keys) => loadUser([],keys)"
            @select="selectNode"
            :show-operate-btn="true"
            @addTreeNode="addTreeNode"
            @updateTreeNode="updateTreeNode"
            @deleteTreeNode="deleteTreeNode"
          >
          </c-tree>

        </a-card>
      </a-layout-sider>
      <a-layout>
        <a-layout-content>
          <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 150px)`,overflow:'auto'}" title="资源列表">
            <InfoTable ref="infoTable"/>
          </a-card>
        </a-layout-content>
      </a-layout>
    </a-layout>

    <EditGroup
      :visible="editVisible"
      :record="record"
      @close="editClose"
      @ok="editOk"
      :tree-data="treeData"/>
  </div>
</template>

<script>
import CTree from '@/components/Custom/Tree/'
import InfoTable from './components/InfoTable.vue'
import { loadGroupTree } from '@/api/asset'
import EditGroup from './components/EditGroup.vue'
import request from '@/utils/request'

export default {
  name: 'AssetInfo',
  components: { EditGroup, CTree, InfoTable },
  data() {
    return {
      loading: false,
      windowHeight: 0,
      deptIds: undefined,
      treeData: [],
      loadGroupTree: () => {
        this.loading = true
        return loadGroupTree().then(res => {
          const { code, data } = res

          if (code !== 200) {
            return res
          }
          this.treeData = [data]
          if (data.children && data.children.length > 0) {
            res.data = data.children
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
  destroyed () {
    window.removeEventListener('resize', this.resize, false)
  },
  methods: {
    resize: function () {
      this.windowHeight = document.body.clientHeight
    },
    loadUser: function(treeData, checkedKeys) {
      this.$refs.infoTable.loadInfo(checkedKeys)
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
          request.post('/asset-group/delete', { ids: [node.id] }).then(res => {
            if (res.code === 200) {
              this.$refs.groupTree.loadData()
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
      this.operateTreeNode = true
      this.record = {
        pid: node.pid,
        id: node.id,
        name: node.name
      }
      this.editVisible = true
    },
    selectNode: function(key, e) {
      if (this.operateTreeNode) {
        return
      }
      this.$refs.infoTable.loadInfo(key)
    },
    editOk: function () {
      this.editClose()
      this.$refs.groupTree.loadData()
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
