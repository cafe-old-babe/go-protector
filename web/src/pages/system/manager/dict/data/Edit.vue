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
        <a-form-model ref="refForm"
                      :model="localRecord"
                      :rules="rules"
                      :label-col="{ span: 6 }"
                      :wrapper-col="{ span: 14 }"
                      layout="horizontal">
          <a-input v-model="localRecord.id" v-show='false'/>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="数据名称"
            prop="dataName"
          >
            <a-input  v-model="localRecord.dataName" placeholder="请输入数据名称"/>
          </a-form-model-item>
          <a-form-model-item
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            label="数据编码"
            prop="dataCode"
          >
            <a-input v-model="localRecord.dataCode"
                     placeholder="请输入数据编码"
             />
          </a-form-model-item>
          <a-form-model-item
            label="数据状态"
            :label-col="formItemLayout.labelCol"
            :wrapper-col="formItemLayout.wrapperCol"
            prop="dataStatus"
          >
            <a-radio-group v-model="localRecord.dataStatus" :default-value="localRecord.dataStatus" button-style="solid">
              <a-radio :value="0">
                启用
              </a-radio>
              <a-radio :value="1">
                停用
              </a-radio>
            </a-radio-group>
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
import {request} from "@/utils/request";

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
  },
  data() {
    return {
      loading: true,
      localRecord: {},
      rules: {
        dataStatus: [
          {required: true, message: '请选择是否启用',}
        ],
        dataCode: [
          {required: true, message: '请输入数据编码',}
        ],
        dataName: [
          {required: true, message: '请输入数据名称',}
        ]
      }
    };
  },
  watch:{
    record(val) {
      this.localRecord = val
      this.loading = false
    }
  },
  computed: {
    formItemLayout() {
      const { formLayout } = this;
      return formLayout === 'horizontal'
        ? {
          labelCol: { span: 2 },
          wrapperCol: { span: 14 },
        }
        : {};
    },
  },
  methods: {
    onClose() {
      this.$emit('close');
    },
    handleSave() {
      this.loading = true
      this.$refs.refForm.validate(valid => {

        if (!valid) {
          this.loading = false
          return false;
        }

        request("/api/dict/data/save",this.localRecord).then(res => {


          const {code, message} = res?.data ?? {}
          if (code === 200) {
            this.$emit('ok');
            this.$message.success(message);
            return

          }
          this.$message.error(message);

        }).finally(() => this.loading = false)
      });

    }
  },
};
</script>
