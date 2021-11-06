import {
    Button,
    Container,
    Divider,
    FormControl,
    IconButton,
    InputLabel,
    MenuItem,
    Select,
    Stack, Tooltip, tooltipClasses,
    TooltipProps
} from "@mui/material";
import React from "react";
import { Header } from "../Header/Header";
import {useRouter} from "next/router";
import {AuthUtil} from "../../util/AuthUtil";
import {useApp} from "../../hooks/useApp";
import {Add} from "@mui/icons-material";
import {getEnvironment} from "../../../environment/environment";
import http from "../../client/http";
import event from "../../util/Event";
import {styled} from "@mui/styles";
import NextLink, { LinkProps as NextLinkProps } from 'next/link';



export default function DashboardMain({children}: any){

    const router = useRouter();
    const [currentApp, isThereAnyApp, setCurrentApp] = useApp();
    const [appList, setAppList] = React.useState(null);

    React.useEffect(() => {
        const info = AuthUtil.getInformation();
        if (info == null){
            router.push('/').then()
            return
        }

        // @ts-ignore
        if (new Date().getTime() >= new Date(info?.expiresAt).getTime()){
            router.push("/").then()
        }

        const query = router.query;
        if (query?.from === 'TryIt'){
            event.emit('runApp', 'TryIt')
        }

    });

    React.useEffect(() => {

    }, []);

    React.useEffect(() => {
    })

    React.useEffect(() => {

        http.get(getEnvironment().rootUrl + 'myApps').then((data) => {
            if (currentApp !== null){
                if (!data.data?.includes(currentApp)){
                    setCurrentApp(null);
                }
            }
            setAppList(data.data);
        });
    }, [currentApp]);

    return (
        <div>
            <Header />
            <Container style={{padding: '45px, 5px'}}>
                <div style={{display: 'flex', justifyContent: 'space-between', padding: '25px 0', alignItems: 'end'}}>
                    <Stack direction={'row'} spacing={2}>
                        <NextLink href={'/dashboard/application'} passHref>
                            <Button component={'a'}>Application</Button>
                        </NextLink>
                        <NextLink href={'/dashboard/logs'} passHref>
                            <Button component={'a'}>Logs</Button>
                        </NextLink>
                        <NextLink href={'/dashboard/terminal'} passHref>
                            <Button component={'a'} href={'/dashboard/terminal'}>Terminal</Button>
                        </NextLink>
                    </Stack>
                    <div style={{display:'flex', justifyContent:'space-between', alignItems:'end', gap: '30px'}}>
                        <FormControl variant={'standard'} sx={{minWidth: '250px'}}>
                            <InputLabel id="currentAppSelector">Current Application</InputLabel>
                            <Select
                                placeholder={'Select an app'}
                                color={'secondary'}
                                label={'Select an app'}
                                labelId={'currentAppSelector'}
                                id={'currentAppSelect'}
                                value={currentApp === null ? '': currentApp}
                                onChange={(e) => setCurrentApp(e.target.value)}
                            >
                                {appList?.map((d) => {
                                    return (
                                        <MenuItem key={d} value={d}>
                                            {d}
                                        </MenuItem>
                                    )
                                })}
                            </Select>
                        </FormControl>
                        <Button onClick={() => event.emit('runApp')} variant={'contained'} startIcon={ <Add/>}  color={'primary'}>
                            Run application
                        </Button>
                    </div>
                </div>
                <Divider />
                {currentApp !== null ? children : (<div>Choose an application...</div>)}
            </Container>
        </div>
    );
}
