import Head from "next/head";
import Dynamic from "next/dynamic"
import React from "react"
import { ButtonBase } from "@material-ui/core";

const MatchEditor = Dynamic(
    () => import("../components/matcherEditor"),
    { ssr: false }
);

export default function TryMatcher() {
    return (
        <div>
            <Head>
                <title>RT-CV home</title>
            </Head>

            <div className="header">
                <ButtonBase focusRipple style={{ borderRadius: 4 }}>
                    <h1 className="title">RT-CV</h1>
                </ButtonBase>
            </div>

            <MatchEditor
                style={{ height: 'calc(100vh - 50px)', width: '100%' }}
            />

            <style jsx>{`
                .header {
                    height: 50px;
                    background-color: #424242;
                    display: flex;
                    justify-content: flex-start;
                    align-items: center;
                    padding: 0 10px;
                }
                .header .title {
                    margin: 0;
                    padding: 5px 20px;
                }
            `}</style>
        </div>
    )
}
