import Head from "next/head";
import Dynamic from "next/dynamic"
import React, { useEffect, useState } from "react"
import { Button, ButtonBase, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from "@material-ui/core"
import Info from '@material-ui/icons/Info'
import ArrowBack from '@material-ui/icons/ArrowBack'
import Link from 'next/link'

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

            <div className="header">
                <div>
                    <Link href="/">
                        <ButtonBase focusRipple style={{ borderRadius: 4 }}>
                            <h1 className="title"><span className="arrowBack"><ArrowBack fontSize="small" /></span> RT-CV</h1>
                        </ButtonBase>
                    </Link>
                </div>
                <div>
                    <Button
                        variant="contained"
                        color="primary"
                        size="small"
                        startIcon={<Info />}
                        onClick={() => setInfoDialogOpen(true)}
                    >
                        Info
                    </Button>
                </div>
            </div>

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
                        On the <b>right side you input the CV</b> and on the <b>left side you'll see the matched profiles</b>.
                    </DialogContentText>
                    <DialogContentText>
                        The profiles available to match are the profiles in the database.
                    </DialogContentText>
                    <DialogContentText>
                        Note that this is made mainly developers for debugging and testing purposes.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={confirmInfoModal} color="primary" autoFocus>Oke</Button>
                </DialogActions>
            </Dialog>

            <style jsx>{`
                .header {
                    height: 50px;
                    background-color: #424242;
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding: 0 10px;
                }
                .header > div {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    height: 100%;
                }
                .header .title {
                    margin: 0;
                    padding: 5px 20px;
                }
                .header .title .arrowBack {
                    position: relative;
                    top: 2px;
                }
            `}</style>
        </div>
    )
}
