import AbstractModel from "./AbstractModel";

export default interface User extends AbstractModel {
    givenName: string
    familyName: string,
}