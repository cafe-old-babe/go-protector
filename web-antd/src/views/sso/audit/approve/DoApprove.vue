<template>
  <div>
    <a-drawer
      title="处理审批"
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
            label="工单号"
            prop="workNum"
          >
            <a-input v-model="localRecord.workNum" :disabled="true"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="审批类型"
            prop="approveTypeName"
          >
            <a-input v-model="localRecord.approveTypeName" :disabled="true"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="申请人"
            prop="applicantUsername"
          >
            <a-input v-model="localRecord.applicantUsername" :disabled="true"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="申请原因"
            prop="applicantContent"
          >
            <a-input v-model="localRecord.applicantContent" type="textarea" :disabled="true"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="审批"
            prop="approveStatus"
          >
            <a-radio-group
              v-model="localRecord.approveStatus"
            >
              <a-radio-button :value="1">
                通过
              </a-radio-button>
              <a-radio-button :value="2">
                拒绝
              </a-radio-button>
            </a-radio-group>
          </a-form-model-item>

          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="审批意见"
            prop="approveContent"
          >
            <a-input v-model="localRecord.approveContent" type="textarea" />
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
          处理
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
        approveStatus: [
          { required: true, message: '请输入选择审批结果' }
        ],
        approveContent: [
          { required: false, message: '不通过需要填写拒绝原因' }
        ]

      }
    }
  },
  watch: {
    record(val) {
      this.localRecord = Object.assign({}, val)
      console.log(this.localRecord)
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
        if (this.localRecord.approveStatus === 2 && !this.localRecord.approveContent) {
          this.$message.warning('不通过需要填写审批意见')
          this.loading = false
          return false
        }
        request.post('/approve-record/doApprove', Object.assign({}, {
          'id': this.localRecord.id,
          'approveStatus': this.localRecord.approveStatus,
          'approveContent': this.localRecord.approveContent
        })).then(res => {
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
