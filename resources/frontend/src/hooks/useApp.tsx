import React from "react";
import {AppUtil} from "../util/AppUtil";
import {ApplicationContext} from "../../pages/_app";

export function useApp(){

    const [currentApp, setCurrentApp] = React.useContext(ApplicationContext);
    const [isThereAnyApp, setIsThereAnyApp] = React.useState<any>(false);
    React.useEffect(() => {
        setCurrentApp(AppUtil.GetCurrentApp());
    }, []);

    React.useEffect(() => {
        if (currentApp != null){
            setIsThereAnyApp(true);
        }
    }, [currentApp]);

    const changeCurrentApp = (name) => {
        if (name == null) {
            AppUtil.DeleteCurrentApp();
            setCurrentApp(null);
        } else {
            AppUtil.SetCurrentApp(name);
            setCurrentApp(name);
        }
    }

    return [
        currentApp,
        isThereAnyApp,
        changeCurrentApp
    ];
}