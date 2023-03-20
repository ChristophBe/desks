import {App, Plugin} from 'vue';
import moment from "moment/moment";
import User from "@/models/User";


class Formatter {
    date(date: Date): string {
        return moment(date).format("DD.MM.YYYY")
    }

    time(date: Date): string {
        return moment(date).format("HH:mm")
    }

    timeRange(start: Date, end: Date): string {
        return `${this.time(start)} - ${this.time(end)}`
    }

    userFullName(user: User){
        return `${user.givenName} ${user.familyName}`
    }
    userInitials(user: User){
        return `${user.givenName.charAt(0)}${user.familyName.charAt(0)}  `
    }

}

export const stringFormat: Plugin = {
    install: (app: App) => {
        app.config.globalProperties.$format = new Formatter()
    }
}


