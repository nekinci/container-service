import {Button, Container, Paper, Typography} from '@mui/material';

export function Hero(){

    return (
      <Container style={{padding: '0', height: '100%'}}>
          <div style={{display:'flex', flexWrap: 'nowrap', gap: 0, justifyContent: 'space-between', padding: '100px 20px'}}>
              <div style={{}}>
                  <Typography fontWeight={'bolder'} variant={'h2'}>Containers ready to <div style={{color: '#a84e32'}}>demo</div></Typography>
                  <Typography color={'secondary'} py={'5px'} variant={'body1'}>Create your own docker image and register it anywhere like Docker Hub. After, create spec and upload and run. It's all 5 minutes.</Typography>
                  <Typography px={'2px'} pb={'4px'} variant={'h6'}>Always free!</Typography>
                  <Button href={'#getstarted'} size={'large'} variant={'outlined'} color={'secondary'}>Get started</Button>
              </div>
              <Paper style={{}} variant={'elevation'} elevation={0}>
                <img width={'600px'} alt={'Containers Hero Image'} src="/banner-1.svg" />
              </Paper>
          </div>
      </Container>
    );
}
