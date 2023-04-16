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
interface RoomAndRange{
    roomId: Room["id"]
    from: MomentInput
    to: MomentInput
}

type LoadingState = "loading" | "loaded";

class Range{

    public readonly from: Moment
    public readonly to: Moment

    constructor(from: MomentInput, to: MomentInput){
        if (moment(from).isAfter(to)){
            throw "invalid range"
        }
        this.from = moment(from)
        this.to = moment(to)
    }

    public isInRange(t: MomentInput):boolean {
        const other =  moment(t)
        const result = other.isBetween(this.from,this.to, "seconds","[]")
        return result
    }
    public isSubrangeOrEqual(other :Range): boolean{
        return this.isInRange(other.from) && this.isInRange(other.to)
    }
}
class BookingRangeLoadingState extends Range{

    public readonly state: LoadingState;
    public readonly roomId: Room['id'];

    constructor(from: MomentInput, to: MomentInput, roomId: Room['id'], state: LoadingState){
        super(from, to)
        this.roomId = roomId
        this.state = state
    }

    public isState(state: LoadingState): boolean{
        return this.state === state
    }

    public identifier() {
        return `${this.roomId}-${this.from.toISOString()}-${this.to.toISOString()}`;
    }
}

export interface BookingsState {
    loadingMyBookings: boolean
    myBookingsFetched: boolean
    bookings: Array<Booking>
    loadingStates:Map<string, BookingRangeLoadingState>
}

const state: BookingsState = {
    loadingMyBookings: false,
    myBookingsFetched: true,
    bookings: [],
    loadingStates: new Map<string,BookingRangeLoadingState>()
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
    setBookings(state: BookingsState, bookings: Array<Booking>) {
        state.bookings = mergeBooking(state.bookings, bookings)
    },
    setLoadingRange(state: BookingsState, loadingState: BookingRangeLoadingState){
        state.loadingStates = new Map<string, BookingRangeLoadingState>(state.loadingStates.set(loadingState.identifier(), loadingState))
    },
    removeBooking(state: BookingsState, idToRemove: Booking['id']) {
        state.bookings = state.bookings.filter(booking => booking.id !== idToRemove)
    }
}

function isLoadingBookingsByRoomAndRange(loadingStates: Array<BookingRangeLoadingState>, roomId: Room["id"], range: Range): boolean {
    return loadingStates
        .filter(r => r.roomId === roomId)
        .filter(r => r.isState("loading"))
        .some(r => r.isSubrangeOrEqual(range))
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
            if (!roomId || !date) {
                return false

            }
            const range = new Range(moment(date).startOf("day"), moment(date).endOf("day"))
            return isLoadingBookingsByRoomAndRange([...state.loadingStates.values()], roomId, range)
        }
    },

    isLoadingBookingsByRoomAndWeek(state: BookingsState) {
        return (roomId?: Room["id"], date?: MomentInput) => {
            if (!roomId || !date) {
                return false
            }
            const range = new Range(moment(date).startOf("week"), moment(date).endOf("week"))
            return isLoadingBookingsByRoomAndRange([...state.loadingStates.values()], roomId, range)
        }
    },
    isLoadingBookingsByRoomAndRange(state: BookingsState) {
        return (roomId?: Room["id"], from?: MomentInput, to?: MomentInput) => {
            if (!roomId || !from || !to) {
                return false
            }
            const range = new Range(from,to)
            return isLoadingBookingsByRoomAndRange([...state.loadingStates.values()], roomId, range)
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
    async fetchBookingsByRoomAndDate({dispatch}: ActionContext<BookingsState, RootState>, {
        roomId,
        date,
    }: RoomAndDate) {
        return dispatch("fetchBookingsForRange", {
            roomId: roomId,
            from: moment(date).startOf("day"),
            to: moment(date).endOf("day")
        })
    },
    async fetchBookingsByRoomAndWeek({dispatch}: ActionContext<BookingsState, RootState>, {
        roomId,
        date,
    }: RoomAndDate) {
        return dispatch("fetchBookingsForRange", {
            roomId: roomId,
            from: moment(date).startOf("week"),
            to: moment(date).endOf("week")
        })
    },

    async fetchBookingsForRange({commit, getters}: ActionContext<BookingsState, RootState>, {
        roomId,
        from, to,
    }: RoomAndRange) {

        if (!roomId || !from || !to || getters.isLoadingBookingsByRoomAndRange(roomId, from, to)) {
            return
        }

        const loadingState = new BookingRangeLoadingState(from,to,roomId, "loading")

        commit("setLoadingRange", loadingState)
        const fromString = moment(from).toISOString()
        const toString = moment(to).toISOString()


        console.log(`fetch bookings for room ${roomId} from ${fromString} to ${toString}`)
        const res = await fetch(`/api/v1.0/rooms/${roomId}/bookings?from=${fromString}&to=${toString}`)
        if (res.status >= 400) {
            console.log(`Failed to fetch bookings for room ${roomId} from ${fromString} to ${toString}`)
            return
        }
        try {
            const bookings: Array<Booking> = await res.json()

            console.log(`Received ${bookings.length} bookings for room ${roomId} from ${fromString} to ${toString}`)
            commit("setBookings", bookings)
            commit("setLoadingRange", new BookingRangeLoadingState(from,to,roomId, "loaded"))
        } catch (e) {
            console.log(`Failed to fetch bookings for room ${roomId} from ${from} to ${to}`, e)
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
