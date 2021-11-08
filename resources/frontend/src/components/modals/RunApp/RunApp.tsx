import {
    Alert,
    Backdrop,
    Button,
    Card, CardActions,
    CardContent, CardHeader,
    CircularProgress,
    Divider,
    IconButton,
    Modal, Snackbar, Tooltip, tooltipClasses, TooltipProps,
    Typography
} from "@mui/material";
import React from "react";
import {Close, InfoOutlined, InfoRounded, RunCircle} from "@mui/icons-material";
import {styled} from "@mui/styles";
import http from "../../../client/http";
import {getEnvironment} from "../../../../environment/environment";
import {useApp} from "../../../hooks/useApp";
import {useRouter} from "next/router";
import GenericSnackbar from "../../Snackbar/GenericSnackbar";
import event from "../../../util/Event";
import {func} from "prop-types";
import {LineType} from "react-terminal-ui";
import App from "next/app";
import {AuthUtil} from "../../../util/AuthUtil";


type AppState = {
    Type: string,
    Content: string,
    ApplicationCount: number,
    Time: Date | string,
    MaxApplicationCount: number
}

const Input = styled('input')({
    display: 'none'
})

const HtmlTooltip = styled(({ className, ...props }: TooltipProps) => (
    <Tooltip {...props} classes={{ popper: className }} />
))(({ theme }) => ({
    [`& .${tooltipClasses.tooltip}`]: {
        backgroundColor: '#556cd6',
        color: 'white',
        maxWidth: 220,
        fontSize: theme.typography.pxToRem(12),
        border: '1px solid #dadde9',
    },
    [`& .${tooltipClasses.arrow}`]: {
        color: '#556cd6',
    },
}));

export function RunApp() {

    const [backdropOpen, setBackdropOpen] = React.useState(false);
    const [selectedFile, setSelectedFile] = React.useState(null);
    const [selectedFileName, setSelectedFileName] = React.useState('');
    const [currentApp, isThereAnyApp, setCurrentApp] = useApp();
    const [open, setOpen] = React.useState(false);
    const [from, setFrom] = React.useState(null);
    const [ws, setWS] = React.useState<WebSocket | any>(null);
    const [applicationCount, setApplicationCount] = React.useState(0);
    const [maxApplicationCount, setMaxApplicationCount] = React.useState(1);
    const router = useRouter();

    React.useEffect(() => {
        if (ws == null){
            if (AuthUtil.getInformation()?.token != null){
                setWS(new WebSocket(getEnvironment().wsUrl + "appState?token=" + AuthUtil.getInformation()?.token))
            }
        }
    })

    React.useEffect(() => {
        event.on('runApp', (fromq: string) => {
            setOpen(true);
            setFrom(fromq);
        })


        return () => {
            if (ws != null) {
                ws.close();
                setWS(null);
            }
        }
    }, []);

    React.useEffect(() => {
        if (ws != null) {
            ws.onopen = function (event) {
                ws.onmessage = function (message){
                    const event: AppState = JSON.parse(message.data) as AppState;
                    setApplicationCount(event.ApplicationCount);
                    setMaxApplicationCount(event.MaxApplicationCount);
                    console.log(event);
                }
            }
        }
    }, [ws]);

    React.useEffect(async () => {
        if (from === 'TryIt' ){
            let response = await fetch('/example.yml');
            let data = await response.blob();
            let file = new File([data], "example.yml");
            setSelectedFile(file);
        }
    }, [from]);

    const onClose = () => {
        setSelectedFileName('');
        setSelectedFile(null);
        setOpen(false);
        setFrom(null);
    }

    const downloadTemplate = () => {
        const a = document.createElement('a');
        a.href = '/example.yml';
        a.setAttribute('download', 'example.yml');
        a.click();
    }

    const chooseFile = (e) => {
        setSelectedFile(e.target.files[0]);
        setSelectedFileName(e.target.files[0].name)
    }

    const apply = () => {
        if (selectedFile === null){
            event.emit('snackbar', 'Please select a valid file.', 'warning')
            return
        }
        const formData = new FormData();
        formData.append('file', selectedFile);
        setBackdropOpen(true);
        setOpen(false);
        http.post(getEnvironment().rootUrl + 'run', formData).then((res) => {
            setBackdropOpen(false);
            setOpen(true);
            setCurrentApp(res.data.appName);
            onClose();
            router.push('/dashboard/application');
        }, (err) => {
            setBackdropOpen(false);
            setOpen(true);
            console.log(err.response)
            event.emit('snackbar', err.response.data.message, 'error')
        });
    }

    return (
        <React.Fragment>
            <Backdrop open={backdropOpen} style={{color: 'white'}}>
                <CircularProgress color="inherit" />
            </Backdrop>
            <Modal open={open} onClose={onClose}>
                <div style={{display: 'flex', justifyContent: 'center'}}>
                    <Card style={{'minWidth': '550px', margin: '60px', position: 'relative'}}>
                        <CardHeader title={'Run an application'} action={
                            <IconButton onClick={onClose} aria-label="settings">
                                <Close />
                            </IconButton>
                        }>
                        </CardHeader>
                        <CardContent>
                            <Divider />
                            <Alert  severity="info">
                                <Typography variant={'body2'}>A maximum of {maxApplicationCount} + 1 applications are allowed on the whole system. Current: {applicationCount}</Typography>
                            </Alert>
                           <div style={{padding:'150px 50px'}}>
                               <Typography component={'div'} variant={'body1'} color={'secondary'}>Choose a yaml file which contain application informations</Typography>
                               <Typography variant={'subtitle2'} align={'center'}>{selectedFileName}</Typography>
                               <div style={{display: 'flex', justifyContent:'center', alignItems: 'center', gap: '10px', marginTop: '10px'}}>
                                   <Button onClick={downloadTemplate} variant={'outlined'}>
                                       Download template
                                   </Button>
                                   <div>
                                       <label htmlFor={'upload-yml'}>
                                           <Input accept={'.yaml, .yml'} onChange={chooseFile} id={'upload-yml'} type={'file'} />
                                           <Button variant={'contained'} component={'span'}>Choose file</Button>
                                       </label>
                                   </div>

                               </div>
                           </div>
                        </CardContent>
                        <CardActions style={{display: 'flex', justifyContent: 'center'}}>
                            {from === 'TryIt' && (
                                <HtmlTooltip arrow={true} open={true} placement={'top'}
                                             title={
                                                 <React.Fragment>
                                                    <div style={{display: 'flex', 'gap': '5px', alignItems: 'center'}}>
                                                        <InfoOutlined />
                                                        <Typography fontWeight={'bold'} variant={'subtitle1'} fontSize={'12px'} color="inherit">Click apply for run application</Typography>
                                                    </div>
                                                 </React.Fragment>
                                             }
                                >
                                    <Button disabled={selectedFile === null} onClick={apply} variant={'outlined'} startIcon={<RunCircle />}>Apply</Button>
                                </HtmlTooltip>
                            )}
                            {from !== 'TryIt' && (
                                    <Button disabled={selectedFile === null} onClick={apply} variant={'outlined'} startIcon={<RunCircle />}>Apply</Button>
                            )}
                        </CardActions>
                    </Card>
                </div>
            </Modal>
        </React.Fragment>
    )
}