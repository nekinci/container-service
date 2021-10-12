import {Container, Typography} from '@mui/material';
import {Build, CloudUpload, Description, DocumentScanner, Slideshow} from '@mui/icons-material';
import {makeStyles} from '@mui/styles';
import React from 'react';


type IconComponent = 'Build' | 'CloudUpload' | 'Description' | 'Demo';
type Direction = 'left' | 'right';

type FeatureProps = {
    title: string,
    description?: string,
    iconComponent?: IconComponent,
    direction?: Direction,
    color?: string
    end?: boolean
}

function Feature({title, description, iconComponent, direction = 'right', color = '#353535', end}: FeatureProps) {
    const reverse = direction === 'left' ? 'right': 'left';
    const margin = direction === 'left' ? 'marginLeft' : 'marginRight';
    const border = direction === 'left' ? 'borderRight' : 'borderLeft';
    const rotate = direction === 'left' ? 'rotate(-25deg)' :  'rotate(25deg)';
    return (
        <div style={{display: 'flex',position: 'relative',alignItems:'center', gap: '30px', padding: '20px', flexDirection: direction === 'left' ? 'row' : 'row-reverse'}}>



            <div style={{display: 'flex', flexDirection: 'column', padding: '10px', width: '250px'}}>
                <Typography variant={'h5'}>
                    {title}
                </Typography>
                <Typography color={'secondary'}>
                    {description}
                </Typography>
            </div>

            <div style={{display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                borderRadius: '999px',
                width: '130px',
                height: '130px',
                background:color,
            }}>
                {iconComponent === 'Build' && <Build style={{color: 'white', fontSize: 60}} />}
                {iconComponent === 'CloudUpload' && <CloudUpload style={{color: 'white', fontSize: 60}} />}
                {iconComponent === 'Description' && <Description style={{color: 'white', fontSize: 60}} />}
                {iconComponent === 'Demo' && <Slideshow style={{color: 'white', fontSize: 60}} />}
            </div>
            {!end && (
                <React.Fragment>
                    <div style={{[margin]: '-32px', marginBottom: '2px', borderTop:'0.1px dashed black', overflow: 'hidden', width: '304.5px'}} />
                    <div style={{position: 'absolute', transform: rotate, padding:'10px', width: '163px', height: '163px', overflow:'hidden', background: 'transparent', [border]: '1px dashed black', borderRadius: '999px', [reverse]: '350px', top: '84px'}} />
                </React.Fragment>
            )}
        </div>
    )
}

export function Features() {

    const styles = useStyles();
    return (
        <React.Fragment>
            <Typography align={'center'} pb={'10px'} variant={'h4'} fontWeight={'bold'}>How it works!</Typography>

            <Container style={{position: 'relative'}}>
                <Feature color={'#e74744'} iconComponent={'Build'} title="Build docker image" description="Build the docker image as usual." />
                <Feature color={'#c69f34'} iconComponent={'CloudUpload'} direction={'left'} title="Upload docker image to image registry" description="Like Dockerhub" />
                <Feature color={'#9d9c52'} iconComponent={'Description'} title="Define specification" description="Define spec that contain docker image and other informations." />
                <Feature color={'#739476'} iconComponent={'CloudUpload'} direction={'left'} title="Upload Spec" description="Upload your spec file to our system." />
                <Feature end color={'#513841'} iconComponent={'Demo'} title="Demo" description="You're ready to demo." />

            </Container>
        </React.Fragment>
    )
}


const useStyles = makeStyles(() => ({
    scl: {
        transform: 'scale(5.8)'
    },
    svg: {
        width: '160px',
        position: 'absolute'
    },
    path: {
        fill: 'none',
        strokeDasharray: '2, 10', /*adjust this to control the number of dots*/
        strokeWidth:'2px'
}
}));

