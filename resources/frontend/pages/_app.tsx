import * as React from 'react';
import Head from 'next/head';
import { AppProps } from 'next/app';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { CacheProvider, EmotionCache } from '@emotion/react';
import theme from '../src/theme';
import createEmotionCache from '../src/createEmotionCache';
import Moment from "react-moment";
import GenericSnackbar from "../src/components/Snackbar/GenericSnackbar";
import {RunApp} from "../src/components/modals/RunApp/RunApp";
import {Login} from "../src/components/modals/Login/Login";
import {Typography} from "@mui/material";


export const ApplicationContext = React.createContext([null, (val) => {}]);

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache();

interface MyAppProps extends AppProps {
  emotionCache?: EmotionCache;
}
Moment.globalFormat = 'DD.MM.yyyy HH:mm:ss'
export default function MyApp(props: MyAppProps) {

    const [currentApp, setCurrentApp] = React.useState(null);
    const value = React.useMemo(() => [currentApp, setCurrentApp], [currentApp])
    const [view, setView] = React.useState(true);

    React.useEffect(() => {

        const listener = () => {
           if (window?.innerWidth < 989) {
               setView(false);
           } else {
               setView(true)
           }
        }

        window.onload = listener;
        window.addEventListener('resize', listener);

        return () => {
            window.removeEventListener('resize', listener);
        }
    }, []);
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

          {view && (
              <ApplicationContext.Provider value={value}>
                  <>
                      <GenericSnackbar />
                      <RunApp />
                      <Login />
                      <Component {...pageProps} />
                  </>
              </ApplicationContext.Provider>
          )}

          {!view && (
              <Typography variant={'h6'}>
                  Please enter from a larger screen for a better user experience.
              </Typography>
          )}
      </ThemeProvider>
    </CacheProvider>
  );
}
