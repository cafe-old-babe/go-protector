<script>
import moment from 'moment'
import request from '@/utils/request'
export default {
 props: {
   visible: {
     type: Boolean,
     default: false
   },
   id: {
     type: Number,
     default: -1
   },
   status: {
     type: Number,
     default: -1
   }

 },
 data() {
   return {
     confirmLoading: false,
     param: {},
     form: this.$form.createForm(this)

   }
 },
  methods: {
   // https://1x.antdv.com/components/date-picker-cn/#components-date-picker-demo-disabled
    moment,
    range(start, end) {
      const result = []
      for (let i = start; i < end; i++) {
        result.push(i)
      }
      return result
    },
    disabledDate(current) {
      // Can not select days before today and today
      return current && current < moment().endOf('day')
    },

    disabledDateTime() {
      return {
        disabledHours: () => this.range(0, 24).splice(4, 20),
        disabledMinutes: () => this.range(30, 60),
        disabledSeconds: () => [55, 56]
      }
    },
    updateStatus: function () {
      this.form.validateFields((err, values) => {
        if (!err) {
          if (values.expirationAt) {
            values.expirationAt = values.expirationAt.format('YYYY-MM-DD HH:mm:ss')
          }
          this.loading = true
          request.post('/user/setStatus',
            Object.assign({ id: this.id, userStatus: (this.status ^ 1) },
              values)).then(res => {
            // console.log(res)
            const { code, message } = res
            if (code === 200) {
              this.$message.info(message)
              this.$emit('ok')
            } else {
              this.$message.warning(message)
            }
          }).finally(() => {
            this.loading = false
          })
        }
      })
    },
    handleCancel: function () {
      this.$emit('cancel')
    }
  }

}
</script>

<template>
  <div>
    <a-modal
      :visible="visible"
      :title="status===0?'锁定用户':'解锁用户'"
      ok-text="确认"
      cancel-text="取消"
      :confirm-loading="confirmLoading"
      @ok="updateStatus"
      @cancel="handleCancel">
      <a-form :form="form" :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" >
        <a-form-item label="锁定原因" v-if="status===0">
          <a-input
            v-decorator="['lockReason', { rules: [{ required: false, message: '请填写锁定原因' }] }]"
            placeholder="请填写锁定原因"
          />
        </a-form-item>
        <a-form-item label="用户有效期" v-else>
          <a-date-picker
            format="YYYY-MM-DD HH:mm:ss"
            :disabled-date="disabledDate"
            :disabled-time="disabledDateTime"
            :show-time="{ defaultValue: moment('00:00:00', 'HH:mm:ss') }"
            v-decorator="['expirationAt', { rules: [{ required: false, message: '请选择有效期' }] }]"
            placeholder="请选择有效期"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<style scoped lang="less">

</style>
