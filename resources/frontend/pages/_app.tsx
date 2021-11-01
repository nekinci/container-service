import * as React from 'react';
import Head from 'next/head';
import { AppProps } from 'next/app';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { CacheProvider, EmotionCache } from '@emotion/react';
import theme from '../src/theme';
import createEmotionCache from '../src/createEmotionCache';
import Moment from "react-moment";



export const ApplicationContext = React.createContext(['asf', (val) => {}]);

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache();

interface MyAppProps extends AppProps {
  emotionCache?: EmotionCache;
}
Moment.globalFormat = 'DD.MM.yyyy HH:mm:ss'
export default function MyApp(props: MyAppProps) {

    const [currentApp, setCurrentApp] = React.useState(null);
    const value = React.useMemo(() => [currentApp, setCurrentApp], [currentApp])


  const { Component, emotionCache = clientSideEmotionCache, pageProps } = props;
  return (
    <CacheProvider value={emotionCache}>
      <Head>
        <title>My page</title>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <ThemeProvider theme={theme}>
        {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
        <CssBaseline />
        <ApplicationContext.Provider value={value}>
            <Component {...pageProps} />
        </ApplicationContext.Provider>
      </ThemeProvider>
    </CacheProvider>
  );
}
