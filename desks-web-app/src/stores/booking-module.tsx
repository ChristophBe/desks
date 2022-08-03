import {ActionContext, ActionTree, GetterTree, Module, MutationTree} from "vuex";
import Booking from "@/models/Booking";
import moment, {Moment} from "moment";
import {RootState} from "@/stores/store";
import BookingUtils from "@/utils/booking-utils";
import Room from "@/models/Room";
import {notify} from "@/stores/notification-module";
import {deleteData, getData, patchData, postData} from "@/utils/request-utils";


export interface BookingsState {
    loading: boolean
    bookingsFetched: boolean
    bookings: Array<Booking>
}

const state: BookingsState = {
    loading: false,
    bookingsFetched: true,
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
        return state.bookings.filter(booking => moment(booking.start).isSame(moment.now(), 'days'))
    },
    getOverlappingBookings(state: BookingsState) {
        if (!state.bookingsFetched) {
            throw "bookings were not fetched"
        }
        return (room: Room, start: Moment, end: Moment, skipBookingIds: number[] = []) => {
            return BookingUtils.findOverlaps(state.bookings, start, end)
                .filter(booking => !skipBookingIds.includes(booking.id))
                .filter(booking => booking.room.id === room.id);
        }
    }
}
const actions: ActionTree<BookingsState, RootState> = {

    async fetchBookings({commit, rootState}: ActionContext<BookingsState, RootState>) {

        if (!rootState.user.user) {
            return
        }
        const response = await getData(`/api/v1.0/users/${rootState.user.user.id}/bookings`)
        if (response.status >= 400) {
            console.log("Failed to fetch bookings", response.status)
            return
        }
        try {
            const bookings = await response.json()
            commit('setBookings', bookings)
        } catch (e) {
            console.log("Failed to fetch bookings", e)
        }

    },
    async fetchBookingDefaults({commit, rootState}: ActionContext<BookingsState, RootState>) {

        if (!rootState.user.user) {
            return
        }
        const response = await getData(`/api/v1.0/users/${rootState.user.user.id}/bookings`)
        if (response.status >= 400) {
            console.log("Failed to fetch bookings", response.status)
            return
        }
        try {
            const bookings = await response.json()
            commit('setBookings', bookings)
        } catch (e) {
            console.log("Failed to fetch bookings", e)
        }

    },

    async deleteBooking({dispatch, state, rootState}: ActionContext<BookingsState, RootState>, booking: Booking) {

        console.log("Booking to delete", booking)
        if (!rootState.user.user) {
            console.log("Failed to create booking because of missing user")

            await notify(dispatch, "Failed to delete the booking.", "error")
            return
        }
        booking.user = rootState.user.user

        const response = await deleteData("/api/v1.0/bookings/" + booking.id)
        if (response.status >= 400) {
            await notify(dispatch, "Failed to delete the booking.", "error")
            return
        }
        await dispatch("fetchBookings")
        await notify(dispatch, "Deleted the booking.")
    },
    saveBooking: async function ({
                                     dispatch,
                                     state,
                                     rootState
                                 }: ActionContext<BookingsState, RootState>, booking: Partial<Booking>) {
        if (!rootState.user.user) {
            console.log("Failed to create booking because of missing user")
            return
        }
        booking.user = rootState.user.user
        const response = booking.id
            ? await patchData("/api/v1.0/bookings/" + booking.id.toString(), booking)
            : await postData("/api/v1.0/bookings", booking);


        if (response.status >= 400) {
            await notify(dispatch, "Failed to save the booking.", "error")
            console.log("Failed to create booking", response.status)
            return
        }
        try {
            await response.json()
            await dispatch("fetchBookings")
        } catch (e) {
            await notify(dispatch, "Failed to save the booking.", "error")
            console.log("Failed to create booking", e)
        }
        await notify(dispatch, "Saved the booking.")
    }
}
export const bookingsModule: Module<BookingsState, RootState> = {
    namespaced: true,
    state,
    mutations,
    getters,
    actions
}
