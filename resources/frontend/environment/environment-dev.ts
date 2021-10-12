import {Environment} from "./environment";

export class EnvironmentDev implements Environment{
    env: string;
    rootUrl: string;
    snackbarHideDuration: number;

    constructor() {
        this.env = 'DEV';
        this.rootUrl = 'http://localhost:8070/';
        this.snackbarHideDuration = 5000;
    }

}