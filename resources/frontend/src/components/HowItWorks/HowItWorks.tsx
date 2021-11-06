import {Container, Divider, Typography} from "@mui/material";
import {Build, BuildCircle, CloudUpload, Description, Slideshow} from "@mui/icons-material";
import React from "react";


function Feature({title, subtitle, backgroundColor, sequence, right = false, iconComponent}){

    return (
        <div style={{display: 'flex', gap: '10px', flexDirection: right ? 'row-reverse': 'row', padding: '10px'}}>
            <div style={{
                display: 'flex',
                justifyContent:'center',
                alignItems: 'center',
                width: '130px',
                height: '130px',
                background: backgroundColor,
                borderRadius: '999px',
                position: 'relative'
            }}>
                <div style={{position: 'absolute',
                    top: 0, [right ? 'left': 'right']: 10,
                    background: 'white',
                    borderRadius: '999px',
                    width: '40px',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems:'center',
                    border: '1px dashed black',
                    height: '40px'}}
                >
                    <Typography variant={'body1'} fontWeight={'bold'}>{sequence}</Typography>
                </div>
                {iconComponent === 'Build' && <Build style={{color: 'white', fontSize: 60}} />}
                {iconComponent === 'CloudUpload' && <CloudUpload style={{color: 'white', fontSize: 60}} />}
                {iconComponent === 'Description' && <Description style={{color: 'white', fontSize: 60}} />}
                {iconComponent === 'Demo' && <Slideshow style={{color: 'white', fontSize: 60}} />}
            </div>
            <div style={{width: '250px'}}>
                <Typography variant={'h5'}>{title}</Typography>
                <Typography color={'secondary'}>{subtitle}</Typography>
            </div>
        </div>
    )
}

export function HowItWorks(){

    return (
        <div id={'getstarted'} style={{padding: '20px'}}>
            <Typography align={'center'} variant={'h4'} fontWeight={'bold'}>
                How it works!
            </Typography>
            <Container style={{display: 'flex', justifyContent: 'center', flexDirection: 'column', marginTop: '14px', width: '1100px'}}>
                <Feature
                    title={'Build'}
                    subtitle={'Build docker image'}
                    backgroundColor={'#e74744'}
                    sequence={'01'}
                    iconComponent={'Build'}
                />

                <Feature
                    title={'Upload docker image to image registry'}
                    subtitle={'Like Dockerhub'}
                    backgroundColor={'#c69f34'}
                    sequence={'02'}
                    right={true}
                    iconComponent={'CloudUpload'}
                />

                <Feature
                    title={'Define specification'}
                    subtitle={'Define spec that contain docker image and other informations.'}
                    backgroundColor={'#9d9c52'}
                    sequence={'03'}
                    iconComponent={'Description'}
                />

                <Feature
                    title={'Upload Spec'}
                    subtitle={'Upload your spec file to our system.'}
                    backgroundColor={'#739476'}
                    sequence={'04'}
                    right={true}
                    iconComponent={'CloudUpload'}
                />

                <Feature
                    title={'Demo'}
                    subtitle={"You're ready to demo."}
                    backgroundColor={'#513841'}
                    sequence={'05'}
                    iconComponent={'Demo'}
                />

            </Container>
        </div>
    )
}