import DashboardMain from "../../src/components/DashboardMain/DashboardMain";
import React from "react";
import {AuthUtil} from "../../src/util/AuthUtil";
import Terminal, {LineType} from "react-terminal-ui";
import {useApp} from "../../src/hooks/useApp";
import Head from "next/head";

export default function TerminalPage() {

    const [prompt, setPrompt] = React.useState<string>("$");
    const [hostname, setHostname] = React.useState("");
    const [whoami, setWhoami] = React.useState("");
    const [pwd, setPwd] = React.useState("");
    const [ws, setWs] = React.useState<WebSocket>(null);
    const [val, setVal] = React.useState(0);
    const [currentApp] = useApp();

    const virtualTerminalHandler = () => {
        if (ws != null){
            // Linux always helps us perfectly :) I ♥️ Linux
            ws.send('echo ":x:pwd: $(pwd)"\n');
            ws.send('echo ":x:whoami: $(whoami)"\n');
            ws.send('echo ":x:hostname: $(hostname)"\n');

        }
    };

    const [lines, setLines] = React.useState([
        {type: LineType.Output, value: 'Welcome to the Container Terminal! This is a very basic terminal that communicate between container and ui.'},
        {type: LineType.Input, value: 'Some previous input received'},
    ]);

    React.useEffect(() => {
        if (currentApp != null){
            setWs(new WebSocket("ws://localhost:8070/terminal?token=" + AuthUtil.getInformation()?.token + '&currentApp=' + currentApp));
        }
    }, [currentApp])


    React.useEffect(() => {
        setPrompt(`${whoami}@${hostname} ${pwd} $ `);
    }, [val]);

    React.useEffect(() => {
        if (ws != null){
            ws.onopen = function (){
                virtualTerminalHandler();
                let i = 0;
                ws.onmessage = function (message){

                    if ((message?.data as string)?.startsWith(":x:hostname: ")){
                        let msg = message.data as string;
                        msg = msg.replace(":x:hostname: ", "");
                        setHostname(msg);
                        i++;
                    }

                    else if ((message?.data as string)?.startsWith(":x:whoami: ")){
                        let msg = message.data as string;
                        msg = msg.replace(":x:whoami: ", "");
                        setWhoami(msg);
                        i++;
                    }

                    else if ((message?.data as string)?.startsWith(":x:pwd: ")){
                        let msg = message.data as string;
                        msg = msg.replace(":x:pwd: ", "");
                        setPwd(msg);
                        i++;
                    }

                    else {
                        setLines(l => [...l, {type: LineType.Output, value: message.data.replace(/\r\n/g, "\n")}]);
                    }

                    if (i >= 3){
                        i = 0;
                        setVal(val => val + 1)
                    }

                }

            }
        }
    }, [ws]);



    return (
       <>
           <Head>
               <title>Terminal</title>
           </Head>
           <DashboardMain>
               {ws != null && (
                   <Terminal name={'Terminal'} prompt={prompt} onInput={ terminalInput => {
                       setLines(l => [...l, {type: LineType.Input, value: terminalInput}])
                       ws.send(terminalInput + "\n")
                       virtualTerminalHandler();
                       if (terminalInput === 'clear'){
                           setLines([])
                       }
                   } } lineData={lines} />
               )}
           </DashboardMain>
       </>
    );
}