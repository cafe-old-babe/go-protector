<template>
  <div>
    <a-input-search v-if="_props.isSearch" style="margin-bottom: 8px" placeholder="Search" @change="onChange" />
    <a-tree
      :loading="loading"
      :checkable="$props.checkable"
      :expanded-keys.sync="expandedKeys"
      :auto-expand-parent="false"
      :replace-fields="$props.replaceFields"
      :data="$props.data"
      :tree-data="treeData"
      @loadDone="$props.loadDone"
      @check="$props.check"
      @select="select"
    >
      <template slot="title" slot-scope="treeNode">
        <template>
          <span
            :style="{width:showOperateBtn ? 'calc(100% - 80px)' : '100%'}"
            :title="treeNode[replaceFields.title]"
            v-if="treeNode[replaceFields.title].indexOf(searchValue) > -1">
            {{ treeNode[replaceFields.title].substr(0, treeNode[replaceFields.title].indexOf(searchValue)) }}
            <span style="color: #f50">{{ searchValue }}</span>
            {{ treeNode[replaceFields.title].substr(treeNode[replaceFields.title].indexOf(searchValue) + searchValue.length) }}
          </span>
          <span :title="treeNode[replaceFields.title]" v-else>{{ treeNode[replaceFields.title] }}</span>
        </template>
        <template v-if="showOperateBtn">
          <a-tooltip title="删除" v-if="showDeleteBtn(treeNode)">
            <a href="#" class="red" style="margin-right: 8px" @click="$emit('deleteTreeNode',treeNode)" >
              <a-icon type="delete" />
            </a>
          </a-tooltip>
          <a-tooltip title="修改" v-if="showUpdateBtn(treeNode)">
            <a href="#" style="margin-right: 8px" @click="$emit('updateTreeNode',treeNode)">
              <a-icon type="edit" />
            </a>
          </a-tooltip>
          <a-tooltip title="新建" v-if="showAddBtn(treeNode)">
            <a href="#" @click="$emit('addTreeNode',treeNode)">
              <a-icon type="plus" />
            </a>
          </a-tooltip>
        </template>
      </template>
    </a-tree>
  </div>
</template>
<script>
const innerProps = {
  loading: {
    type: Boolean,
    required: false,
    default: () => false
  },
  // 加载数据 Promise
  data: {
    type: Function,
    required: true
  },
  // 加载完毕
  loadDone: {
    type: Function,
    required: false,
    default: (treeData, checkedKeys) => {
    }
  },
  // 节点前添加 Checkbox 复选框
  checkable: {
    type: Boolean,
    default: false
  },
  isSearch: {
    type: Boolean,
    default: true
  },
  // 选中回调
  check: {
    type: Function,
    required: false,
    default: (checkedKeys) => {
      console.log(checkedKeys)
    }
  },
  // 替换默认字段
  replaceFields: {
    type: Object,
    default: () => {
      return {
        value: 'id',
        key: 'id',
        title: 'name'
      }
    }
  },
  // 是否显示操作字段
  showOperateBtn: {
    type: Boolean,
    default: () => true
  },
  showDeleteBtn: {
    type: Function,
    default: (treeNode) => {
      return treeNode.children && treeNode.children.length <= 0
    }
  },
  showUpdateBtn: {
    type: Function,
    default: (treeNode) => {
      return true
    }
  },
  showAddBtn: {
    type: Function,
    default: (treeNode) => {
      return true
    }
  }
}
export default {
  name: 'CTree',

  props: Object.assign({}, innerProps),
  watch: {
    treeData() {
      if (this.isSearch) {
        this.generateList(this.treeData)
      }
      const temp = []
      this.treeData.forEach(e => {
        if (e[this.replaceFields.key].selected) {
          temp.push(e[this.replaceFields.key])
        }
      })

      if (temp.length <= 0) {
        this.treeData.forEach(elem => {
          if (elem.children.length > 0) {
            temp.push(elem[this.replaceFields.key])
            elem.children.forEach(e => {
              temp.push(e[this.replaceFields.key])
            })
          }
        })
      }
      this.expandedKeys = temp
    }
  },
  data() {
    return {
      // 数据
      treeData: [],
      // 搜索
      searchValue: '',
      // 用于搜索的数据副本
      dataList: [],
      // 选中的key
      expandedKeys: []
    }
  },
  created() {
    this.loadData()
  },
  methods: {
    loadData() {
      const result = this.data()
      result.then((res) => {
        const { code, data, message } = res
        if (code === 200) {
          this.treeData = data
        } else {
          this.$message.warn(message)
        }
      }).finally(() => {
        this.loadDone(this.treeData, this.expandedKeys)
      })
    },
    // 生成用于搜索时的数据
    generateList (data) {
      for (let i = 0; i < data.length; i++) {
        const node = data[i]
        const id = node[this.replaceFields.key]
        const name = node[this.replaceFields.title]

        this.dataList.push({ [this.replaceFields.key]: id, [this.replaceFields.title]: name })
        if (node.children) {
          this.generateList(node.children)
        }
      }
    },
    // 匹配搜索的数据的父节点
    getParentKey (key, tree) {
      let parentKey
      for (let i = 0; i < tree.length; i++) {
        const node = tree[i]
        if (node.children) {
          if (node.children.some(item => item[this.replaceFields.key] === key)) {
            parentKey = node[this.replaceFields.key]
          } else if (this.getParentKey(key, node.children)) {
            parentKey = this.getParentKey(key, node.children)
          }
        }
      }
      return parentKey
    },
    // 搜索
    onChange(e) {
      const value = e.target.value
      const expandedKeys = this.dataList
        .map(item => {
          if (item[this.replaceFields.title].indexOf(value) > -1) {
            return this.getParentKey(item[this.replaceFields.key], this.treeData)
          }
          return null
        })
        .filter((item, i, self) => item && self.indexOf(item) === i)
      Object.assign(this, {
        expandedKeys,
        searchValue: value,
        autoExpandParent: true
      })
    },
    select: function (selectedKeys, e) {
      this.$emit('select', selectedKeys, e)
    }
  }

}
</script>

<style scoped >

</style>
