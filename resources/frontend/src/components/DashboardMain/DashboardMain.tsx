import { Button, Container, Divider, Stack } from "@mui/material";
import React from "react";
import { Header } from "../Header/Header";
import {useRouter} from "next/router";
import {AuthUtil} from "../../util/AuthUtil";

export function DashboardMain({children}: any){

    const router = useRouter();

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


    });

    return (
        <div>
            <Header />
            <Container style={{padding: '45px, 5px'}}>
                <div style={{display: 'flex', justifyContent: 'center', padding: '25px'}}>
                    <Stack direction={'row'} spacing={2}>
                        <Button>Application</Button>
                        <Button>Logs</Button>
                        <Button>Terminal</Button>
                    </Stack>
                </div>
                <Divider />
                {children}
            </Container>
        </div>
    );
}
