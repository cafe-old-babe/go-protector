// import {LOGIN, ROUTES} from '@/services/api'
import {request, METHOD, removeAuthorization} from '@/utils/request'
const url = {
  login: '/api/user/login',
  routes: '/api/system/routes'
}

/**
 * 登录服务
 * @param name 账户名
 * @param password 账户密码
 * @returns {Promise<AxiosResponse<T>>}
 */
export async function login({loginName, password, cid, code}) {
  return request(url.login, {
    loginName: loginName,
    password: password,
    cid: cid,
    code: code
  })
}

export async function getRoutesConfig() {
  return request(url.routes,{}, METHOD.GET)
}

/**
 * 退出登录
 */
export function logout() {
  localStorage.removeItem(process.env.VUE_APP_ROUTES_KEY)
  localStorage.removeItem(process.env.VUE_APP_PERMISSIONS_KEY)
  localStorage.removeItem(process.env.VUE_APP_ROLES_KEY)
  removeAuthorization()
}
export default {
  login,
  logout,
  getRoutesConfig
}
