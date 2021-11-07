import {EnvironmentDev} from "./environment-dev";
import {EnvironmentProd} from "./environment-prod";

export interface Environment {
    rootUrl: string;
    env: string;
    snackbarHideDuration: number;
    wsUrl: string;
}

export function getEnvironment(): Environment{
    return new EnvironmentProd();
}