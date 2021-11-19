import {EnvironmentDev} from "./environment-dev";
import {EnvironmentProd} from "./environment-prod";
import {EnvironmentDocker} from "./environment-docker";

export interface Environment {
    rootUrl: string;
    env: string;
    snackbarHideDuration: number;
    wsUrl: string;
}
const systemEnv = process.env.HOST_ENV || 'DEV';

export function getEnvironment(): Environment{
    switch (systemEnv){
        case 'DOCKER':
            return new EnvironmentDocker();
        case 'PROD':
            return new EnvironmentProd();
        case 'DEV':
            return new EnvironmentDev();
        default:
            return new EnvironmentDev();
    }
}