<template>
  <div class="main">
    <a-form id="formLogin" class="user-layout-login" ref="formLogin" :form="form" @submit="handleSubmit">
      <a-alert
        v-if="isLoginAlert"
        :type="loginAlertType"
        showIcon
        style="margin-bottom: 24px"
        :message="this.loginAlertMessage"
      />
      <a-tabs
        :activeKey="customActiveKey"
        :animated="false"
        :tabBarStyle="{ textAlign: 'center', borderBottom: 'unset' }"
        @change="handleTabClick"
      >
        <a-tab-pane key="firstTab" v-if="state.firstTab" :tab="$t('user.login.tab-login-credentials')">
          <a-form-item>
            <a-input
              size="large"
              type="text"
              :placeholder="$t('user.login.username.placeholder')"
              v-decorator="[
                'username',
                {
                  rules: [
                    { required: true, message: $t('user.userName.required') },
                    { validator: handleUsernameOrEmail },
                  ],
                  validateTrigger: 'change',
                },
              ]"
            >
              <a-icon slot="prefix" type="user" :style="{ color: 'rgba(0,0,0,.25)' }" />
            </a-input>
          </a-form-item>

          <a-form-item>
            <a-input-password
              size="large"
              :placeholder="$t('user.login.password.placeholder')"
              v-decorator="[
                'password',
                { rules: [{ required: true, message: $t('user.password.required') }], validateTrigger: 'blur' },
              ]"
            >
              <a-icon slot="prefix" type="lock" :style="{ color: 'rgba(0,0,0,.25)' }" />
            </a-input-password>
          </a-form-item>
          <a-form-item>
            <a-input
              size="large"
              placeholder="请输入验证码..."
              v-decorator="['code', { rules: [{ required: true, message: '请输入验证码!' }], validateTrigger: 'blur' }]"
            >
              <a-icon slot="prefix" type="safety" :style="{ color: 'rgba(0,0,0,.25)' }" />
              <span slot="suffix" @click="refreshCode" title="点击刷新验证码">
                <img style="width: 100px" :src="this.b64s" alt="点击刷新验证码" />
              </span>
            </a-input>
          </a-form-item>
        </a-tab-pane>
        <a-tab-pane key="emailTab" v-if="state.emailTab" tab="邮箱认证">
          <!--
          <a-form-item>
            <a-input
              size="large"
              type="text"
              :placeholder="$t('user.login.mobile.placeholder')"
              v-decorator="[
                'mobile',
                {
                  rules: [{ required: true, pattern: /^1[34578]\d{9}$/, message: $t('user.login.mobile.placeholder') }],
                  validateTrigger: 'change',
                },
              ]"
            >
              <a-icon slot="prefix" type="mobile" :style="{ color: 'rgba(0,0,0,.25)' }" />
            </a-input>
          </a-form-item>
-->

          <a-row :gutter="16">
            <a-col class="gutter-row" :span="16">
              <a-form-item>
                <a-input
                  size="large"
                  type="text"
                  :placeholder="$t('user.login.mobile.verification-code.placeholder')"
                  v-decorator="[
                    'captcha',
                    {
                      rules: [{ required: true, message: $t('user.verification-code.required') }],
                      validateTrigger: 'blur',
                    },
                  ]"
                >
                  <a-icon slot="prefix" type="mail" :style="{ color: 'rgba(0,0,0,.25)' }" />
                </a-input>
              </a-form-item>
            </a-col>
            <a-col class="gutter-row" :span="8">
              <a-button
                class="getCaptcha"
                tabindex="-1"
                :disabled="state.smsSendBtn"
                @click.stop.prevent="(e)=>sendEmailCode(e,'email')"
                v-text="(!state.smsSendBtn && $t('user.register.get-verification-code')) || state.time + ' s'"
              ></a-button>
            </a-col>
          </a-row>
        </a-tab-pane>
        <a-tab-pane key="otpTab" v-if="state.otpTab" tab="动态令牌认证">
          <a-row :gutter="16">
            <a-col class="gutter-row" :span="16">
              <a-form-item>
                <a-input
                  size="large"
                  type="text"
                  :placeholder="$t('user.login.mobile.verification-code.placeholder')"
                  v-decorator="[
                    'captcha',
                    {
                      rules: [{ required: true, message: $t('user.verification-code.required') }],
                      validateTrigger: 'blur',
                    },
                  ]"
                >
                  <a-icon slot="prefix" type="scan" :style="{ color: 'rgba(0,0,0,.25)' }" />
                </a-input>
              </a-form-item>
            </a-col>
            <a-col class="gutter-row" :span="8">
              <a-button
                class="getCaptcha"
                tabindex="-1"
                :disabled="state.smsSendBtn"
                @click.stop.prevent="(e)=>sendEmailCode(e,'otp')"
                v-text="(!state.smsSendBtn && $t('user.register.get-verification-code')) || state.time + ' s'"
              ></a-button>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>

      <a-form-item>
        <!--        <a-checkbox v-decorator="['rememberMe', { valuePropName: 'checked' }]">{{
          $t('user.login.remember-me')
        }}</a-checkbox>-->
        <router-link :to="{ name: 'recover', params: { user: 'aaa' } }" class="forge-password" style="float: right">{{
          $t('user.login.forgot-password')
        }}</router-link>
      </a-form-item>

      <a-form-item style="margin-top: 24px">
        <a-button
          size="large"
          type="primary"
          htmlType="submit"
          class="login-button"
          :loading="state.loginBtn"
          :disabled="state.loginBtn"
        >{{ $t('user.login.login') }}</a-button
        >
      </a-form-item>

      <div class="user-login-other">
        <span>{{ $t('user.login.sign-in-with') }}</span>
        <a>
          <a-icon class="item-icon" type="alipay-circle"></a-icon>
        </a>
        <a>
          <a-icon class="item-icon" type="taobao-circle"></a-icon>
        </a>
        <a>
          <a-icon class="item-icon" type="weibo-circle"></a-icon>
        </a>
        <router-link class="register" :to="{ name: 'register' }">{{ $t('user.login.signup') }}</router-link>
      </div>
    </a-form>

    <two-step-captcha
      v-if="requiredTwoStepCaptcha"
      :visible="stepCaptchaVisible"
      @success="stepCaptchaSuccess"
      @cancel="stepCaptchaCancel"
    ></two-step-captcha>
  </div>
