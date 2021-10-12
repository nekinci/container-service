import {Avatar, Backdrop, Button, CircularProgress, Container, Menu, MenuItem, Typography} from '@mui/material';
import React from 'react';
import {Login, UserInformation} from '../modals/Login/Login';
import {AuthUtil} from "../../util/AuthUtil";
import {isNullOrUndefined} from "util";
import http from "../../client/http";
import {getEnvironment} from "../../../environment/environment";
import {deepOrange} from "@mui/material/colors";
import {useRouter} from "next/router";

export function Header() {

    const [loginOpen, setLoginOpen] = React.useState(false);
    const [isLoggedIn, setIsLoggedIn] = React.useState(false);
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const [backdropOpen, setBackdropOpen] = React.useState(false);
    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(event.currentTarget);
    };
    const handleClose = () => setAnchorEl(null);
    const menuOpen = Boolean(anchorEl);

    React.useEffect(() => {
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
    });

    const router = useRouter()

    const logout = () => {

        setBackdropOpen(true);
        AuthUtil.removeInformation();
        setIsLoggedIn(false);
        setAnchorEl(null);
        setBackdropOpen(false);
        router.push("/").then()
    }

    return (
        <Container style={{padding: '10px 0'}}>
            <Backdrop style={{color: 'white'}} open={backdropOpen}>
                <CircularProgress color="inherit" />
            </Backdrop>
            <Login open={loginOpen} setOpen={setLoginOpen}/>
            <div style={{display:'flex', justifyContent: 'space-between'}}>
                <Typography variant={'h6'}>Container Cloud</Typography>
                <div style={{display: 'flex'}}>
                    <Button component={'a'} target={'_blank'} href={'https://github.com/nekinci'} variant={'text'} color={'secondary'}>
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                    </Button>
                    {!isLoggedIn && <Button onClick={() => setLoginOpen(true)} variant={'text'} color={'secondary'}>Login</Button>}
                    {isLoggedIn && (
                        <>
                            <Button onClick={handleClick}>
                                <Avatar sx={{width: 36, height: 36, bgcolor: deepOrange[500]}}>
                                    {AuthUtil.getInformation()?.email?.charAt(0).toUpperCase()}
                                </Avatar>
                            </Button>
                            <Menu
                                open={menuOpen}
                                anchorEl={anchorEl}
                                onClose={handleClose}
                                MenuListProps={{
                                    'aria-labelledby': 'basic-button',
                                }}
                        >
                                <MenuItem component={'a'} href={'/dashboard'} onClick={() => {}}>Dashboard</MenuItem>
                                <MenuItem onClick={logout}>Logout</MenuItem>
                        </Menu></>
                    )}
                </div>
            </div>
        </Container>
    )
}



