import { Link, ButtonBase } from "@mui/material"
import ArrowBack from "@mui/icons-material/ArrowBack"
import React, { CSSProperties } from "react"

interface HeaderProps {
    children?: React.ReactNode
    arrowBackStyle?: CSSProperties
}

export default function Header({ children, arrowBackStyle }: HeaderProps) {
    return (
        <div className="header">
            <div>
                <Link href="/">
                    <ButtonBase focusRipple style={{ borderRadius: 4 }}>
                        <h1 className="title"><span className="arrowBack" style={arrowBackStyle}><ArrowBack fontSize="small" /></span> RT-CV</h1>
                    </ButtonBase>
                </Link>
            </div>
            <div>
                {children}
            </div>

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
                .header :global(a) {
                    color: white;
                    text-decoration: none;
                }
                .header .title {
                    color: white;
                    text-decoration: none;
                    margin: 0;
                    padding: 5px 20px;
                    font-size: 22px;
                    font-weight: bold;
                }
                .header .title .arrowBack {
                    position: relative;
                    top: 2px;
                }
            `}</style>
        </div>
    )
}
