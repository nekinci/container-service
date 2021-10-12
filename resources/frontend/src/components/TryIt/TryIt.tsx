import React from 'react';
import {Button, Card, CardActions, CardContent, Container, Typography} from '@mui/material';

export function TryIt() {

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
                            name: example-container
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
                        <Button>Upload</Button>
                    </CardActions>
                </Card>
            </Container>
        </React.Fragment>
    )
}
