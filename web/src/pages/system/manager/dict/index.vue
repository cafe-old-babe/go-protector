
<template>
  <a-row :gutter="[8,8]" type="flex" :style="{ margin: 0 }">
    <a-col :span="12" >
      <TypeList @updateTypeCode="updateTypeCode"/>
    </a-col>
    <a-col :span="12">
      <DataList :type-code="typeCode"/>
    </a-col>
  </a-row>
</template>

<script>

import DataList from "@/pages/system/manager/dict/data/DataList";
import TypeList from "@/pages/system/manager/dict/type/TypeList";
export default {
  name: "typeList",
  components: {TypeList, DataList},
  data() {
    return {
      typeCode: '',
    }
  },
  methods: {
    updateTypeCode: function (typeCode) {
      this.typeCode = typeCode;
    }
  }
}
</script>
<style scoped lang="less">

</style>