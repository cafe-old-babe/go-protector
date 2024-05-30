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
            label="主帐号"
            prop="userId"
          >
            <select-user
              ref="selectUserRef"
              v-model="localRecord.userId"
              :show-label="localRecord.userAcc"
              :show-operate="!localRecord.id"
              @callback="selectUser"
              placeholder="请选择主帐号"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="从帐号"
            prop="assetAccId"
          >
            <select-asset-acc
              ref="selectAssetAcc"
              v-model="localRecord.assetAccId"
              :show-label="localRecord.accountLabel"
              :show-operate="!localRecord.id"
              :user-id="localRecord.userId"
              :placeholder="rules.assetAccId.message"
              @callback="(data) => bindAssetAcc(data[0])"
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="授权生效时间"
            prop="takeEffectDate"
          >
            <a-range-picker
              v-model="localRecord.takeEffectDate"
              :placeholder="rules.takeEffectDate.message"
              @change="changeTakeEffect"/>

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
import SelectAssetAcc from '@/components/Custom/Select/AssetAcc'
import SelectUser from '@/components/Custom/Select/User/index.vue'
import TagSelectOption from '@/components/TagSelect/TagSelectOption'
import moment from 'moment'
export default {
  components: { TagSelectOption, SelectUser, SelectAssetAcc },
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
      assetAccList: [],
      rules: {
        userId: [
          { required: true, message: '请选择主帐号' }
        ],
        assetAccId: [
          { required: true, message: '请选择从帐号' }
        ],
        takeEffectDate: [
          { required: false, message: '请选择生效时间' }
        ]
      }
    }
  },
  watch: {
    record(val) {
      this.localRecord = Object.assign({}, val)
      if (this.localRecord.id) {
        this.localRecord.accountLabel = `${this.localRecord.assetAcc}[${this.localRecord.assetName}(${this.localRecord.assetIp})]`
        if (this.localRecord.startDate && this.localRecord.endDate) {
          this.localRecord.takeEffectDate = [moment(this.localRecord.startDate), moment(this.localRecord.endDate)]
        }
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
  methods: {
    onClose() {
      this.$refs.selectAssetAcc.remove()
      this.localRecord = {}
      this.$emit('close')
    },
    handleSave() {
      this.loading = true
      this.$refs.form.validate(valid => {
        if (!valid) {
          this.loading = false
          return false
        }
        const saveObj = Object.assign({}, this.localRecord)
        delete saveObj.takeEffectDate
        request.post('/asset-auth/save', saveObj).then(res => {
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
    changeTakeEffect: function (date, dateString) {
      this.localRecord.startDate = dateString[0]
      this.localRecord.endDate = dateString[1]
    },
    bindAssetAcc: function (data) {
      this.localRecord.assetAccId = data?.id
      this.localRecord.assetAcc = data?.account
      this.localRecord.assetId = data?.assetBasic?.id
      this.localRecord.assetName = data?.assetBasic?.assetName
      this.localRecord.assetIp = data?.assetBasic?.ip
    },
    selectUser(data) {
      this.localRecord.userAcc = data[0].loginName
      this.localRecord.userId = data[0].id
      this.$refs.selectAssetAcc.remove()
      this.bindAssetAcc({})
    }

  }
}
</script>
