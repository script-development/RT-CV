import Head from "next/head";
import Dynamic from "next/dynamic"
import React, { useEffect, useState } from "react"
import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from "@material-ui/core"
import Info from '@material-ui/icons/Info'
import Header from '../components/header'

const MatchEditor = Dynamic(
    () => import("../components/matcherEditor"),
    { ssr: false }
);

export default function TryMatcher() {
    const [infoDialogOpen, setInfoDialogOpen] = useState(false)

    const confirmInfoModal = () => {
        closeInfoModal()
        localStorage.setItem('rtcv_confirmed_try_matcher_info', 'true')
    }

    const closeInfoModal = () => {
        setInfoDialogOpen(false)
    }

    useEffect(() => {
        if (localStorage.getItem('rtcv_confirmed_try_matcher_info') != 'true')
            setInfoDialogOpen(true)
    }, [])

    return (
        <div>
            <Head><title>RT-CV try matcher</title></Head>

            <Header>
                <Button
                    variant="contained"
                    color="primary"
                    size="small"
                    startIcon={<Info />}
                    onClick={() => setInfoDialogOpen(true)}
                >
                    Info
                </Button>
            </Header>

            <MatchEditor
                top="50px"
                height="calc(100vh - 50px)"
            // expose={values => setExposedState}
            />

            <Dialog
                open={infoDialogOpen}
                onClose={closeInfoModal}
            >
                <DialogTitle id="alert-dialog-title">Try out the matcher API</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        On this page you can try out the matcher and see the matching profiles.
                        On the <b>right side you input the CV</b> and on the <b>left side you'll see the matched profiles</b> and why it was matched.
                    </DialogContentText>
                    <DialogContentText>
                        The profiles available to match are the profiles in the database.
                    </DialogContentText>
                    <DialogContentText>
                        Note that this is made to be used by developers for debugging and testing purposes.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={confirmInfoModal} color="primary" autoFocus>Oke</Button>
                </DialogActions>
            </Dialog>


        </div>
    )
}
