import Button from '@material-ui/core/Button'
import LinearProgress from '@material-ui/core/LinearProgress'
import TextField from '@material-ui/core/TextField'
import Head from 'next/head'
import React, { useState } from 'react'
import { fetcher } from '../src/auth'

export default function Home() {
	const [apiKey, setApiKey] = useState('')
	const [apiKeyId, setApiKeyId] = useState('')
	const [loading, setLoading] = useState(false)
	const [error, setError] = useState('')

	const submit = (e: React.FormEvent<HTMLFormElement>) => {
		try {
			e.preventDefault()
			setLoading(true)
			fetcher.login(apiKey, apiKeyId)
		} catch (e) {
			setError(e.message)
		} finally {
			setLoading(false)
		}
	}

	return (
		<div className="container">
			<Head>
				<title>RT-CV Login</title>
				<link rel="icon" href="/favicon.ico" />
			</Head>

			<h1>RT-CV Login</h1>
			<form noValidate onSubmit={submit} >
				<p>Insert a api key with the <b>Information Obtainer</b> and <b>Controller</b> role</p>
				<TextField
					fullWidth
					id="key"
					label="API Key"
					variant="filled"
					onChange={e => setApiKey(e.target.value)}
					disabled={loading}
				/>
				<TextField
					fullWidth
					id="id"
					label="API Key ID"
					variant="filled"
					onChange={e => setApiKeyId(e.target.value)}
					disabled={loading}
				/>
				<div className="actions">
					<Button
						color="secondary"
						type="submit"
						disabled={loading}
					>Login</Button>
				</div>
				<LinearProgress hidden={!loading} />
			</form>
			<style jsx>{`
				.container {
					min-height: 100vh;
					display: flex;
					justify-content: center;
					align-items: center;
					flex-direction: column;
				}
				form {
					padding: 30px 10px;
					margin: 10px;
					width: 300px;
					border-radius: 10px;
					background-color: #424242;
				}
				.actions {
					padding-top: 10px;
				}
				form p {
					margin-bottom: 6px;
				}
			`}</style>
		</div>
	)
}
