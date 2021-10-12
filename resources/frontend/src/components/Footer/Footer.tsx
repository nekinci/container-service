import { Typography } from "@mui/material";
import React from "react";


export function Footer() {

    return (
        <div style={{display: 'flex', justifyContent: 'center', padding: '50px 4px', alignItems:'center', flexDirection: 'column'}}>
            <Typography variant={'h6'} component={'div'}>
                The project that my fun and profit project.
            </Typography>
            <Typography variant={'subtitle1'} color={'secondary'}>
                Niyazi Ekinci
            </Typography>
        </div>
    )
}
