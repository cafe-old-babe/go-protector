<template>
  <div>
    <a-drawer
      title="菜单授权"
      :width="600"
      :visible="visible"
      :body-style="{ paddingBottom: '80px' }"
      @close="onClose"
    >

      <a-table
        ref="PermissionTable"
        rowKey="id"
        :loading="loading"
        :columns="columns"
        :data-source="menuTreeData"
        :rowSelection="rowSelection"
        :pagination="false"
        :defaultExpandAllRows="true"
      />

      <div
        :style="{
          position: 'absolute',
          right: 0,
          bottom: 0,
          width: '100%',
          borderTop: '1px solid #e9e9e9',
          padding: '10px 16px',
          background: '#fff',
          textAlign: 'right',
          zIndex: 1,
        }"
      >
        <a-button :style="{ marginRight: '8px' }" @click="onClose">
          取消
        </a-button>
        <a-button type="primary" :loading="loading" @click="handleSave">
          保存
        </a-button>
      </div>
    </a-drawer>
  </div>
</template>
<script>
import request from '@/utils/request'
import { MenuColumns } from './column'
export default {
  props: {
    visible: {
      type: Boolean,
      required: false
    },
    roleId: {
      type: Number || String,
      required: true
    }
  },
  data() {
    return {
      loading: true,
      menuTreeData: [],
      selectedRowKeys: [],
      columns: MenuColumns,
      loadData: (roleId) => {
        this.loading = true
        console.log('loadData', this.visible)
        request.post('/menu/list').then((res) => {
          const { code, data, message } = res
          if (code === 200) {
            this.menuTreeData = data.children
          } else {
            this.$message.error(message)
          }
        }).finally(() => {
          this.loadPermission(roleId)
        })
      }
    }
  },
  watch: {
    roleId(val) {
      if (!this.visible) {
        return
      }
      this.loadData(val)
      this.loading = false
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
  methods: {
    onSelectChange(selectedRowKeys) {
      this.selectedRowKeys = selectedRowKeys
    },
    onClose() {
      this.$emit('close')
    },
    handleSave() {
      this.loading = true
      request.post(`/role/savePermission/${this.roleId}`,
        { ids: this.selectedRowKeys }).then(res => {
        const { code, message } = res
        if (code === 200) {
          this.$emit('ok')
          this.$message.success(message)
          return
        }
        this.$message.error(message)
      }).finally(() => {
        this.loading = false
      })
    },
    loadPermission(roleId) {
      request.post(`/role/getPermission/${roleId}`).then(res => {
        const { code, message, data } = res
        data.push(1)
        if (code === 200) {
          this.selectedRowKeys = data
        } else {
          this.$message.error(message)
        }
      }).finally(() => {
        this.loading = false
        this.$forceUpdate()
      })
    }
  }

}
</script>
