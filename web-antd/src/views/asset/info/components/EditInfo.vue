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
            label="资源名称"
            prop="assetName"
          >
            <a-input v-model="localRecord.assetName" placeholder="请输入资源名称" />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="所属资源组"
            prop="groupId"
          >
            <a-tree-select
              v-model="localRecord.groupId"
              style="width: 100%"
              :dropdown-style="{ maxHeight: '200px', overflow: 'auto' }"
              :replace-fields="replaceFields"
              :tree-data="groupTreeData"
              placeholder="请选择所属所属资源组"
              tree-default-expand-all
            />
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="资源IP"
            prop="ip"
          >
            <a-input v-model="localRecord.ip" placeholder="请输入资源IP"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="端口"
            prop="port"
          >
            <a-input-number
              :min="2"
              :max="65536"
              step="1"
              :precision="0"
              v-model="localRecord.port"
              placeholder="请输入端口"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="特权帐号"
            prop="account"
          >
            <a-input v-model="localRecord.account" placeholder="请输入特权帐号"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="特权帐号密码"
            prop="password"
          >
            <a-input-password v-model="localRecord.password" placeholder="请输入特权帐号密码"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="资源管理员"
            prop="managerUserId"
          >
            <select-user
              ref="selectUserRef"
              v-model="localRecord.managerUserId"
              :show-label="localRecord.managerUsername"
              placeholder="请选择资源管理员"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="网关"
            prop="gatewayId"
          >
            <a-select
              v-model="localRecord.gatewayId"
              :default-value="localRecord.gatewayId"
              placeholder="请选择网关">
              <a-select-option v-for="elem in gatewayData" :key="elem.id">
                {{ elem.agName + '('+ elem.agIp+')' }}
              </a-select-option>
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
          zIndex: 1
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
import { loadGatewayList, loadGroupTree } from '@/api/asset'
import SelectUser from '@/components/Custom/Select/User/'

export default {
  components: { SelectUser },
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
    // this.loadGroupTree()
    // this.sexData = await loadDictData('sex')
  },
  data() {
    return {
      loading: true,
      localRecord: {},
      rules: {
        assetName: [
          { required: true, message: '请输入资源名称' }
        ],
        groupId: [
          { required: true, message: '请选择所属资源组' }
        ],
        ip: [
          { required: true, message: '请输入资源IP' }
        ],
        port: [
          { required: true, message: '请输入端口' }
        ],
        account: [
          { required: true, message: '请输入特权帐号' }
        ],
        password: [
          { required: true, message: '请输入特权帐号密码' }
        ],
        managerUserId: [
          { required: true, message: '请选择责任人' }
        ]
      },
      replaceFields: {
        key: 'id',
        value: 'id',
        title: 'name'
      },
      groupTreeData: [],
      gatewayData: []
    }
  },
  watch: {
     record(val) {
       // 为打开当前页面 不查询
       if (!this.visible) {
         return
       }
       this.loadGroupTree()
       this.loadGateway()
       this.localRecord = Object.assign({}, val)
       console.log(this.localRecord)
       this.rules.password[0].required = !this.localRecord.id
       this.loading = false
    }
  },
  computed: {
    isEdit() {
      return !this.localRecord.id
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
    onClose() {
      this.$refs.selectUserRef.removeUser()
      this.$emit('close')
    },
    handleSave() {
      this.loading = true
      this.$refs.form.validate(valid => {
        if (!valid) {
          this.loading = false
          return false
        }
        request.post('/asset-info/save', this.localRecord).then(res => {
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
    loadGroupTree() {
      loadGroupTree().then(res => {
        const { code, data, message } = res
        if (code !== 200) {
          this.$message.error(message)
          return
        }
        this.groupTreeData = data.children
      })
    },
    async loadGateway() {
      await loadGatewayList().then(res => {
        const { code, data, message } = res
        if (code !== 200) {
          this.$message.error(message)
          return
        }
        this.gatewayData = data
      })
    }
  }
}
</script>
