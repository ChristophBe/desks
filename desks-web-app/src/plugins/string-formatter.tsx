import {App, Plugin} from 'vue';
import moment from "moment/moment";


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
}

export const stringFormat: Plugin = {
    install: (app: App) => {
        app.config.globalProperties.$format = new Formatter()
    }
}


