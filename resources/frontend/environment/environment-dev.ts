import {Environment} from "./environment";

export class EnvironmentDev implements Environment{
    env: string;
    rootUrl: string;
    snackbarHideDuration: number;
    wsUrl: string;

    constructor() {
        this.env = 'DEV';
        this.rootUrl = 'http://api.localhost/';
        this.snackbarHideDuration = 5000;
        this.wsUrl = "ws://api.localhost/"
    }


}