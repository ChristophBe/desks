import Vuex, {ActionContext, Store, useStore as baseUseStore} from "vuex";

import {InjectionKey} from "vue";
import User from "../models/User";
import Booking from "@/models/Booking";
import Room from "@/models/Room";
import moment from "moment";


export interface State {
    user: User|null
    loading: boolean
    bookings: Array<Booking>
    rooms: Array<Room>
}

// define injection key
export const key: InjectionKey<Store<State>> = Symbol()

export const store = new Vuex.Store<State>({
    state (){
        return {
            user: null,
            loading: false,
            rooms: [],
            bookings: []
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
        },
        setRooms(state:State, rooms: Array<Room>){
            state.rooms = rooms
        },
        setBookings(state:State, bookings: Array<Booking>){
            state.bookings = bookings
        }
    },
    getters:{
        upcomingBookings(state:State) {
            return state.bookings.filter(booking => moment(booking.end).isAfter(moment.now()))
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

        },
        async fetchRooms({commit, state}: ActionContext<State, State>) {

            if(state.rooms.length > 0){
                return
            }
            const response = await fetch("/api/v1.0/rooms" )
            if(response.status >= 400 ){
                console.log("Failed to fetch rooms", response.status )
                return
            }
            try {
                const rooms = await response.json()
                commit('setRooms', rooms)
            }
            catch (e){
               console.log("Failed to fetch rooms", e)
            }

        },
        async fetchBookings({commit, state}: ActionContext<State, State>) {

            if(!state.user){
                return
            }
            const response = await fetch(`/api/v1.0/users/${state.user.id}/bookings` )
            if(response.status >= 400 ){
                console.log("Failed to fetch bookinges", response.status )
                return
            }
            try {
                const bookings = await response.json()
                commit('setBookings', bookings)
            }
            catch (e){
               console.log("Failed to fetch bookings", e)
            }

        },
        async bookDesk({dispatch, state}: ActionContext<State, State>, booking:Booking) {

            const response = await postData("/api/v1.0/bookings",booking )
            if(response.status >= 400 ){
                console.log("Failed to create booking", response.status )
                return
            }
            try {
                await response.json()
                await dispatch("fetchBookings")
            }
            catch (e){
               console.log("Failed to create booking", e)
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