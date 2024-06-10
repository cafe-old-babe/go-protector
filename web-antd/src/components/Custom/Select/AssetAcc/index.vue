<!--https://v2.cn.vuejs.org/v2/guide/components-custom-events.html#%E8%87%AA%E5%AE%9A%E4%B9%89%E7%BB%84%E4%BB%B6%E7%9A%84-v-model-->
<script>
import selectIndex from './components/SelectIndex.vue'
export default {
  name: 'SelectAssetAcc',
  components: { selectIndex },
  model: {
    prop: 'value',
    event: 'change.value'
  },
  props: {
    showLabel: {
      type: [String],
      required: false,
      default: ''
    },
    value: {
      type: [String, Number],
      required: false,
      default: ''
    },
    userId: {
      type: [Number],
      required: false,
      default: 0
    },
    showOperate: {
      type: Boolean,
      default: true
    },
    placeholder: {
      type: String,
      required: false,
      default: '请选择资源从帐号'
    }
  },
  watch: {
    showLabel(val) {
      this.localShowLabel = val
    }
  },
  computed: {
    calcWidth() {
      return this.showOperate ? 'calc(100% - 68px)' : 'calc(100% + 0px)'
    }
  },
  data() {
    return {
      localValue: this.value,
      localShowLabel: this.showLabel
    }
  },
  methods: {
    showSelect: function () {
      const innerSelf = this
      this.$dialog(selectIndex, {
        userId: this.userId,
        on: {
          select (rowIds, rowData) {
            innerSelf.selectChange(rowIds, rowData)
          }
        }
      }, {
        title: '选择从帐号',
        okText: '确认',
        cancelText: '取消',
        width: '70%',
        centered: true
      })
    },
    remove: function () {
      this.localValue = ''
      this.localShowLabel = ''
    },
    selectChange: function (rowIds, rowData) {
      this.$emit('change.value', rowIds[0])
      this.localShowLabel = `${rowData[0].account} [${rowData[0].assetBasic.assetName}(${rowData[0].assetBasic.ip})]`
      this.$emit('callback', rowData)
    }
  }
}
</script>

<template>
  <div>
    <a-input
      type="hidden"
      :value="value"/>
    <a-input
      disabled
      v-model="localShowLabel"
      :placeholder="placeholder"
      :style="{width:`${calcWidth}`,marginRight: `2px`}"/>
    <template v-if="showOperate">
      <a-tooltip placement="top" title="选择从帐号" :style="{marginRight: `2px`}">
        <a-button @click="showSelect" type="primary" icon="select"></a-button>
      </a-tooltip>
      <a-tooltip placement="top" title="清除从帐号">
        <a-button @click="remove" type="danger" icon="delete"></a-button>
      </a-tooltip>
    </template>
  </div>

</template>

<style scoped lang="less">

</style>
