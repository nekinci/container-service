import {Alert, Snackbar} from "@mui/material";
import {getEnvironment} from "../../../environment/environment";
import React from "react";
import event from "../../util/Event";

export default function GenericSnackbar(){

    const [snackbarOpen, setSnackbarOpen] = React.useState(false);
    const [snackbarContent, setSnackbarContent] = React.useState("");
    const [alertType, setAlertType] = React.useState("error")
    const [anchor, setAnchor] = React.useState({vertical: 'top', horizontal: 'right'})

    React.useEffect(() => {

        event.on('snackbar', (msg, type = 'success', anchor = {vertical: 'top', horizontal: 'right'}) => {
            console.log(msg)
            setSnackbarOpen(true)
            setSnackbarContent(msg)
            setAlertType(type)
            setAnchor(anchor)
        })
    }, []);

    return (
       <React.Fragment>
           <Snackbar
               anchorOrigin={anchor}
               open={snackbarOpen}
               autoHideDuration={getEnvironment().snackbarHideDuration}
               onClose={() => setSnackbarOpen(false)}
           >
              <Alert severity={alertType} color={alertType}>
                  {snackbarContent}
              </Alert>
           </Snackbar>
       </React.Fragment>
    )
}