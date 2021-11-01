import axios, {AxiosRequestConfig} from 'axios';
import {AuthUtil} from "../util/AuthUtil";

const skippedEndpoints = ['register', 'login', 'refreshToken'];

axios.interceptors.request.use(function (cfg: AxiosRequestConfig){
    let url = cfg.url
    if (url) {
        let splitted = url.split("/")
        if (splitted.length > 3){
            url = splitted[3]
        }
    }

    if (!skippedEndpoints.includes(url)){
        // @ts-ignore
        cfg.headers.Authorization = `Bearer ${AuthUtil.getInformation()?.token}`
    }

    return {
        ...cfg
    }
})

export default {
    get: axios.get,
    post: axios.post,
    put: axios.put,
    delete: axios.delete,
    patch: axios.patch
}
