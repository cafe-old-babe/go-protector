<template>
  <div>
    <a-layout>
      <a-layout-sider width="280">
        <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 210px)`,overflow:'auto'}" title="资源组">
          <c-tree
            :loading="loading"
            ref="groupTree"
            :data="loadGroupTree"
            :load-done="loadData"
            :check="(keys) => loadData([],keys)"
            @select="selectNode"
            :show-operate-btn="false"
          >
          </c-tree>

        </a-card>
      </a-layout-sider>
      <a-layout>
        <a-layout-content>
          <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 210px)`,overflow:'auto'}" title="资源列表">
            <InfoTable ref="infoTable" @change="onSelectChange"/>
          </a-card>
        </a-layout-content>
      </a-layout>
    </a-layout>

  </div>
</template>

<script>
import CTree from '@/components/Custom/Tree/'
import InfoTable from './InfoTable.vue'
import { loadGroupTree } from '@/api/asset'

export default {
  name: 'SelectIndex',
  components: { CTree, InfoTable },
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
    loadData: function(treeData, checkedKeys) {
      this.$refs.infoTable.loadInfo(checkedKeys)
    },
    selectNode: function(key, e) {
      if (this.operateTreeNode) {
        return
      }
      this.$refs.infoTable.loadInfo(key)
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
