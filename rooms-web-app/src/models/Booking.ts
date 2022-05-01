import AbstractModel from "@/models/AbstractModel";
import User from "@/models/User";
import Room from "@/models/Room";

export default interface Booking extends AbstractModel{
    start: Date;
    end: Date;
    user: User;
    room: Room;
}