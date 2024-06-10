<script>
import { accountColumns } from './component.js'

export default {

  name: 'SelectAccount',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    accounts: {
      type: Array,
      required: true,
      default: () => []
    }
  },
  watch: {
  },
  data() {
    return {
      columns: accountColumns

    }
  },
  methods: {
    handleCancel: function () {
      this.$emit('close')
    },
    dailColor(dailStatus) {
      switch (dailStatus) {
        case '0':
          return 'red'
        case '1':
          return 'green'
        default:

          return 'grey'
      }
    },
    login: function (current) {
      this.$emit('login', current)
    }
  }
}
</script>

<template>
  <a-modal
    v-model="visible"
    title="请选择从帐号"
    :closable="false"
    :maskClosable="false"
    :width="800">
    <template v-slot:footer>
      <a-button key="back" @click="handleCancel">
        取消
      </a-button>
    </template>

    <a-table
      ref="table"
      rowKey="id"
      size="default"
      :show-pagination="false"
      :showSizeChanger="false"
      :columns="columns"
      :data-source="accounts"
    >
      <template v-slot:action="text,current">

        <a style="margin-right: 8px" @click="login(current)">
          <a-icon type="code"/>
          登录
        </a>
        <!--          </span>-->
      </template>
      <template v-slot:dailStatus="text,current">
        <a-tooltip
          placement="left"
          :title="current.assetAccount.dailMsg"
          :get-popup-container="(trigger) => trigger.parentElement">
          <a-tag
            :color="dailColor(current.assetAccount.dailStatus)">
            {{
              current.assetAccount.dailStatusText
            }}
          </a-tag>
        </a-tooltip>
      </template>
    </a-table>
  </a-modal>

</template>

<style scoped lang="less">

</style>
