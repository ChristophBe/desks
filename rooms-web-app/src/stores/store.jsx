import Vuex from 'vuex';

export const store =  new Vuex.Store({
    state () {
        return {
            user: null,
            loading: false
        }
    },
    mutations: {
        login (state ,user) {
            state.user = user
            state.loading = false
        },
        logout (state) {
            state.user = null
            state.loading = false
        },
        loading(state){
            state.loading = true
        }
    },
    actions:{
        async login({commit}, user) {
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
        async fetchCurrentUser({commit}) {
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