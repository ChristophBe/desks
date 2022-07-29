import {ActionContext, ActionTree, Dispatch, Module, MutationTree} from "vuex";
import User from "../models/User";
import {RootState} from "@/stores/store";


export interface Notification {
    notificationType: "normal"|"error"
    text: string

}
export interface NotificationState  extends Notification {
    time: number
}

const notificationState: NotificationState = {
    notificationType: "normal",
    text: "",
    time: Date.now()
}

const notificationMutations: MutationTree<NotificationState> = {
    notify(state: NotificationState, notification: Notification) {
        state.notificationType = notification.notificationType
        state.text = notification.text
        state.time = Date.now()
    },
};

const notificationActions: ActionTree<NotificationState, RootState> = {
    async notify({commit}: ActionContext<NotificationState, RootState>, notification: Notification ) {
        commit('notify',notification)


    },
}
export const notificationModule: Module<NotificationState, RootState> = {
    namespaced: true,
    state: notificationState,
    mutations: notificationMutations,
    actions: notificationActions
}


export async function notify(dispatch:Dispatch, text: string, type: Notification['notificationType']= "normal"){
    await dispatch("notification/notify", {
        text: text,
        notificationType:type
    }, {
        root:true
    });
}
