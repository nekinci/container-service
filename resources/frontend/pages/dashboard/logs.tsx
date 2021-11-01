import React from "react";
import Terminal, {LineType} from "react-terminal-ui";
import DashboardMain from "../../src/components/DashboardMain/DashboardMain";
import {AuthUtil} from "../../src/util/AuthUtil";
import {useApp} from "../../src/hooks/useApp";


export default function Logs(){

    const [ws, setWs] = React.useState<WebSocket>(null);
    const [currentApp] = useApp();
    const [lines, setLines] = React.useState([
        {type: LineType.Output, value: 'Welcome to the React Terminal UI Demo!'},
        {type: LineType.Input, value: 'Some previous input received'},
    ]);

    React.useEffect(() => {
        if (currentApp != null){
            setWs(new WebSocket("ws://localhost:8070/logs?token=" + AuthUtil.getInformation()?.token + '&currentApp=' + currentApp));
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
        <DashboardMain>
            {ws != null && (
                <Terminal name={'Logs'} prompt={'$'} onInput={ terminalInput => {
                    console.log(terminalInput);
                    ws.send(terminalInput + "\n")
                } } lineData={lines} />
            )}
        </DashboardMain>
    )
}
