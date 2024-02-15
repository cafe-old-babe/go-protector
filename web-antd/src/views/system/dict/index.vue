<template>
  <a-row :gutter="[8,8]" type="flex" :style="{ margin: 0 }">
    <a-col :span="12">
      <TypeList @changeTypeCode="changeTypeCode"/>
    </a-col>
    <a-col :span="12">
      <DataList :type-code="typeCode"/>
    </a-col>
  </a-row>
</template>

<script>
import DataList from './data'
import TypeList from './type'
export default {
  name: 'Dict',
  components: { TypeList, DataList },
  data() {
    return {
      typeCode: ''
    }
  },
  methods: {
    changeTypeCode: function (typeCode) {
      this.typeCode = typeCode
    }
  }
}
</script>
<style scoped >
.space-align-block .mock-block {
  display: inline-block;
  padding: 32px 8px 16px;
  background: rgba(150, 150, 150, 0.2);
}
</style>
