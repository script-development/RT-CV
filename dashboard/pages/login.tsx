import Button from '@material-ui/core/Button'
import TextField from '@material-ui/core/TextField'
import Head from 'next/head'
import React, { useState } from 'react'
import { login } from '../src/auth'

export default function Home() {
	const [apiToken, setApiToken] = useState('')

	const submit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		login(apiToken)
	}

	return (
		<div className="container">
			<Head>
				<title>Home</title>
				<link rel="icon" href="/favicon.ico" />
			</Head>

			<h1>Login</h1>
			<form noValidate onSubmit={submit} >
				<TextField fullWidth id="token" label="API Token" variant="filled" onChange={e => setApiToken(e.target.value)} />
				<div className="actions">
					<Button color="secondary" type="submit">Login</Button>
				</div>
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
			`}</style>
		</div>
	)
}
