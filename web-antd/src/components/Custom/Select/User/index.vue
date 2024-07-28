<!--https://v2.cn.vuejs.org/v2/guide/components-custom-events.html#%E8%87%AA%E5%AE%9A%E4%B9%89%E7%BB%84%E4%BB%B6%E7%9A%84-v-model-->
<script>
// 6-10	【实战】资源编辑页面-掌握封装自定义选择器公共组件搞定v-model双向绑定）
import selectIndex from './components/SelectIndex.vue'
export default {
  name: 'SelectUser',
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
    placeholder: {
      type: String,
      required: false,
      default: '请选择用户'
    },
    showOperate: {
      type: Boolean,
      default: true
    }
  },
  watch: {
    showLabel(val) {
      this.localShowLabel = val
    }
  },
  data() {
    return {
      localValue: this.value,
      localShowLabel: this.showLabel
    }
  },
  computed: {
    calcWidth() {
      return this.showOperate ? 'calc(100% - 68px)' : 'calc(100% + 0px)'
    }
  },
  methods: {
    showSelectUser: function () {
      const innerSelf = this
      this.$dialog(selectIndex, {
        on: {
          select (rowIds, rowData) {
            innerSelf.selectChange(rowIds, rowData)
          }
        }
      }, {
        title: '请选择用户',
        okText: '确认',
        cancelText: '取消',
        width: '70%',
        centered: true
      })
    },
    removeUser: function () {
      this.localValue = ''
      this.localShowLabel = ''
    },
    selectChange: function (rowIds, rowData) {
      this.$emit('change.value', rowIds[0])
      this.localShowLabel = rowData[0].username
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
      <a-tooltip placement="top" title="选择用户" :style="{marginRight: `2px`}">
        <a-button @click="showSelectUser" type="primary" icon="select"></a-button>
      </a-tooltip>
      <a-tooltip placement="top" title="清除用户">
        <a-button @click="removeUser" type="danger" icon="delete"></a-button>
      </a-tooltip>
    </template>
  </div>

</template>

<style scoped lang="less">

</style>
