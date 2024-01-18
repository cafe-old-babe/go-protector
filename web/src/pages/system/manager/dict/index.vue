
<template>
  <a-row :gutter="[8,8]" type="flex" :style="{ margin: 0 }">
    <a-col :span="12" >
      <TypeList/>
    </a-col>
    <a-col :span="12">
      <DataList/>
    </a-col>
  </a-row>
</template>

<script setup>

import DataList from "@/pages/system/manager/dict/DataList";
import TypeList from "@/pages/system/manager/dict/TypeList";
import column from './column'
export default {
    name: "typeList",
    components: { TypeList,DataList},
    data () {
        return {
            column: {
                typeColumn : column.typeColumn,
                dataColumn : column.dataColumn,
            },

        }
    }
}
</script>
<style scoped lang="less">

</style>