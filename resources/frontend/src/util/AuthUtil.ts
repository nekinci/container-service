import {UserInformation} from "../components/modals/Login/Login";

export class AuthUtil {

    public static setInformation(info: UserInformation): void {
        localStorage.setItem('userInformation', JSON.stringify(info))
    }

    public static getInformation(): UserInformation | null {
        let info = null
        if (typeof window !== "undefined"){
            info = localStorage.getItem('userInformation')
        }
        if (info == null) {
            return null
        }
        return JSON.parse(info) as UserInformation
    }

    public static removeInformation(){
        localStorage.removeItem('userInformation')
    }
}
