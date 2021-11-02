import { Button, Icon } from '@material-ui/core'
import Head from 'next/head'
import Link from 'next/link'
import React, { useEffect } from 'react'
import KeysCard from '../components/keysCard'
import SecretsCard from '../components/secretsCard'
import Statistics from '../components/statistics'
import { fetcher } from '../src/auth'

function getWebsocketUrl() {
	const url = fetcher.getAPIPath(`/api/v1/events/ws/${fetcher.authorizationValue}`, true)
	if (url[0] == '/') {
		return `ws${location.protocol == 'https:' ? 's' : ''}://${location.host}${url}`
	} else {
		return url
	}
}

export default function Home() {
	const connectToSocket = () => {
		try {
			const socket = new WebSocket(getWebsocketUrl())

			let open = true
			const close = () => {
				if (!open) { return }
				open = false
				socket.onmessage = null
				socket.onopen = null
				socket.onerror = null
				socket.onclose = null
				setTimeout(() => connectToSocket(), 5000)
			}

			socket.onmessage = (ev: MessageEvent<any>) => {
				console.log('received message', ev)
			}
			socket.onopen = () => {
				console.log('connected')
			}
			socket.onerror = (e) => {
				console.log('disconnected from websocket, error:', e)
				close()
			}
			socket.onclose = (e) => {
				console.log('disconnected from websocket, error:', e)
				close()
			}
			return () => socket.close(1000, 'navigating to different route')
		} catch (e) {
			console.error(e)
			return () => { }
		}
	}

	useEffect(() => {
		let closeConn = undefined
		closeConn = connectToSocket()
		return closeConn
	}, [])

	return (
		<div>
			<Head><title>RT-CV home</title></Head>

			<main>
				<h1>RT-CV</h1>

				<div className="appLinks">
					<Link href="/tryMatcher">
						<Button style={{ marginRight: 5, marginBottom: 5 }} color="primary" variant="outlined" startIcon={<Icon>construction</Icon>}>Try the CV matcher</Button>
					</Link>
					<Link href="/docs">
						<Button style={{ marginRight: 5, marginBottom: 5 }} color="primary" variant="outlined" startIcon={<Icon>menu_book</Icon>}>API docs</Button>
					</Link>
				</div>

				<KeysCard />

				<SecretsCard />

				<Statistics />
			</main>

			<style jsx>{`
				main {
					padding: 50px 20px;
					display: flex;
					justify-content: center;
					flex-direction: column;
					align-items: center;
				}
				.appLinks {
					padding: 10px;
					width: 700px;
					box-sizing: border-box;
					max-width: calc(100vw - 20px);
					display: flex;
					flex-wrap: wrap;
				}
			`}</style>
		</div >
	)
}
