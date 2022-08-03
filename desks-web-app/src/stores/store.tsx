import Vuex, {Store} from "vuex";

import {InjectionKey} from "vue";
import {UserState, userModule} from "@/stores/user-module";
import {roomsModule} from "@/stores/rooms-module";
import {bookingsModule, BookingsState} from "@/stores/booking-module";
import {configurationModule, ConfigurationState} from "@/stores/configuration-module";
import {notificationModule} from "@/stores/notification-module";
import {defaultsModule, DefaultsState} from "@/stores/defaults-module";


export interface RootState {
    user: UserState
    rooms: RootState
    bookings: BookingsState
    configuration: ConfigurationState
    notification: ConfigurationState
    defaults: DefaultsState
}

// define injection key
export const key: InjectionKey<Store<RootState>> = Symbol()

export const store = new Vuex.Store<RootState>({
        modules: {
            user: userModule,
            rooms: roomsModule,
            bookings: bookingsModule,
            configuration: configurationModule,
            notification: notificationModule,
            defaults: defaultsModule
        }
    }
);
