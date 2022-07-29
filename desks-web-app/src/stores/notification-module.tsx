import {ActionContext, ActionTree, Module, MutationTree} from "vuex";
import User from "../models/User";
import {RootState} from "@/stores/store";


export interface UserState {
    user: User | null
    loading: boolean
}

const userState: UserState = {
    user: null,
    loading: false,
}

const userMutations: MutationTree<UserState> = {
    login(state: UserState, user: User) {
        state.user = user
        state.loading = false
    },
    logout(state: UserState) {
        state.user = null
        state.loading = false
    },
    loading(state: UserState) {
        state.loading = true
    },
};

const userActions: ActionTree<UserState, RootState> = {
    async fetchCurrentUser({commit}: ActionContext<UserState, RootState>) {
        commit('loading')

        const response = await fetch("/api/v1.0/users/me")
        if (response.status >= 400) {
            commit('logout')

            window.location.replace("/auth/login");
            return
        }
        try {
            const user = await response.json()
            commit('login', user)
        } catch (e) {
            commit('logout')
        }

    },
}
export const userModule: Module<UserState, RootState> = {
    namespaced: true,
    state: userState,
    mutations: userMutations,
    actions: userActions
}


// Example POST method implementation:
async function postData(url = '', data = {}) {
    // Default options are marked with *
    return await fetch(url, {
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