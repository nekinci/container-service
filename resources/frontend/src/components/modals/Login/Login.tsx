import { Close } from "@mui/icons-material";
import {
    Modal,
    Card,
    CardContent,
    Typography,
    Divider,
    TextField,
    Button,
    IconButton,
    Snackbar, Alert, Backdrop, CircularProgress
} from "@mui/material";
import React from "react";
import http from "../../../client/http";
import {getEnvironment} from "../../../../environment/environment";
import {AuthUtil} from "../../../util/AuthUtil";

export interface LoginModal {
    open: boolean;
    setOpen: any;
}

export interface UserInformation {
    email?: string;
    token?: string;
    expiresAt?: Date;
    refreshToken?: string;
}

export function Login({open, setOpen}: LoginModal){

    const [email, setEmail] = React.useState('')
    const [password, setPassword] = React.useState('')
    const [snackbarOpen, setSnackbarOpen] = React.useState(false)
    const [snackbarContent, setSnackbarContent] = React.useState(null)
    const [backdropOpen, setBackdropOpen] = React.useState(false);

    React.useEffect(() => {
        setEmail('')
        setPassword('')
    }, [open]);

    const onLogin = () => {
        setBackdropOpen(true);
        http.post(getEnvironment().rootUrl + "login", {email, password})
            .then((response) => {
                const data = response.data as any;
                const info = {token: data.token, refreshToken: data.refresh_token, expiresAt: data.expires_at} as UserInformation;
                info.email = email;
                AuthUtil.setInformation(info);
                setOpen(false)
                setBackdropOpen(false);
            }, (err) => {
                setBackdropOpen(false)
                setSnackbarOpen(true)
                // @ts-ignore
                setSnackbarContent(<Alert severity={'error'} onClose={() => setSnackbarOpen(false)}>
                    {err.response.data}
                </Alert>)
            });
    }

    const onRegister = () => {
        setBackdropOpen(false);
        http.post(getEnvironment().rootUrl + "register", {email, password})
            .then(() => {
                setBackdropOpen(false);
                onLogin();
            }, err => {
                setBackdropOpen(false);
                setSnackbarOpen(true)
                // @ts-ignore
                setSnackbarContent(<Alert severity={'error'} onClose={() => setSnackbarOpen(false)}>
                    {err.response.data}
                </Alert>)
            });
    }

    const onClose = () => {
        setOpen(false);
    };

    // @ts-ignore
    // @ts-ignore
    return (
       <React.Fragment>
           <Backdrop style={{color: 'white'}} open={backdropOpen}>
               <CircularProgress color="inherit" />
           </Backdrop>
           <Snackbar
               anchorOrigin={{vertical: 'top', horizontal: 'right'}}
               open={snackbarOpen}
               autoHideDuration={getEnvironment().snackbarHideDuration}
               onClose={() => setSnackbarOpen(false)}
           >
               {snackbarContent}
           </Snackbar>

           <Modal open={open} onClose={onClose}>
               <div style={{display: 'flex', justifyContent: 'center'}}>
                   <Card style={{'minWidth': '380px', margin: '60px', position: 'relative'}}>
                       <CardContent>
                           <Typography align={'center'} style={{padding: '5px'}} variant={'h6'}>
                               Login or Register
                           </Typography>
                           <IconButton onClick={() => setOpen(false)} style={{position: 'absolute', right: '5px', top: '15px'}}>
                               <Close/>
                           </IconButton>
                           <Divider style={{marginTop: '8px'}} />

                           <div style={{display: 'flex', flexDirection: 'column', padding: '45px 15px 20px 15px', gap: '15px'}}>
                               <TextField type={'email'} value={email} onChange={(e) => setEmail(e.target.value)} autoComplete={'off'} label={'E-mail'}/>
                               <TextField value={password} onChange={(e) => setPassword(e.target.value)} autoComplete={'off'} type={'password'} label={'Password'}/>
                           </div>
                           <div style={{margin: '0 auto', textAlign: 'center', display: 'flex', flexDirection: 'column', padding: '10px', gap: '10px'}}>
                               <Button onClick={onRegister} variant={'outlined'} color={'primary'}>Register and Login</Button>
                               <Button onClick={onLogin} color={'primary'} variant={'contained'}>Login</Button>
                           </div>
                       </CardContent>
                   </Card>
               </div>
           </Modal>
       </React.Fragment>
    )
}
