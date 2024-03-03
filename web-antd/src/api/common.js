
import request from '@/utils/request'

const api = {
  loadDictData: '/dict/dataList/'
}
// 加载 字典
export async function loadDictData(dictType) {
  if (dictType === '') {
    return []
  }
  let dictData
  try {
    const res = await request.post(api.loadDictData + dictType)
    const { code, data, message } = res
    if (code !== 200) {
      throw new Error(message)
    }
    dictData = data
  } catch (error) {
    console.error(error.message)
    dictData = []
  }

  return dictData
}
