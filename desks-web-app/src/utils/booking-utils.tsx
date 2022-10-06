import Booking from "@/models/Booking";
import moment, {Moment} from "moment";

export default class BookingUtils {
    static findOverlaps(bookings: Booking[], start: Moment, end: Moment): Booking[] {

        return bookings.filter((booking) => {

            const bookingStart = moment(booking.start)
            const bookingEnd = moment(booking.end)

            return (start.isSameOrBefore(bookingStart) && end.isSameOrAfter(bookingEnd))
                || start.isBetween(bookingStart, bookingEnd)
                || end.isBetween(bookingStart, bookingEnd)
        })
    }

    static roundToNextFiveMinutes(): moment.Moment{
        return this.roundToNextMinutes(5)
    }
    static roundToNextMinutes(nMinutes: number): moment.Moment{
        const now = moment().startOf("hour")
        const min = (nMinutes * Math.floor(moment().diff(now, "minute") / nMinutes)) + nMinutes
        return now.add(min, "minutes")
    }
    static roundToNextMinute(): moment.Moment{
       return this.roundToNextMinutes(1)
    }

}
