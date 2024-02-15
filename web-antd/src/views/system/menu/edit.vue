<template>
  <div>
    <a-drawer
      :title="localRecord.id?'编辑':'新增'"
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
          <a-input v-model="localRecord.id" v-show="false"/>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="父级"
            prop="pid"
          >
            <a-tree-select
              v-model="localRecord.pid"
              style="width: 100%"
              :dropdown-style="{ maxHeight: '200px', overflow: 'auto' }"
              :tree-data="menuTreeData"
              :replace-fields="replaceFields"
              placeholder="请选择父级"
              :defaultValue="localRecord.pid"
              tree-default-expand-all
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="名称"
            prop="name"
          >
            <a-input v-model="localRecord.name" placeholder="请输入名称"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="权限标识"
            prop="permission"
          >
            <a-input
              v-model="localRecord.permission"
              placeholder="权限标识"
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="类型"
            prop="menuType"
          >
            <a-radio-group
              v-model="localRecord.menuType"
              :default-value="localRecord.menuType"
              @select="selectMenuType"
            >
              <a-radio-button :value="0">
                目录
              </a-radio-button>
              <a-radio-button :value="1">
                菜单
              </a-radio-button>
              <a-radio-button :value="2">
                按钮
              </a-radio-button>
            </a-radio-group>
          </a-form-model-item>
          <a-form-model-item
            label="显示状态"
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            prop="hidden"
            v-if="localRecord.menuType!==2"
          >
            <a-radio-group v-model="localRecord.hidden" :default-value="localRecord.hidden??0" button-style="solid">
              <a-radio :value="0">
                显示
              </a-radio>
              <a-radio :value="1">
                隐藏
              </a-radio>
            </a-radio-group>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="组件名称"
            prop="component"
            v-if="localRecord.menuType!==2"
          >
            <a-input
              v-model="localRecord.component"
              placeholder="组件名称"
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
    menuTreeData: {
      type: Array,
      requires: true,
      default: () => []
    },
    // 替换默认字段
    replaceFields: {
      type: Object,
      default: () => {
        return {
          value: 'id',
          key: 'id',
          title: 'name'
        }
      }
    }
  },
  data() {
    return {
      loading: false,
      localRecord: {},
      rules: {
        pid: [
          { required: true, message: '请选择父级' }
        ],
        permission: [
          { required: true, message: '请输入权限标识' }
        ],
        name: [
          { required: true, message: '请输入名称' }
        ],
        menuType: [
          { required: true, message: '请选择类型' }
        ],
        hidden: [
          { required: true, message: '请选择显示状态' }
        ],
        component: [
          { required: true, message: '请选择显示状态' }
        ]
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

        request.post('/menu/save', this.localRecord).then(({ code, message }) => {
          if (code === 200) {
            this.$emit('ok')
            this.$message.success(message)
          } else {
          this.$message.error(message)
          }
        }).finally(() => {
          this.loading = false
        })
      })
    },
    selectMenuType: function (value) {
      console.log(value)
    }
  }
}
</script>
