import React from 'react';
import {Button, Card, CardActions, CardContent, Container, Typography} from '@mui/material';
import event from "../../util/Event";
import {useRouter} from "next/router";
import {AuthUtil} from "../../util/AuthUtil";
import http from "../../client/http";
import {getEnvironment} from "../../../environment/environment";
import {UserInformation} from "../modals/Login/Login";

export function TryIt() {

    const router = useRouter();
    const [isLoggedIn, setIsLoggedIn] = React.useState(false);
    const listen = () => {
        const info = AuthUtil.getInformation();

        if (info != null){
            // @ts-ignore
            if (new Date().getTime() >= new Date(info?.expiresAt).getTime()){
                http.post(getEnvironment().rootUrl + 'refreshToken', {refresh_token: info.refreshToken})
                    .then((res) => {
                        const data = {...res.data} as any;
                        const newInfo = {...info, refreshToken: data.refresh_token, token: data.token, expiresAt: data.expires_at} as UserInformation;
                        AuthUtil.removeInformation();
                        AuthUtil.setInformation(newInfo)
                        setIsLoggedIn(true);
                    });
            } else {
                setIsLoggedIn(true);
            }
        }
    }

    React.useEffect(() => {
        listen()
    })

    React.useEffect(() => {
        event.on('loggedIn', (from) => {
            if (from === 'TryIt'){
                router.push('/dashboard/application?from=TryIt');
            }
        })
    })

    return (
        <React.Fragment>
            <Typography align={'center'} pb={'10px'} pt={'10px'} variant={'h4'} fontWeight={'bold'}>Let's try together!</Typography>

            <Container style={{display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '15px'}}>
                <Typography variant={'subtitle1'} color={'secondary'}>
                    Just touch upload button
                </Typography>
                <Card style={{minWidth: '500px', margin: '0 auto', padding: '20px 10px', boxSizing: 'border-box'}}>
                    <CardContent>
                        <Typography>
                            version: 1
                        </Typography>
                        <Typography>
                            name: nginxapp
                        </Typography>
                        <Typography>
                            image: nginx
                        </Typography>
                        <Typography>
                            port: 80
                        </Typography>
                        <Typography>
                            type: docker
                        </Typography>
                    </CardContent>
                    <CardActions style={{display: 'flex', justifyContent: 'center'}}>
                        <Button onClick={() => {
                            if (isLoggedIn){
                                router.push('/dashboard/application?from=TryIt')
                            } else {
                                event.emit('login', 'TryIt');
                            }
                        }}>Upload</Button>
                    </CardActions>
                </Card>
            </Container>
        </React.Fragment>
    )
}
