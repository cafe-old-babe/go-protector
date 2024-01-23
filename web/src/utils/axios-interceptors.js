// import Cookie from 'js-cookie'
import store from '@/store'
// 401拦截
const resp401 = {
  /**
   * 响应数据之前做点什么
   * @param response 响应对象
   * @param options 应用配置 包含: {router, i18n, store, message}
   * @returns {*}
   */
  onFulfilled(response, options) {
    const {message} = options
    if (response.status === 401) {
      message.error('无此权限')
    }
    return response
  },
  /**
   * 响应出错时执行
   * @param error 错误对象
   * @param options 应用配置 包含: {router, i18n, store, message}
   * @returns {Promise<never>}
   */
  onRejected(error, options) {
    const {message} = options
    const {response} = error
    if (response.status === 401) {
      message.error('无此权限')
    }
    return Promise.reject(error)
  }
}

const resp403 = {
  onFulfilled(response, options) {
    const {message} = options
    if (response.status === 403) {
      message.error('请求被拒绝')
    }
    return response
  },
  onRejected(error, options) {
    const {message} = options
    const {response} = error
    if (response.status === 403) {
      message.error('请求被拒绝')
    }
    return Promise.reject(error)
  }
}

const resp200 = {
  onFulfilled(response, options) {
    const {message} = options
    if (response.status !== 200) {
      message.error('请求失败')
      return response
    }

    if (response.data.code === 403 || response.data.code === 401) {
      let {router,message} = options;
      router.push('/login').then(() => {
        message.warn("登录信息已失效,请重新登录", 3)
      })
      return
    }
    let token = store.getters["account/token"];
    if (response.headers.authorization && token !== response.headers.authorization ) {
      this.setToken(response.headers.authorization)
    }
    return response
  },
  onRejected(error, options) {
    const {message} = options
    const {response} = error
    // 统一包装失败信息
    if (response.status !== 200) {
      message.error('failure request: ' + response.config.url + ', status: ' + response.status +', message: '+ response.statusText)
    }
    return Promise.reject(error)
  }
}

const reqCommon = {
  /**
   * 发送请求之前做些什么
   * @param config axios config
   * @param options 应用配置 包含: {router, i18n, store, message}
   * @returns {*}
   */
  onFulfilled(config) {
    let token = store.getters["account/token"];
    if (token) {
      // 获取 token 加入的header中
      config.headers['Authorization'] = 'Bearer ' + token
    }

    /*
    const {message} = options;
    const {url, xsrfCookieName} = config

    if (url.indexOf('login') === -1 && xsrfCookieName && !Cookie.get(xsrfCookieName)) {
      message.warning('认证 token 已过期，请重新登录')
    }*/
    return config
  },
  /**
   * 请求出错时做点什么
   * @param error 错误对象
   * @param options 应用配置 包含: {router, i18n, store, message}
   * @returns {Promise<never>}
   */
  onRejected(error, options) {
    const {message} = options
    message.error(error.message)
    return Promise.reject(error)
  }
}

export default {
  request: [reqCommon], // 请求拦截
  response: [resp401, resp403, resp200] // 响应拦截
}
