import {Environment} from "./environment";

export class EnvironmentDocker implements Environment{
    env: string;
    rootUrl: string;
    snackbarHideDuration: number;
    wsUrl: string;

    constructor() {
        this.env = 'DOCKER';
        this.rootUrl = 'http://api.host.docker.internal/';
        this.snackbarHideDuration = 5000;
        this.wsUrl = "ws://api.host.docker.internal/"
    }


}