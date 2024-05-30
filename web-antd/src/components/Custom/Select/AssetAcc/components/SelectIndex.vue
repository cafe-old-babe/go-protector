<template>
  <InfoTable ref="infoTable" :user-id="userId" @change="onSelectChange"/>
</template>

<script>
import CTree from '@/components/Custom/Tree/'
import InfoTable from './InfoTable.vue'

export default {
  name: 'SelectIndex',
  components: { CTree, InfoTable },
  props: {
    userId: {
      type: [Number],
      required: false,
      default: 0
    }
  },
  data() {
    return {
      loading: false,
      windowHeight: 0,
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
