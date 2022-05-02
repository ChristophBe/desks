import Vuex, {ActionContext, ActionTree, Module, MutationTree, Store, useStore as baseUseStore} from "vuex";

import {InjectionKey} from "vue";
import User from "../models/User";
import Booking from "@/models/Booking";
import Room from "@/models/Room";
import moment from "moment";
import {RootState} from "@/stores/store";


export interface RoomsState {
    loading: boolean
    rooms: Array<Room>
}

const state: RoomsState = {
    loading: false,
    rooms: [],
}

const mutations: MutationTree<RoomsState> = {
    loading(state: RoomsState) {
        state.loading = true
    },
    setRooms(state: RoomsState, rooms: Array<Room>) {
        state.rooms = rooms
        state.loading = false
    },
}
const actions: ActionTree<RoomsState, RootState> = {

    async fetchRooms({commit, state}: ActionContext<RoomsState, RootState>) {

        if (state.rooms.length > 0) {
            return
        }
        commit("loading")

        const response = await fetch("/api/v1.0/rooms")
        if (response.status >= 400) {
            console.log("Failed to fetch rooms", response.status)
            return
        }
        try {
            const rooms = await response.json()
            commit('setRooms', rooms)
        } catch (e) {
            console.log("Failed to fetch rooms", e)
        }

    }

}
export const roomsModule: Module<RoomsState, RootState> = {
    namespaced: true,
    state,
    actions,
    mutations
}