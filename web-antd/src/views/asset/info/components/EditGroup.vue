<template>
  <div>
    <a-drawer
      :title="localRecord.id?'编辑资源组':'新增资源组'"
      :width="500"
      :visible="visible"
      :body-style="{ paddingBottom: '80px' }"
      @close="onClose"
    >
      <a-spin :spinning="loading">
        <a-form-model
          ref="form"
          :model="localRecord"
          :rules="rules"
          :label-col="{ span: 6 }"
          :wrapper-col="{ span: 14 }"
          layout="horizontal">
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="上级资源组"
            prop="pid"
          >
            <a-tree-select
              v-model="localRecord.pid"
              style="width: 100%"
              :dropdown-style="{ maxHeight: '200px', overflow: 'auto' }"
              :tree-data="treeData"
              :replace-fields="replaceFields"
              placeholder="请选择上级资源组"
              :defaultValue="defaultValue"
              @select="selectTreeNode"
              tree-default-expand-all
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="资源组名称"
            prop="name"
          >
            <a-input
              v-model="localRecord.name"
              placeholder="请输入资源组名称"
            />
          </a-form-model-item>
        </a-form-model>
      </a-spin>
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

export default {
  props: {
    visible: {
      type: Boolean,
      required: false
    },
    record: {
      type: Object,
      default: () => null
    },
    treeData: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      loading: false,
      localRecord: {},
      rules: {
        pid: [
          { required: true, message: '请选择上级资源组' }
        ],
        name: [
          { required: true, message: '请输入资源组名称' }
        ]
      },
      replaceFields: {
        key: 'id',
        value: 'id',
        title: 'name'
      }
    }
  },
  watch: {
    record(val) {
      this.localRecord = val
      if (this.treeData) {
        this.treeData[0].children = this.disabledLocal(this.treeData[0].children, this.localRecord.id)
      }
      this.loading = false
    }
  },
  computed: {
    formItemLayout() {
      const { formLayout } = this
      return formLayout === 'horizontal'
        ? {
          labelCol: { span: 2 },
          wrapperCol: { span: 14 }
        }
        : {}
    },
    defaultValue() {
      return this.localRecord.pid
    }
  },
  methods: {
    onClose() {
      this.$emit('close')
    },
    handleSave() {
      this.loading = true
      this.$refs.form.validate(valid => {
        if (!valid) {
          this.loading = false
          return false
        }

        request.post('/asset-group/save', this.localRecord).then(res => {
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
      })
    },
    selectTreeNode: function (value) {
      this.localRecord.pid = value
    },
    disabledLocal: function (data, id) {
      if (data && id) {
        data.forEach((item) => {
          item.disabled = item.id === id
          if (item.children) {
            item.children = this.disabledLocal(item.children, id)
          }
        })
      }

      return data
    }
  }
}
</script>
