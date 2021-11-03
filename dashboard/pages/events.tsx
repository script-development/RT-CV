import Head from 'next/head'
import { useEffect, useState } from 'react'
import { fetcher } from '../src/auth'
import Header from '../components/header'

function getWebsocketUrl() {
    const url = fetcher.getAPIPath(`/api/v1/events/ws/${fetcher.authorizationValue}`, true)
    if (url[0] == '/') {
        return `ws${location.protocol == 'https:' ? 's' : ''}://${location.host}${url}`
    } else {
        return url
    }
}

export default function Events() {
    const [connectionStatus, setConnectionStatus] = useState({
        connected: false,
        msg: '',
    })

    const connectToSocket = () => {
        try {
            const socket = new WebSocket(getWebsocketUrl())

            let open = true
            const close = () => {
                if (!open) { return }
                setConnectionStatus({
                    connected: false,
                    msg: 'trying to reconnect to websocket in 5 seconds',
                })
                open = false
                socket.onmessage = null
                socket.onopen = null
                socket.onerror = null
                socket.onclose = null
                setTimeout(() => {
                    setConnectionStatus({
                        connected: false,
                        msg: 'reconnecting..',
                    })
                    connectToSocket()
                }, 5000)
            }

            socket.onmessage = (ev: MessageEvent<any>) => {
                console.log('received message', ev)
            }
            socket.onopen = () => {
                setConnectionStatus({
                    connected: true,
                    msg: '',
                })
                console.log('connected to websocket')
            }
            socket.onerror = (e) => {
                console.log('disconnected from websocket, error:', e)
                close()
            }
            socket.onclose = () => {
                console.log('websocket connection closed')
                close()
            }
            return () => {
                open = false;
                socket.close(1000, 'navigating to different route')
            }
        } catch (e) {
            console.error(e)
            setConnectionStatus({
                connected: false,
                msg: 'unable to connect with websocket',
            })
            return () => { }
        }
    }

    useEffect(() => {
        const closeConn = connectToSocket()
        return closeConn
    }, [])

    return (
        <div>
            <Header />

            <Head><title>RT-CV events</title></Head>

            <div className="status">
                <div className="dot" style={{ backgroundColor: connectionStatus.connected ? '#8bc34a' : '#ff5722' }} />
                {connectionStatus.connected ? 'connected to server' : connectionStatus.msg ? 'disconnected, ' + connectionStatus.msg : 'disconnected'}
            </div>

            <style jsx>{`
                .status {
                    padding: 10px;
                    text-align: center;
                    color: rgba(255, 255, 255, 0.7);
                }
                .status .dot {
                    display: inline-block;
                    height: 10px;
                    width: 10px;
                    background-color: white;
                    border-radius: 5px;
                    margin-right: 4px;
                }
            `}</style>
        </div>
    )
}
