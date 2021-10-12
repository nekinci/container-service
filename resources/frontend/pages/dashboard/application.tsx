import { Typography } from "@mui/material";
import React from "react";
import { DashboardMain } from "../../src/components/DashboardMain/DashboardMain";

export default function Application() {

    return (
        <DashboardMain>
            <div style={{display: 'flex', padding: '20px', gap: '20px'}}>
                <Typography variant={'subtitle1'} color={'secondary'}>
                    Application Information:
                </Typography>
                <div id={'applicationInformations'}>
                    <Typography variant={'subtitle1'} color={'secondary'}>
                        Name: {'AppName'}
                    </Typography>
                    <Typography variant={'subtitle1'} color={'secondary'}>
                        Name: {'AppName'}
                    </Typography>
                    <Typography variant={'subtitle1'} color={'secondary'}>
                        Name: {'AppName'}
                    </Typography>
                </div>
            </div>
        </DashboardMain>
    )
}
