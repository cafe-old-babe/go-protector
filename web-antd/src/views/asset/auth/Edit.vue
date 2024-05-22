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
            label="所属资源"
            prop="assetId"
          >
            <select-asset
              ref="selectAsset"
              v-model="localRecord.assetId"
              :show-label="localRecord.assetInfoName"
              :show-operate="!localRecord.id"
              :placeholder="rules.assetId.message"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="从帐号类型"
            prop="accountType"
          >
            <a-radio-group
              v-model="localRecord.accountType"
              :default-value="localRecord.accountType"
            >
              <a-radio-button value="1">
                管理从帐号
              </a-radio-button>
              <a-radio-button value="2">
                普通从帐号
              </a-radio-button>
            </a-radio-group>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="从帐号"
            prop="account"
          >
            <a-input v-model="localRecord.account" :placeholder="rules.account.message"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="从帐号"
            prop="password"
          >
            <a-input-password v-model="localRecord.password" :placeholder="rules.password.message"/>
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
import SelectAsset from '@/components/Custom/Select/Asset'

export default {
  components: { SelectAsset },
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
      rules: {
        assetId: [
          { required: true, message: '请输入从帐号名称' }
        ],
        account: [
          { required: true, message: '请输入从帐号' }
        ],
        accountType: [
          { required: true, message: '请选择从帐号类型' }
        ],
        password: [
          { required: true, message: '请输入从帐号名称' }
        ]
      }
    }
  },
  watch: {
    record(val) {
      this.localRecord = Object.assign({}, val)
      this.rules.password[0].required = !this.localRecord.id
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
      this.$refs.selectAsset.remove()
      this.$emit('close')
    },
    handleSave() {
      this.loading = true
      this.$refs.form.validate(valid => {
        if (!valid) {
          this.loading = false
          return false
        }
        if (!this.localRecord.accountStatus) {
          this.localRecord.accountStatus = '0'
        }
        request.post('/asset-account/save', this.localRecord).then(res => {
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
    }
  }
}
</script>
