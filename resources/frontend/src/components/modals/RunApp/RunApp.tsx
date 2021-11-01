import {
    Backdrop,
    Button,
    Card, CardActions,
    CardContent, CardHeader,
    CircularProgress,
    Divider,
    IconButton,
    Modal,
    Typography
} from "@mui/material";
import React from "react";
import {Close, RunCircle} from "@mui/icons-material";
import {styled} from "@mui/styles";
import http from "../../../client/http";
import {getEnvironment} from "../../../../environment/environment";
import {useApp} from "../../../hooks/useApp";
import {useRouter} from "next/router";


const Input = styled('input')({
    display: 'none'
})

export function RunApp({open, setOpen}: any) {

    const [backdropOpen, setBackdropOpen] = React.useState(false);
    const [selectedFile, setSelectedFile] = React.useState(null);
    const [selectedFileName, setSelectedFileName] = React.useState('');
    const [currentApp, isThereAnyApp, setCurrentApp] = useApp();
    const router = useRouter();

    const onClose = () => {
        setSelectedFileName('');
        setSelectedFile(null);
        setOpen(false)
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
        }, () => {
            setBackdropOpen(false);
            setOpen(true);
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
                           <div style={{padding:'150px 50px'}}>
                               <Typography component={'div'} variant={'body1'} color={'secondary'}>Choose a yaml file which contain application informations</Typography>
                               <Typography variant={'subtitle2'} align={'center'}>{selectedFileName}</Typography>
                               <div style={{display: 'flex', justifyContent:'center', alignItems: 'center', gap: '10px', marginTop: '10px'}}>
                                   <Button onClick={downloadTemplate} variant={'outlined'}>
                                       Download template
                                   </Button>
                                   <div>
                                       <label htmlFor={'upload-yml'}>
                                           <Input onChange={chooseFile} id={'upload-yml'} type={'file'} />
                                           <Button variant={'contained'} component={'span'}>Choose file</Button>
                                       </label>
                                   </div>

                               </div>
                           </div>
                        </CardContent>
                        <CardActions style={{display: 'flex', justifyContent: 'center'}}>
                            <Button onClick={apply} variant={'outlined'} startIcon={<RunCircle />}>Apply</Button>
                        </CardActions>
                    </Card>
                </div>
            </Modal>
        </React.Fragment>
    )
}