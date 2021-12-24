import { Button, Icon } from '@material-ui/core'
import Head from 'next/head'
import Link from 'next/link'
import React from 'react'
import KeysCard from '../components/keysCard'
import SecretsCard from '../components/secretsCard'
import Statistics from '../components/statistics'

const ButtonStyle = { marginRight: 5, marginBottom: 5 }

export default function Home() {
	return (
		<div>
			<Head><title>RT-CV home</title></Head>

			<main>
				<h1>RT-CV</h1>

				<div className="appLinks">
					<Link href="/tryMatcher">
						<Button style={ButtonStyle} color="primary" variant="outlined" startIcon={<Icon>construction</Icon>}>Try the CV matcher</Button>
					</Link>
					<Link href="/docs">
						<Button style={ButtonStyle} color="primary" variant="outlined" startIcon={<Icon>menu_book</Icon>}>API docs</Button>
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
