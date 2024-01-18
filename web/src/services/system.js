import {METHOD, request} from "@/utils/request";

const url = {
    captcha: "/api/system/captcha"
}

export async function getCaptcha(){
    return request(url.captcha,{}, METHOD.GET)
}

export default {

}
