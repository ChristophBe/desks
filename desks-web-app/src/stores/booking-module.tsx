import {ActionContext, ActionTree, GetterTree, Module, MutationTree} from "vuex";
import Booking from "@/models/Booking";
import moment, {Moment} from "moment";
import {RootState} from "@/stores/store";
import BookingUtils from "@/utils/booking-utils";
import Room from "@/models/Room";
import {notify} from "@/stores/notification-module";
import {deleteData, getData, patchData, postData} from "@/utils/request-utils";
import {MomentInput} from "moment/moment";

interface RoomAndDate {
    roomId: Room["id"]
    date: MomentInput
}

interface SetBookingsPayload extends RoomAndDate {
    bookings: Array<Booking>
}

export interface BookingsState {
    loadingMyBookings: boolean
    myBookingsFetched: boolean
    bookings: Array<Booking>
    loadingByRoomAndDate: Map<string, boolean>
}

const state: BookingsState = {
    loadingMyBookings: false,
    myBookingsFetched: true,
    bookings: [],
    loadingByRoomAndDate: new Map<string, boolean>()
}

function mergeBooking(bookings: Array<Booking>, bookingsToUpdate: Array<Booking>): Array<Booking> {
    const idsOfUpdatedBookings = bookingsToUpdate ? bookingsToUpdate.map(booking => booking.id) : []
    return [
        ...bookings.filter(booking => !idsOfUpdatedBookings.includes(booking.id)),
        ...bookingsToUpdate
    ];
}

const mutations: MutationTree<BookingsState> = {
    loadingMyBookings(state: BookingsState) {
        state.loadingMyBookings = true
    },
    setMyBookings(state: BookingsState, bookings: Array<Booking>) {
        state.bookings = mergeBooking(state.bookings, bookings)
        state.loadingMyBookings = false
        state.myBookingsFetched = true
    },
    setBookings(state: BookingsState, {bookings, roomId, date}: SetBookingsPayload) {
        state.bookings = mergeBooking(state.bookings, bookings)

        state.loadingByRoomAndDate.set(`${roomId}-${moment(date).format('YYYY-MM-DD')}`, false)

    },
    loadingForRoomAndDate(state: BookingsState, {roomId, date}: RoomAndDate) {
        state.loadingByRoomAndDate.set(`${roomId}-${moment(date).format('YYYY-MM-DD')}`, true)
    },
    removeBooking(state: BookingsState, idToRemove: Booking['id']) {
        state.bookings = state.bookings.filter(booking => booking.id !== idToRemove)
    }
}

const getters: GetterTree<BookingsState, RootState> = {

    getBookingsByRoomAndDay(state: BookingsState) {
        return (room: Room["id"], date: MomentInput): Booking[] => {
            return state.bookings
                .filter(booking => booking.room.id === room)
                .filter((booking: Booking) => moment(booking.start).isSame(date, 'days'))
        }
    },
    getByRoomAndDay(state: BookingsState) {
        return (room: Room["id"], date: MomentInput) => {
            return state.bookings
                .filter(booking => booking.room.id === room)
                .filter((booking: Booking) => moment(booking.start).isSame(date, 'days'))
        }
    },
    isLoadingBookingsByRoomAndDay(state: BookingsState) {
        return (roomId?: Room["id"], date?: MomentInput) => {
            if(!roomId && !date){
                return false
            }
            const key = `${roomId}-${moment(date).format('YYYY-MM-DD')}`
            return state.loadingByRoomAndDate.has(key) && state.loadingByRoomAndDate.get(key)
        }
    },
    myBookings(state: BookingsState, _, rootState: RootState): Array<Booking> {
        if (!rootState.user.user) {
            return []
        }
        return state.bookings.filter(booking => booking.user.id === rootState.user.user?.id)
    },
    myUpcomingBookings(state: BookingsState, getters) {
        return getters.myBookings.filter((booking: Booking) => moment(booking.end).isAfter(moment.now()))
    },
    myBookingsOfTheDay(state: BookingsState, getters) {
        return getters.myBookings.filter((booking: Booking) => moment(booking.start).isSame(moment.now(), 'days'))
    },
    getMyOverlappingBookings(state: BookingsState, getters) {
        if (!state.myBookingsFetched) {
            throw "bookings were not fetched"
        }
        return (room: Room, start: Moment, end: Moment, skipBookingIds: number[] = []) => {
            return room ? BookingUtils.findOverlaps(getters.myBookings, start, end)
                .filter(booking => !skipBookingIds.includes(booking.id))
                .filter(booking => booking.room.id === room.id) : false;
        }
    }
}
const actions: ActionTree<BookingsState, RootState> = {

    async fetchMyBookings({commit, rootState}: ActionContext<BookingsState, RootState>) {
        commit('loadingMyBookings')
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
            commit('setMyBookings', bookings)
        } catch (e) {
            console.log("Failed to fetch bookings", e)
        }

    },
    async fetchBookingsByRoomAndDate({commit, getters}: ActionContext<BookingsState, RootState>, {
        roomId,
        date
    }: RoomAndDate) {
        console.log("load bookings for room and day")
        if (!roomId || !date || getters.isLoadingBookingsByRoomAndDay(roomId, date)) {
            return
        }
        commit("loadingForRoomAndDate", {roomId: roomId, date: date})
        const dateString = moment(date).format('YYYY-MM-DD')
        const res = await fetch(`/api/v1.0/rooms/${roomId}/bookings?date=${dateString}`)
        if (res.status >= 400) {
            console.log(`Failed to fetch bookings for room ${roomId} on ${dateString} `)
            return
        }
        try {
            const bookings: Array<Booking> = await res.json()
            commit("setBookings", {bookings: bookings, roomId: roomId, date: date})
        } catch (e) {
            console.log(`Failed to fetch bookings for room ${roomId} on ${dateString} `, e)
        }
    },

    async deleteBooking({
                            commit,
                            dispatch,
                            state,
                            rootState
                        }: ActionContext<BookingsState, RootState>, booking: Booking) {

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
        await commit("removeBooking", booking.id)
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
            await dispatch("fetchMyBookings")
        } catch (e) {
            await notify(dispatch, "Failed to save the booking.", "error")
            console.log("Failed to create booking", e)
        }
        await notify(dispatch, "Saved the booking.")
    },
    leaveEarly: async function ({dispatch}: ActionContext<BookingsState, RootState>, booking: Partial<Booking>) {
        const updateBooking: Partial<Booking> = {
            id: booking.id,
            end: BookingUtils.roundToNextMinute().toDate()
        }
        await dispatch("saveBooking", updateBooking)
    }
}
export const bookingsModule: Module<BookingsState, RootState> = {
    namespaced: true,
    state,
    mutations,
    getters,
    actions
}
