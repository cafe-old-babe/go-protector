<template>
  <div>
    <a-card :bordered="false" style="margin-bottom: 24px" size="small"  title="字典数据列表">
      <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 14 }" :form="searchForm" @submit="handleSearch">
          <a-row :gutter="[8,1]">
              <a-col :span="12">

                  <a-form-item label="数据名称">
                      <a-input
                          v-decorator="[ 'dataName', {
                            rules: [
                              {
                                required: false,
                                message: 'Input something!',
                              },
                            ],
                          },
                        ]" placeholder="请输入数据名称"
                      />
                  </a-form-item>

              </a-col>
              <a-col :span="12">
                  <a-form-item label="数据编码">
                      <a-input
                          v-decorator="[ 'dataCode', {
                            rules: [
                              {
                                required: false,
                                message: 'Input something!',
                              },
                            ],
                          },
                        ]" placeholder="请输入数据编码"
                      />
                  </a-form-item>
              </a-col>
          </a-row>

          <a-row :gutter="[8,8]">
              <a-col :span="12" :style="{ textAlign: 'left' }">
                  <a-space :size="8">
                    <a-button @click="addNew" type="primary">新建</a-button>
                    <a-button type="danger">批量删除</a-button>

<!--                            <a-button >批量操作</a-button>
                      <a-dropdown>
                          &lt;!&ndash;                    <a-menu @click="handleMenuClick" slot="overlay">&ndash;&gt;
                          &lt;!&ndash;                        <a-menu-item key="delete">删除</a-menu-item>&ndash;&gt;
                          &lt;!&ndash;                        <a-menu-item key="audit">审批</a-menu-item>&ndash;&gt;
                          &lt;!&ndash;                    </a-menu>&ndash;&gt;
                          <a-button>
                              更多操作 <a-icon type="down" />
                          </a-button>
                      </a-dropdown>-->
                  </a-space>
              </a-col>
              <a-col :span="12" :style="{ textAlign: 'right' }">
                  <a-space :size="8">
                      <a-button type="primary" html-type="submit">
                          查询
                      </a-button>
                      <a-button @click="handleReset">
                          Clear
                      </a-button>

                  </a-space>
              </a-col>
          </a-row>
      </a-form>
    </a-card>

    <a-card :bordered="false" style="">
        <StandardTable
            :columns="columns"
            :loading="loading"
            rowKey="id"
            :data-source=dataSource
            :selectedRows.sync="selectedRows"
            @selectedRowChange="onSelectChange"
            :pagination="{...pagination, onChange: onPageChange}"
        >

        <div slot="status" slot-scope="{text,record}">
          <a >
            <a-tooltip placement="left" :title="record.dataStatus==='0'?'点击停用':'点击锁定'"
                       :get-popup-container="getPopupContainer">
              <a-tag :color="record.dataStatus===0?'green':'red'"
                     @click="changeStatus(record)"
              >
                {{ record.dataStatusText }}
              </a-tag>
            </a-tooltip>
          </a>
        </div>


      <div slot="action" slot-scope="{text, record}">
        <a style="margin-right: 8px"
           @click="editRecord(record)"
        >
          <a-icon type="edit"/>编辑
        </a>
        <a @click="deleteRecord(record.id)">
          <a-icon type="delete" />删除
        </a>
      </div>
        </StandardTable>
    </a-card>
    <Edit
      :visible="editVisible"
      :record="record"
      @close="editClose"
      @ok="editOk"
    />
  </div>

</template>

<script>
import StandardTable from '@/components/table/StandardTable'
import {request} from '@/utils/request'
import column from "../column";
import Edit from "@/pages/system/manager/dict/data/Edit.vue";

export default {
  name: 'DataList',
  components: {StandardTable, Edit},
  props: {
    typeCode: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      searchForm: this.$form.createForm(this),
      searchFormParam: {},
      loading: false,
      columns: column.dataColumn,
      dataSource: [],
      selectedRows: [],
      editVisible: false,
      record: {},
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0
      }
    };
  },
  watch: {
    typeCode() {
      this.getData();
    }
  },
  mounted() {
    this.getData()
  },
  computed: {},
  methods: {
    editOk() {
      this.editClose()
      this.getData()
    },
    editClose() {
      this.editVisible = false
      this.record = {}
    },
    editRecord(record) {
      this.record = record;
      this.editVisible = true
    },

    changeStatus(record) {
      let url = "/api/dict/data/"+record.id+"/"+(record.dataStatus ^ 1)
      console.log(url)
      request(url).then(res=> {
        console.log(res)
        let resData = res?.data ?? {};
        if (resData.code === 200) {
          this.getData();
        } else {
          this.$message.warning(resData.message);
        }
      })
      console.log(record, "changeStatus")
    },
    getPopupContainer(trigger) {
      return trigger.parentElement;
    },
    onPageChange(page, pageSize) {
      this.pagination.current = page
      this.pagination.pageSize = pageSize
      this.getData()
    },

    handleSearch(e) {
      e.preventDefault();
      this.searchForm.validateFields((error, values) => {
        if (error) {
          console.log('error', error)
          return
        }
        this.searchFormParam = values;
        this.getData();
      });
    },
    getData() {
      if (this.typeCode === '') {
        return;
      }
      this.loading = true;

      request('/api/dict/data',
        {...this.pagination, ...this.searchFormParam, typeCode: this.typeCode}).then(res => {
        const {code, message, data: {list, current, pageSize, total}} = res?.data ?? {}
        if (code !== 200) {
          this.$message.warning(message)
          return
        }
        this.dataSource = list;
        this.pagination.current = current
        this.pagination.pageSize = pageSize
        this.pagination.total = total
      }).finally(() => this.loading = false)
    },
    handleReset() {
      this.searchForm.resetFields()
    },
    deleteRecord(ids) {
      console.log(ids, "delete")

    },
    remove() {
      this.dataSource = this.dataSource.filter(item => this.selectedRows.findIndex(row => row.key === item.key) === -1)
      this.selectedRows = []
    },
    onClear() {
      this.$message.info('您清空了勾选的所有行')
    },


    onSelectChange() {
      console.log(this.selectedRows)
      this.$message.info('选中行改变了')
    },
    addNew() {
      this.$message.info("addNew");
    },
  }
};
</script>

<style lang="less" scoped>
.search{
    //margin-bottom: 54px;
}
.fold{
    width: calc(100% - 216px);
    display: inline-block
}
.operator{
    margin: 10px;
    //margin-bottom: 18px;
}
@media screen and (max-width: 900px) {
    .fold {
        width: 100%;
    }
}

.ant-advanced-search-form {
    //padding: 12px;
    //background: #fbfbfb;
    //border: 1px solid #d9d9d9;
    //border-radius: 6px;
}

.ant-advanced-search-form .ant-form-item {
    display: flex;
    //display: flow;
}

.ant-advanced-search-form .ant-form-item-control-wrapper {
    flex: 1;
}

#components-form-demo-advanced-search .ant-form {
    max-width: none;
}

</style>
