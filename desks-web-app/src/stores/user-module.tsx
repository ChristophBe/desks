import {ActionContext, ActionTree, Module, MutationTree} from "vuex";
import User from "../models/User";
import {RootState} from "@/stores/store";
import {getData} from "@/utils/request-utils";


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

        const response = await getData("/api/v1.0/users/me")
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

