import {ActionContext, ActionTree, GetterTree, Module, MutationTree} from "vuex";
import Booking from "@/models/Booking";
import moment from "moment";
import {RootState} from "@/stores/store";


export interface BookingsState {
    loading: boolean
    bookings: Array<Booking>
}

const state: BookingsState = {
    loading: false,
    bookings: []
}

const mutations: MutationTree<BookingsState> = {
    loading(state: BookingsState) {
        state.loading = true
    },
    setBookings(state: BookingsState, bookings: Array<Booking>) {
        state.bookings = bookings
        state.loading = false
    }
}
const getters: GetterTree<BookingsState, RootState> = {
    upcomingBookings(state: BookingsState) {
        return state.bookings.filter(booking => moment(booking.end).isAfter(moment.now()))
    },
    todaysBookings(state: BookingsState) {
        return state.bookings.filter(booking => moment(booking.start).diff(moment.now(), 'days') === 0)
    }
}
const actions: ActionTree<BookingsState, RootState> = {

    async fetchBookings({commit, rootState}: ActionContext<BookingsState, RootState>) {

        if (!rootState.user.user) {
            return
        }
        const response = await fetch(`/api/v1.0/users/${rootState.user.user.id}/bookings`)
        if (response.status >= 400) {
            console.log("Failed to fetch bookinges", response.status)
            return
        }
        try {
            const bookings = await response.json()
            commit('setBookings', bookings)
        } catch (e) {
            console.log("Failed to fetch bookings", e)
        }

    },
    async bookDesk({dispatch, state, rootState}: ActionContext<BookingsState, RootState>, booking: Booking) {
        if (!rootState.user.user) {
            console.log("Failed to create booking because of missing user")
            return
        }
        booking.user = rootState.user.user
        const response = await postData("/api/v1.0/bookings", booking)
        if (response.status >= 400) {
            console.log("Failed to create booking", response.status)
            return
        }
        try {
            await response.json()
            await dispatch("fetchBookings")
        } catch (e) {
            console.log("Failed to create booking", e)
        }

    }
}
export const bookingsModule: Module<BookingsState, RootState> = {
    namespaced: true,
    state,
    mutations,
    getters,
    actions
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
