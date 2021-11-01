import {Environment} from "./environment";

export class EnvironmentDev implements Environment{
    env: string;
    rootUrl: string;
    snackbarHideDuration: number;

    constructor() {
        this.env = 'DEV';
        this.rootUrl = 'http://api.localhost:7888/';
        this.snackbarHideDuration = 5000;
    }

}