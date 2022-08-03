import Room from "@/models/Room";

export default interface BookingDefaultDto{
    room: Room|null
    start: Date
    end: Date
}
