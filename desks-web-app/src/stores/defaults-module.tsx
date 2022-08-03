import {ActionContext, ActionTree, Module, MutationTree} from "vuex";
import {RootState} from "@/stores/store";
import FrontendConfiguration from "@/dtos/FrontendConfigurationDto";
import {getData} from "@/utils/request-utils";
import BookingDefaultDto from "@/dtos/BookingDefaultDto";
import {notify} from "@/stores/notification-module";

export interface DefaultsState {
    bookingDefaultsLoading: boolean
    bookingDefaultsFetched: boolean
    bookingDefaults: BookingDefaultDto | null
}

const state: DefaultsState = {
    bookingDefaultsLoading: false,
    bookingDefaultsFetched: false,
    bookingDefaults: null

}

const mutations: MutationTree<DefaultsState> = {
    bookingDefaultsLoading(state: DefaultsState) {
        state.bookingDefaultsLoading = true
    },
    setBookingDefaults(state: DefaultsState, defaults: BookingDefaultDto) {
        state.bookingDefaultsLoading = false
        state.bookingDefaultsFetched = true
        state.bookingDefaults = defaults
    }
}

const actions: ActionTree<DefaultsState, RootState> = {

    async fetchBookingDefaults({state,commit, rootState}: ActionContext<DefaultsState, RootState>) {
        commit("bookingDefaultsLoading");

        if (!rootState.user.user) {
            console.log("Failed to create booking defaluts because of missing user")
            return
        }
        const response = await getData(`/api/v1.0/users/${rootState.user.user.id}/bookings/defaults`)
        if (response.status >= 400) {
            console.log("Failed to fetch configuration", response.status)

            commit('setBookingDefaults', state.bookingDefaults)
            return
        }
        try {
            const bookingDefaults = await response.json()
            commit('setBookingDefaults', bookingDefaults)
        } catch (e) {
            console.log("Failed to fetch configuration", e)
        }

    },
}
export const defaultsModule: Module<DefaultsState, RootState> = {
    namespaced: true,
    state,
    mutations,
    actions
}