</template>

<script>
// import md5 from 'md5'
import TwoStepCaptcha from '@/components/tools/TwoStepCaptcha'
import { mapActions } from 'vuex'
import { timeFix } from '@/utils/util'
import { getCaptcha } from '@/api/login'

export default {
  components: {
    TwoStepCaptcha
  },
  data () {
    return {
      customActiveKey: 'firstTab',
      loginBtn: false,
      // login type: 0 email, 1 username, 2 telephone
      loginType: 0,
      isLoginAlert: false,
      loginAlertType: 'error',
      requiredTwoStepCaptcha: false,
      stepCaptchaVisible: false,
      b64s: '',
      cid: '',
      loginPolicyParam: {
        loginName: '',
        policyParam: {
          sessionId: '',
          policyCode: '',
          operate: 0,
          val: ''
        }
      },
      form: this.$form.createForm(this),
      loginAlertMessage: '',
      state: {
        time: 60,
        loginBtn: false,
        // login type: 0 email, 1 username, 2 telephone
        loginType: 0,
        smsSendBtn: false,
        firstTab: true,
        emailTab: false,
        otpTab: false
      }
    }
  },
  created () {
    this.refreshCode()
    /*
    get2step({})
      .then((res) => {
        this.requiredTwoStepCaptcha = res.result.stepCode
      })
      .catch(() => {
        this.requiredTwoStepCaptcha = false
      })
      */
    // this.requiredTwoStepCaptcha = true
  },
  methods: {
    ...mapActions(['Login', 'Logout']),
    refreshCode () {
      getCaptcha().then(res => {
        if (!res.code || res.code !== 200) {
          this.$message.error('请求服务器异常,请联系管理员')
          return
        }
        this.b64s = res.data.b64s
        this.cid = res.data.cid
      })
    },
    // handler
    handleUsernameOrEmail (rule, value, callback) {
      const { state } = this
      const regex = /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$/
      if (regex.test(value)) {
        state.loginType = 0
      } else {
        state.loginType = 1
      }
      callback()
    },
    handleTabClick (key) {
      this.customActiveKey = key
      this.form.resetFields()
      switch (key) {
        case 'emailTab':
          this.loginPolicyParam.policyParam.policyCode = 'email'
          break
        case 'otpTab':
          this.loginPolicyParam.policyParam.policyCode = 'otp'
          break
      }
    },
    doLogin: function (loginParams) {
      this.Login(loginParams)
        .then((res) => {
          const { code, message, data } = res
          if (code === 200) {
            this.loginSuccess(res)
            return
          } else if (code === 201) {
            // 加载策略
            this.loginPolicyParam = {
              loginName: loginParams.loginName,
              policyParam: {
                sessionId: data.sessionId
              }
            }
            this.state.time = -1
            this.state.firstTab = false
            this.state.emailTab = false
            this.state.otpTab = false
            if (data.policyCode.includes('otp')) {
              this.state.otpTab = true
              this.handleTabClick('otpTab')
              // this.$forceUpdate()
            }
            if (data.policyCode.includes('email')) {
              this.state.emailTab = true
              this.handleTabClick('emailTab')
              // this.$forceUpdate()
            }
          } else if (code === 203) {
            if (loginParams.policyParam && loginParams.policyParam.operate === 0) {
              const {
                state
              } = this
              const interval = window.setInterval(() => {
                if (state.time-- <= 0) {
                  state.time = 60
                  state.smsSendBtn = false
                  window.clearInterval(interval)
                }
              }, 1000)
              this.state.smsSendBtn = true
            }
          }
          this.isLoginAlert = true
          this.loginAlertType = 'info'
          this.loginAlertMessage = message
        }).catch((err) => this.requestFailed(err))
        .finally(() => {
          this.state.loginBtn = false
        })
    },
    handleSubmit (e) {
      e.preventDefault()
      const {
        form: { validateFields },
        state,
        customActiveKey,
        cid
      } = this

      state.loginBtn = true

      const validateFieldsKey = customActiveKey === 'firstTab' ? ['username', 'password', 'code'] : ['captcha']

      validateFields(validateFieldsKey, { force: true }, (err, values) => {
        if (!err) {
          let loginParams
          if (customActiveKey === 'firstTab') {
            loginParams = { ...values, cid: cid }
            delete loginParams.username
            // console.log('login form', loginParams)
            loginParams.loginName = values.username
          } else {
            this.loginPolicyParam.policyParam.val = values.captcha
            this.loginPolicyParam.policyParam.operate = 1
            loginParams = this.loginPolicyParam
          }
          // const loginParams = { ...values, cid: cid }
          // delete loginParams.username
          // loginParams.loginName = values.username

          // console.log('login form', loginParams)
          // loginParams[!state.loginType ? 'email' : 'username'] = values.username
          // loginParams.password = md5(values.password)
          this.doLogin(loginParams)
        } else {
          setTimeout(() => {
            state.loginBtn = false
          }, 600)
        }
      })
    },
    sendEmailCode (e, code) {
      e.preventDefault()

      this.loginPolicyParam.policyParam.operate = 0
      this.loginPolicyParam.policyParam.policyCode = code
      this.doLogin(this.loginPolicyParam)

      // const hide = this.$message.loading('验证码发送中..', 0)
     /* const {
        form: { validateFields },
        state
      } = this */

      /* validateFields(['mobile'], { force: true }, (err, values) => {
        if (!err) {
          state.smsSendBtn = true

          const interval = window.setInterval(() => {
            if (state.time-- <= 0) {
              state.time = 60
              state.smsSendBtn = false
              window.clearInterval(interval)
            }
          }, 1000)

          const hide = this.$message.loading('验证码发送中..', 0)
          getSmsCaptcha({ mobile: values.mobile })
            .then((res) => {
              setTimeout(hide, 2500)
              this.$notification['success']({
                message: '提示',
                description: '验证码获取成功，您的验证码为：' + res.result.captcha,
                duration: 8
              })
            })
            .catch((err) => {
              setTimeout(hide, 1)
              clearInterval(interval)
              state.time = 60
              state.smsSendBtn = false
              this.requestFailed(err)
            })
        }
      }) */
    },
    stepCaptchaSuccess () {
      this.loginSuccess()
    },
    stepCaptchaCancel () {
      this.Logout().then(() => {
        this.loginBtn = false
        this.stepCaptchaVisible = false
      })
    },
    loginSuccess (res) {
      // console.log(res)
      const { data } = res
      // check res.homePage define, set $router.push name res.homePage
      // Why not enter onComplete
      /*
      this.$router.push({ name: 'analysis' }, () => {
        console.log('onComplete')
        this.$notification.success({
          message: '欢迎',
          description: `${timeFix()}，欢迎回来`
        })
      })
      */
      this.$router.push({ path: '/' })
      // 延迟 1 秒显示欢迎信息
      setTimeout(() => {
        this.$notification.success({
          message: '欢迎 ' + data.user.userName,
          description: `${timeFix()}，欢迎回来`
        })
      }, 1000)
      this.isLoginAlert = false
    },
    requestFailed (err) {
      this.isLoginAlert = true
      this.loginAlertType = 'error'
      this.loginAlertMessage = err.message
      if (this.customActiveKey === 'firstTab') {
        this.refreshCode()
      }
     /* this.$notification['error']({
        message: '错误',
        description: ((err.response || {}).data || {}).message || '请求出现错误，请稍后再试',
        duration: 4
      }) */
    }
  }
}
</script>

<style lang="less" scoped>
.user-layout-login {
  label {
    font-size: 14px;
  }

  .getCaptcha {
    display: block;
    width: 100%;
    height: 40px;
  }

  .forge-password {
    font-size: 14px;
  }

  button.login-button {
    padding: 0 15px;
    font-size: 16px;
    height: 40px;
    width: 100%;
  }

  .user-login-other {
    text-align: left;
    margin-top: 24px;
    line-height: 22px;

    .item-icon {
      font-size: 24px;
      color: rgba(0, 0, 0, 0.2);
      margin-left: 16px;
      vertical-align: middle;
      cursor: pointer;
      transition: color 0.3s;

      &:hover {
        color: #1890ff;
      }
    }

    .register {
      float: right;
    }
  }
}
</style>
