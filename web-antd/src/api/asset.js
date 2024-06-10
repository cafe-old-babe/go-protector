import request from '@/utils/request'

const api = {
  loadGroupTree: '/asset-group/tree',
  loadAsset: '/asset-info/page',
  loadAuthAsset: '/asset-info/auth/page',
  loadGatewayList: '/asset-gateway/list'
}
// 加载 资源组
export function loadGroupTree() {
  return request.post(api.loadGroupTree)
}

export function loadAsset(data) {
  return request.post(api.loadAsset, data)
}

export function loadGatewayList() {
  return request.post(api.loadGatewayList)
}

export function loadAuthAsset(data) {
  return request.post(api.loadAuthAsset, data)
}
