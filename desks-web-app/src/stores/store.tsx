import Vuex, {Store} from "vuex";

import {InjectionKey} from "vue";
import {UserState, userModule} from "@/stores/user-module";
import {roomsModule} from "@/stores/rooms-module";
import {bookingsModule} from "@/stores/booking-module";
import {configurationModule} from "@/stores/configuration-module";


export interface RootState {
    user: UserState
}

// define injection key
export const key: InjectionKey<Store<RootState>> = Symbol()

export const store = new Vuex.Store<RootState>({
        modules: {
            user: userModule,
            rooms: roomsModule,
            bookings: bookingsModule,
            configuration: configurationModule
        }
    }
);