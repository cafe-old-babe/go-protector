<template>
  <div>
    <a-drawer
      :title="localRecord.id?'编辑岗位':'新增岗位'"
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
            prop="deptIds"
          >
            <!--            4-21	【实战】岗位页面实战开发-掌握TreeSelect组件、VUE-组件之间同步操作-->
            <a-tree-select
              v-model="localRecord.deptIds"
              style="width: 100%"
              :dropdown-style="{ maxHeight: '200px', overflow: 'auto' }"
              :tree-data="deptTreeData"
              :replace-fields="replaceFields"
              placeholder="请选择上级部门"
              :defaultValue="localRecord.deptIds"
              multiple
              tree-default-expand-all
              :show-checked-strategy="showStrategy"
              tree-checkable
              @change="selectTreeNode"
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="岗位名称"
            prop="name"
          >
            <a-input
              v-model.trim="localRecord.name"
              placeholder="请输入岗位名称"
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="岗位代码"
            prop="code"
          >
            <a-input
              v-model.trim="localRecord.code"
              placeholder="请输入岗位代码"
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
import { loadDept } from '@/api/user'
import { TreeSelect } from 'ant-design-vue'

export default {
  props: {
    visible: {
      type: Boolean,
      required: false
    },
    record: {
      type: Object,
      default: () => null
    }
  },
  data() {
    return {
      loading: false,
      localRecord: {},
      showStrategy: TreeSelect.SHOW_ALL,
      deptTreeData: [],
      rules: {
        deptIds: [
          { required: true, message: '请选择上级部门', trigger: 'change' }
        ],
        name: [
          { required: true, message: '请输入岗位名称' }
        ],
        code: [
          { required: true, message: '请输入岗位代码' }
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
      this.localRecord = Object.assign({}, val)
      if (this.localRecord.deptIds && this.localRecord.deptIds.length > 0) {
        this.localRecord.deptIds = this.localRecord.deptIds?.split(',').map(Number)
      } else {
        this.localRecord.deptIds = []
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
    }
  },
  mounted() {
    this.loading = true
    loadDept().then(res => {
      const { code, data } = res

      if (code !== 200) {
        return res
      }
      if (data.children && data.children.length > 0) {
        this.deptTreeData = data.children
      }
    }).finally(() => {
      this.loading = false
    })
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
        request.post('/post/save', this.localRecord).then(res => {
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
      this.localRecord.deptIds = value
      this.$forceUpdate()
    }
  }
}
</script>
