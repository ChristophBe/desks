import Vuex, {ActionContext, Store, useStore as baseUseStore} from "vuex";

import {InjectionKey} from "vue";
import User from "../models/User";


export interface State {
    user: User|null
    loading: boolean
}

// define injection key
export const key: InjectionKey<Store<State>> = Symbol()

export const store = new Vuex.Store<State>({
    state (){
        return {
            user: null,
            loading: false
        }
    },
    mutations: {
        login (state:State ,user:User) {
            state.user = user
            state.loading = false
        },
        logout (state:State) {
            state.user = null
            state.loading = false
        },
        loading(state:State){
            state.loading = true
        }
    },
    actions:{
        async login({commit}:ActionContext<State, State>, user:User) {
            commit('loading')

            const response = await postData("/api/v1.0/users/login",user)
            if (response.status >= 400) {
                commit('logout')
                return
            }
            try {
                const user = await response.json()
                commit('login', user)
            } catch (e) {
                commit('logout')
            }
        },
        async fetchCurrentUser({commit}: ActionContext<State, State>) {
            commit('loading')

            const response = await fetch("/api/v1.0/users/me")
            if(response.status >= 400 ){
                commit('logout')
                return
            }
            try {
                const user = await response.json()
                commit('login', user)
            }
            catch (e){
                commit('logout')
            }

        }
    }

})


// Example POST method implementation:
async function postData(url = '', data = {}) {
    // Default options are marked with *
    return  await fetch(url, {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, *same-origin, omit
        headers: {
            'Content-Type': 'application/json'
            // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: JSON.stringify(data) // body data type must match "Content-Type" header
    });
}

export function useStore () {
    return baseUseStore(key)
}