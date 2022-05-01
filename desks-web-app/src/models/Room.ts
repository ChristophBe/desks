import AbstractModel from "@/models/AbstractModel";

export default interface Room extends AbstractModel{
    name: string
    capacity: number
}