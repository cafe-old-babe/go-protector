<template>
  <div>

    <a-card :bordered="false" :style="{height:`calc(${windowHeight}px - 150px)`,overflow:'auto'}">
      <a-skeleton :loading="skeleton">
        <a-spin :spinning="loading">
          <a-tabs type="card" :default-active-key="activeKey" @change="(key) => activeKey = key">
            <a-tab-pane key="global" :tab="localRecord.global.name??'全局认证配置' ">
              <a-alert :message="`${localRecord.global.name}只控制[${localRecord.email.name}]与[${localRecord.otp.name}]`" banner />

              <a-form-model
                ref="global"
                :model="localRecord.global"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 6 }"
                :rules="rules"
                layout="horizontal">

                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="是否启用"
                  prop="enable"
                >
                  <a-radio-group
                    :default-value="localRecord.global.enable??0"
                    button-style="solid"
                    v-model="localRecord.global.enable"
                  >
                    <a-radio-button value="0">
                      停用
                    </a-radio-button>
                    <a-radio-button value="1">
                      启用
                    </a-radio-button>
                  </a-radio-group>
                </a-form-model-item>
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="通过方式"
                  prop="mode"
                >
                  <a-radio-group
                    button-style="solid"
                    v-model="localRecord.global.mode">
                    <a-radio-button value="0">
                      通过全部
                    </a-radio-button>
                    <a-radio-button value="1">
                      通过一种即可
                    </a-radio-button>
                  </a-radio-group>
                </a-form-model-item>
              </a-form-model>
            </a-tab-pane>
            <a-tab-pane key="email" :tab="localRecord.email.name??'登录邮箱'">
              <a-form-model
                ref="email"
                :model="localRecord.email"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 6 }"
                :rules="rules"
                layout="horizontal">
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="是否启用"
                  prop="enable"
                >
                  <a-radio-group
                    :default-value="localRecord.email.enable??0"
                    button-style="solid"
                    v-model="localRecord.email.enable"
                  >
                    <a-radio-button value="0">
                      停用
                    </a-radio-button>
                    <a-radio-button value="1">
                      启用
                    </a-radio-button>
                  </a-radio-group>
                </a-form-model-item>
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="过期时间(分钟)"
                  prop="expireTime"
                >
                  <a-input-number
                    v-model="localRecord.email.expireTime"
                    :min="2"
                    :max="10"
                    step="1"
                    :precision="0" />
                </a-form-model-item>
              </a-form-model>
            </a-tab-pane>
            <a-tab-pane key="otp" :tab="localRecord.otp.name??'动态密码认证策略'">
              <a-form-model
                ref="otp"
                :model="localRecord.otp"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 6 }"
                :rules="rules"
                layout="horizontal">
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="是否启用"
                  prop="enable"
                >
                  <a-radio-group
                    :default-value="localRecord.otp.enable??0"
                    button-style="solid"
                    v-model="localRecord.otp.enable"
                  >
                    <a-radio-button value="0">
                      停用
                    </a-radio-button>
                    <a-radio-button value="1">
                      启用
                    </a-radio-button>
                  </a-radio-group>
                </a-form-model-item>
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="签发机构名称"
                  prop="issuer"
                >
                  <a-input
                    v-model.trim="localRecord.otp.issuer"
                    placeholder="请输入签发机构名称"
                  />
                </a-form-model-item>
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="动态密码有效时间(秒)"
                  prop="period"
                >
                  <a-input-number
                    v-model="localRecord.otp.period"
                    :min="30"
                    :max="60"
                    step="10"
                    :precision="0"/>
                </a-form-model-item>
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="秘钥长度"
                  prop="secretSize"
                >
                  <a-input-number
                    v-model="localRecord.otp.secretSize"
                    :min="12"
                    :max="36"
                    step="12"
                    :precision="0"/>
                </a-form-model-item>
              </a-form-model>
            </a-tab-pane>
            <a-tab-pane key="share" :tab="localRecord.share.name??'共享登录策略'">
              <a-alert message="登录间隔时间设置为:0,表示不限制" banner />
              <a-form-model
                ref="share"
                :model="localRecord.share"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 6 }"
                :rules="rules"
                layout="horizontal">
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="是否启用"
                  prop="enable"
                >
                  <a-radio-group
                    :default-value="localRecord.share.enable??0"
                    button-style="solid"
                    v-model="localRecord.share.enable"
                  >
                    <a-radio-button value="0">
                      停用
                    </a-radio-button>
                    <a-radio-button value="1">
                      启用
                    </a-radio-button>
                  </a-radio-group>
                </a-form-model-item>
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="登录间隔时间(分)"
                  prop="loginInterval"
                >
                  <a-input-number
                    v-model="localRecord.share.loginInterval"
                    :min="0"
                    :max="10"
                    step="1"
                    :precision="0"/>
                </a-form-model-item>
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="单用户在线"
                  prop="singleOnline"
                >
                  <a-radio-group
                    :default-value="localRecord.share.singleOnline??0"
                    button-style="solid"
                    v-model="localRecord.share.singleOnline"
                  >
                    <a-radio-button :value="0">
                      停用
                    </a-radio-button>
                    <a-radio-button :value="1">
                      启用
                    </a-radio-button>
                  </a-radio-group>
                </a-form-model-item>
                <a-form-model-item
                  v-if="localRecord.share.singleOnline===1"
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="异地多人登录时"
                  prop="singleOnlineOperate"
                >
                  <a-radio-group
                    :default-value="localRecord.share.singleOnlineOperate??0"
                    button-style="solid"
                    v-model="localRecord.share.singleOnlineOperate"
                  >
                    <a-radio-button :value="0">
                      禁止当前用户登录
                    </a-radio-button>
                    <a-radio-button :value="1">
                      踢掉在线用户
                    </a-radio-button>
                  </a-radio-group>
                </a-form-model-item>
              </a-form-model>
            </a-tab-pane>
            <a-tab-pane key="intruder" :tab="localRecord.intruder.name??'防爆破登录策略'">
              <a-alert message="当用户输入密码错误到达阈值时,将锁定用户" banner />
              <a-form-model
                ref="intruder"
                :model="localRecord.intruder"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 6 }"
                :rules="rules"
                layout="horizontal">
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="是否启用"
                  prop="enable"
                >
                  <a-radio-group
                    :default-value="localRecord.intruder.enable??0"
                    button-style="solid"
                    v-model="localRecord.intruder.enable"
                  >
                    <a-radio-button value="0">
                      停用
                    </a-radio-button>
                    <a-radio-button value="1">
                      启用
                    </a-radio-button>
                  </a-radio-group>
                </a-form-model-item>
                <a-form-model-item
                  :label-col="formItemLayout.labelCol"
                  :wrapper-col="formItemLayout.wrapperCol"
                  label="密码错误阈值"
                  prop="failThreshold"
                >
                  <a-input-number
                    v-model="localRecord.intruder.failThreshold"
                    :min="5"
                    :max="10"
                    step="1"
                    :precision="0"/>
                </a-form-model-item>
              </a-form-model>
            </a-tab-pane>
          </a-tabs>

        </a-spin>
      </a-skeleton>
    </a-card>
    <div
      :style="{
        position: 'absolute',
        right: 0,
        bottom: 0,
        width: '100%',
        borderTop: '1px solid #e9e9e9',
        padding: '10px 16px',
        background: '#fff',
        textAlign: 'center',
        zIndex: 1,
      }"
    >
      <a-button type="primary" :loading="loading" @click="handleSave">
        保存
      </a-button>
    </div>
  </div>
