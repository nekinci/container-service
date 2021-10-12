import React from "react";
import Terminal, {LineType} from "react-terminal-ui";
import {DashboardMain} from "../../src/components/DashboardMain/DashboardMain";
import {AuthUtil} from "../../src/util/AuthUtil";


export default function Logs(){

    const [ws, setWs] = React.useState<WebSocket>(null);
    const [lines, setLines] = React.useState([
        {type: LineType.Output, value: 'Welcome to the React Terminal UI Demo!'},
        {type: LineType.Input, value: 'Some previous input received'},
    ]);

    React.useEffect(() => {
        setWs(new WebSocket("ws://localhost:8070/ws?token=" + AuthUtil.getInformation()?.token));
    }, []);

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
            {lines.map(x => <div key={x}>{x.value}</div>)}
        </DashboardMain>
    )
}
