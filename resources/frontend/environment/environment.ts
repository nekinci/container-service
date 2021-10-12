import {EnvironmentDev} from "./environment-dev";

export interface Environment {
    rootUrl: string;
    env: string;
    snackbarHideDuration: number;
}

export function getEnvironment(): Environment{
    return new EnvironmentDev();
}