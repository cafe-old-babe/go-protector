<script>

import STable from '@/components/Table'

// yarn add  asciinema-player@3.4.0
import * as AsciinemaPlayer from 'asciinema-player'
import 'asciinema-player/dist/bundle/asciinema-player.css'
import request from '@/utils/request'
export default {
  name: 'PlayBack',
  components: { STable },
  methods: {
    onClose: function (e) {
      this.$emit('close')
    },
    seek: function (elem) {
      this.player.seek(elem.timeStamp).then(() => {
        this.player.play()
      })
    }
  },
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    castData: {
      type: String,
      required: true
    },
    record: {
      type: Object,
      default: () => undefined,
      required: true
    }
  },
  data() {
    return {
      player: undefined,
      queryParam: {},
      columns: [
        {
          title: '序号',
          width: 40,
          // scopedSlots: { customRender: 'serial' }
          customRender: (text, record, index) => index + 1
        },
        {
          title: '指令',
          scopedSlots: { customRender: 'cmd' }
          // customRender: (text, record) => record.cmd
        }
      ],
      loadData: (parameter) => {
        const id = this.record.id
        const promise = request.post(`/sso-operation/page/${id}`, Object.assign(parameter, this.queryParam)).then((res) => {
          const { code, data, message } = res
          if (code === 200) {
            return data
          }
          this.$message.error(message)
        })
        return promise.catch((error) => {
          this.$message.error(error.message)
          return {
            data: [],
            pageNo: 1,
            pageSize: 10,
            totalCount: 0,
            totalPage: 0
          }
        })
      }
    }
  },
  watch: {
    visible(val) {
      if (!val) {
        return
      }

      this.$nextTick(() => {
        // const element = document.getElementById('player')
        this.player = AsciinemaPlayer.create({ data: this.castData }, this.$refs.player)
        this.$refs.table.refresh(true)
      })
    }
  }

}
</script>

<template>
  <a-modal
    :visible="visible"
    title="录像回放"
    :footer="null"
    width="85%"
    height="60%"
    :dialog-style="{ top: '20px' }"
    @cancel="onClose"
  >
    <a-layout>
      <a-layout-content>
        <div v-if="visible" ref="player"></div>
      </a-layout-content>
      <a-layout-sider
        :width="300"
        style="{background: #cccccc}"
      >
        <a-card :bordered="false" :style="{overflow:'auto'}">
          <div class="table-page-search-wrapper">
            <a-form layout="inline" >
              <a-row :gutter="[10]">
                <a-col :span="16" >
                  <a-form-item >
                    <a-input v-model="queryParam.cmd" placeholder="" />
                  </a-form-item>
                </a-col>
                <a-col :span="5">
                  <a-button type="primary" @click="$refs.table.refresh(true)">查询</a-button>
                </a-col>
              </a-row>
            </a-form>
          </div>
          <s-table
            ref="table"
            rowKey="id"
            size="small"
            :show-pagination="true"
            :showSizeChanger="true"
            :columns="columns"
            :data="loadData"
            :scroll="{y:'calc(60vh - 30px)'}"
            :auto-load="false"
          >

            <template v-slot:cmd="text,current">
              <span @dblclick="seek(current)">{{ current.cmd }}</span>
            </template>
            <template v-slot:action="text,current">

            </template>

          </s-table>
        </a-card>
      </a-layout-sider>
    </a-layout>
  </a-modal>
</template>

<style scoped lang="less">
.ant-layout-sider{
  background-color: #f0f2f5;
}
</style>
