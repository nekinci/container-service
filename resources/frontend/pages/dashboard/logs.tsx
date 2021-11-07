import React from "react";
import Terminal, {LineType} from "react-terminal-ui";
import DashboardMain from "../../src/components/DashboardMain/DashboardMain";
import {AuthUtil} from "../../src/util/AuthUtil";
import {useApp} from "../../src/hooks/useApp";
import Head from "next/head";
import {getEnvironment} from "../../environment/environment";


export default function Logs(){

    const [ws, setWs] = React.useState<WebSocket | any>(null);
    const [currentApp] = useApp();
    const [lines, setLines] = React.useState([
        {type: LineType.Input, value: 'Some previous input received'},
    ]);

    React.useEffect(() => {
        if (currentApp != null){
            setWs(new WebSocket(getEnvironment().wsUrl + "logs?token=" + AuthUtil.getInformation()?.token + '&currentApp=' + currentApp));
            return () => {
                if (ws != null){
                    ws.close()
                }
            }
        }
    }, [currentApp]);

    React.useEffect(() => {
        if (ws != null){
            ws.onopen = function (){
                ws.onmessage = function (message){
                    console.log(JSON.stringify(message.data))
                    setLines(l => [...l, {type: LineType.Output, value: message.data.replace(/\r\n/g, "\n")}])
                }
            }
        }
    }, [ws]);

    return (
       <>
           <Head>
               <title>Logs</title>
           </Head>
           <DashboardMain>
               {ws != null && (
                   <Terminal name={'Logs'} prompt={'$'} onInput={ terminalInput => {
                       console.log(terminalInput);
                       ws.send(terminalInput + "\n")
                   } } lineData={lines} />
               )}
           </DashboardMain>
       </>
    )
}
