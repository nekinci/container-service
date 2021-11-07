import {EnvironmentDev} from "./environment-dev";

export interface Environment {
    rootUrl: string;
    env: string;
    snackbarHideDuration: number;
    wsUrl: string;
}

export function getEnvironment(): Environment{
    return new EnvironmentDev();
}