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
            label="用户姓名"
            prop="username"
          >
            <a-input v-model="localRecord.username" placeholder="请输入用户姓名" />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="用户帐号"
            prop="loginName"
          >
            <a-input v-model="localRecord.loginName" placeholder="请输入用户帐号"/>
          </a-form-model-item>
          <a-form-model-item
            v-if="!localRecord.id"
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="用户密码"
            prop="password"
          >
            <a-input-password v-model="localRecord.password" placeholder="请输入用户密码"/>
          </a-form-model-item>

          <a-form-model-item
            :disabled="localRecord.id"
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="用户有效期"
            prop="expirationAt"
          >
            <a-date-picker
              format="YYYY-MM-DD HH:mm:ss"
              show-time
              placeholder="请选择有效期"
              @change="changeExpirationAt"
            />
          </a-form-model-item>

          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="性别"
            prop="sex"
          >
            <a-radio-group :default-value="localRecord.sex" button-style="solid" v-model="localRecord.sex">
              <a-radio-button v-for="sex in sexData" :key="sex.id" :value="sex.code">
                {{ sex.text }}
              </a-radio-button>
            </a-radio-group>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="角色"
            prop="roleIds"
          >
            <a-select
              mode="multiple"
              placeholder="请选择角色"
              v-model="localRecord.roleIds"
              style="width: 100%"

            >
              <a-select-option v-for="item in roleList" :key="item.id" :value="item.id">
                {{ item.roleName }}
              </a-select-option>
            </a-select>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="所属部门"
            prop="deptId"
          >
            <a-tree-select
              v-model="localRecord.deptId"
              style="width: 100%"
              :dropdown-style="{ maxHeight: '200px', overflow: 'auto' }"
              :replace-fields="replaceFields"
              :tree-data="deptTreeData"
              placeholder="请选择所属部门"
              @select="selectDept"
              tree-default-expand-all
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="岗位"
            prop="postIds"
          >
            <a-select
              mode="multiple"
              placeholder="请选择岗位"
              v-model="localRecord.postIds"
              style="width: 100%"
              @change="changePost"
              :options="postList"
            >
              <!--              <a-select-option-->
              <!--                v-for="item in postList"-->
              <!--                :key="item.id"-->
              <!--                :value="item.id">-->
              <!--                {{ item.name }}-->
              <!--              </a-select-option>-->
            </a-select>
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
// import { Radio } from 'ant-design-vue'
import { loadDept } from '@/api/user'
import { loadDictData } from '@/api/common'

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
  async mounted() {
    this.loadDeptData()
    this.sexData = await loadDictData('sex')
    this.loadRoleData()
  },
  data() {
    return {
      loading: true,
      localRecord: {},
      sexData: [],
      postList: [],
      roleList: [],
      rules: {
        loginName: [
          { required: true, message: '请输入登录帐号' }
        ],
        username: [
          { required: true, message: '请输入类型名称' }
        ],
        sex: [
          { required: true, message: '请输入性别' }
        ],
        roleIds: [
          { required: true, message: '请选择角色' }
        ],
        deptId: [
          { required: true, message: '请选择部门' }
        ],
        postIds: [
          { required: true, message: '请选择岗位' }
        ],
        password: [
          { required: true, message: '请输入用户密码' }
        ]
      },
      replaceFields: {
        key: 'id',
        value: 'id',
        title: 'name'
      },
      deptTreeData: []
    }
  },
  watch: {
    record(val) {
      this.localRecord = Object.assign({}, val)
      this.loadPostData(this.localRecord.deptId)
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
        console.log(this.localRecord)
        request.post('/user/save', this.localRecord).then(res => {
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
    loadDeptData() {
      loadDept().then(res => {
        const { code, data, message } = res
        if (code !== 200) {
          this.$message.error(message)
        }
        this.deptTreeData = data.children
      })
    },
    loadPostData(deptId) {
      if (!deptId) {
        this.postList = []
        return
      }
      this.loading = true

      request.post('/post/list/' + deptId).then(res => {
        const { code, data, message } = res
        if (code !== 200) {
          this.$message.error(message)
          return
        }
        // this.postList = data
        data.forEach(e => {
          this.postList.push({
            value: e.id,
            label: e.name
          })
        })
        this.$forceUpdate()
      }).finally(() => {
        this.loading = false
      })
    },
    selectDept(key) {
      this.postList = []
      this.localRecord.postIds = []
      this.loadPostData(key)
    },
    changePost(key) {
      console.log(key)
      this.localRecord.postIds = key
      this.$forceUpdate()
    },
    loadRoleData() {
      request.post('/role/list').then(res => {
        const { code, data, message } = res
        if (code !== 200) {
          this.$message.error(message)
          return
        }
        this.roleList = data
      })
    },
    changeExpirationAt(date, dateString) {
      console.log(date)
      this.localRecord.expirationAt = dateString
    }
  }
}
</script>
