import request from '@/utils/request'

const api = {
  loadUser: '/user/page',
  loadDept: '/user/dept/tree'
}
// 加载 部门
export function loadDept() {
  return request.post(api.loadDept)
}

// 加载用户
export function loadUser(parameter) {
  return request.post(api.loadUser, parameter)
}