</template>

<script>
import request from '@/utils/request'
// 5-5	【实战】完成策略页面实战-掌握antd Tabs、Skeleton组件；掌握使用策略模式完成校验和提交
export default {
  data() {
    return {
      loading: false,
      skeleton: true,
      windowHeight: 0,
      activeKey: 'global',
      localRecord: {
        global: { name: '1', enable: 0, mode: 0 },
        email: { name: '1', enable: 0, expireTime: 0 },
        otp: { name: '1', enable: 0, issuer: '', period: 0, secretSize: 0 },
        share: { name: '1', enable: 0, loginInterval: 0 }
      },
      rules: {
        enable: [{ required: true, message: '请选择是否启用', trigger: 'blur' }],
        mode: [{ required: true, message: '请选择校验方式', trigger: 'blur' }],
        expireTime: [{ required: true, message: '请输入邮箱校验过期时间', trigger: 'blur' }],
        period: [{ required: true, message: '请输入动态密码有效时间', trigger: 'blur' }],
        secretSize: [{ required: true, message: '请输入秘钥长度', trigger: 'blur' }],
        issuer: [{ required: true, message: '请输入签发机构名称', trigger: 'blur' }],
        loginInterval: [{ required: true, message: '请输入登录间隔限制', trigger: 'blur' }],
        singleOnlineOperate: [{ required: true, message: '请选择多人登录时操作', trigger: 'blur' }],
        singleOnline: [{ required: true, message: '请选择单用户在线是否启用', trigger: 'blur' }],
        failThreshold: [{ required: true, message: '请输入密码错误阈值', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.loadData()
  },
  mounted() {
    this.resize()
    window.addEventListener('resize', this.resize, false)
  },
  destroyed () {
    window.removeEventListener('resize', this.resize, false)
  },
  computed: {
    disabled() {
      return this.localRecord.global.enable === 0
    },
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
    resize () {
      this.windowHeight = document.body.clientHeight
    },
    loadData() {
      this.skeleton = true
      request.post('/sys-login-policy/info').then(res => {
        const { code, data, message } = res
        if (code !== 200) {
          this.$message.error(message)
        }
        this.localRecord = data
      }).finally(() => {
        this.skeleton = false
      })
    },
    handleSave() {
      // 获取 form表单中指定数据
      // this.$refs.form.validateField('globalEnable')
       this.loading = true

      this.$refs[this.activeKey].validate(valid => {
        if (!valid) {
          this.loading = false
          return false
        }
        const param = {}
        param[this.activeKey] = this.localRecord[this.activeKey]
         request.post('/sys-login-policy/save', param).then(res => {
          const { code, message } = res
          if (code === 200) {
            this.$emit('ok')
            this.$message.success(message)
            this.loadData()
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

<style scoped lang="less">

</style>
