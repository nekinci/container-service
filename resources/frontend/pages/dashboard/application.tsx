import {Button, Typography} from "@mui/material";
import React from "react";
import DashboardMain from "../../src/components/DashboardMain/DashboardMain";
import {useApp} from "../../src/hooks/useApp";
import {getEnvironment} from "../../environment/environment";
import http from "../../src/client/http";
import Moment from "react-moment";
import Head from "next/head";

interface AppInfo {
    containerId: string;
    environments: any;
    image: string;
    name: string;
    owner: string;
    startTime: Date;
    status: any;
    url: string;
}

export default function Application() {

    const [currentApp] = useApp();
    const [appInfo, setAppInfo] = React.useState<AppInfo>(null);

    React.useEffect(() => {
        if (currentApp){
            http.get(getEnvironment().rootUrl + 'info/' + currentApp)
                .then((res) => {
                    if (res.data){
                        const app: AppInfo = res.data;
                        setAppInfo(app);
                    }
                })
        }
    }, [currentApp]);

    return (
        <>
            <Head>
                <title>Application Informations</title>
            </Head>
            <DashboardMain>

                <div style={{display: 'flex', padding: '20px', gap: '20px'}}>
                    <Typography variant={'subtitle1'} color={'secondary'}>
                        Application Information:
                    </Typography>
                    <div id={'applicationInformations'}>
                        {appInfo && (
                            <>
                                <Typography variant={'subtitle1'} color={'secondary'}>
                                    Container Id: {appInfo.containerId}
                                </Typography>
                                <Typography variant={'subtitle1'} color={'secondary'}>
                                    Image: {appInfo.image}
                                </Typography>
                                <Typography variant={'subtitle1'} color={'secondary'}>
                                    Name: {appInfo.name}
                                </Typography>
                                <Typography variant={'subtitle1'} color={'secondary'}>
                                    Owner: {appInfo.owner}
                                </Typography>

                                <Typography variant={'subtitle1'} color={'secondary'}>
                                    Start Time: <Moment>{appInfo.startTime}</Moment>
                                </Typography>
                                <Typography variant={'subtitle1'} color={'secondary'}>
                                    Status: {appInfo.status}
                                </Typography>
                                <Typography variant={'subtitle1'} color={'secondary'}>
                                    Url: <Button target={'_blank'} component={'a'} href={appInfo.url} style={{textTransform: 'none'}}>{appInfo.url}</Button>
                                </Typography>

                            </>
                        )}
                    </div>
                </div>
            </DashboardMain>
        </>
    )
}
