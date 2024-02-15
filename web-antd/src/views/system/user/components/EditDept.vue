<template>
  <div>
    <a-drawer
      :title="localRecord.id?'编辑部门':'新增部门'"
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
            label="上级部门"
            prop="pid"
          >
            <a-tree-select
              v-model="localRecord.pid"
              style="width: 100%"
              :dropdown-style="{ maxHeight: '200px', overflow: 'auto' }"
              :tree-data="deptTreeData"
              :replace-fields="replaceFields"
              placeholder="请选择上级部门"
              :defaultValue="defaultValue"
              @select="selectTreeNode"
              tree-default-expand-all
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="部门名称"
            prop="deptName"
          >
            <a-input
              v-model="localRecord.deptName"
              placeholder="请输入部门名称"
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
    deptTreeData: {
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
          { required: true, message: '请选择上级部门' }
        ],
        deptName: [
          { required: true, message: '请输入部门名称' }
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
      return this.localRecord.parentId
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

        request.post('/user/dept/save', this.localRecord).then(res => {
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
      console.log('value', value)
      this.localRecord.parentId = value
    }
  }
}
</script>
