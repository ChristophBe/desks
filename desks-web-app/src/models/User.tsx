import AbstractModel from "./AbstractModel";

export default interface User extends AbstractModel{
    username:string
    firstname:string
    lastname:string,
}