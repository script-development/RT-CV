import Button from '@material-ui/core/Button'
import LinearProgress from '@material-ui/core/LinearProgress'
import TextField from '@material-ui/core/TextField'
import Head from 'next/head'
import { useRouter } from 'next/router'
import React, { useEffect, useState } from 'react'
import { fetcher } from '../src/auth'

export default function Home() {
	const [apiKeyId, setApiKeyId] = useState('')
	const [apiKey, setApiKey] = useState('')
	const [loading, setLoading] = useState(false)
	const [error, setError] = useState('')
	const router = useRouter()

	const submit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		try {
			setLoading(true)
			setError('')
			await fetcher.login(apiKey, apiKeyId)
			router.push('/')
		} catch (e: any) {
			setError(e?.message || e)
		} finally {
			setLoading(false)
		}
	}

	useEffect(() => {
		setApiKeyId(fetcher.getApiKeyId)
		setApiKey(fetcher.getApiKey)
	}, [])

	return (
		<div className="container">
			<Head>
				<title>RT-CV Login</title>
			</Head>

			<h1>RT-CV Login</h1>
			<form noValidate onSubmit={submit} >
				<p>Insert a api key with the <b>Dashboard</b> role</p>
				<TextField
					value={apiKeyId}
					fullWidth
					id="id"
					label="API Key ID"
					variant="filled"
					onChange={e => setApiKeyId(e.target.value)}
					disabled={loading}
					error={!!error}
				/>
				<div className="marginTop" >
					<TextField
						value={apiKey}
						fullWidth
						id="key"
						label="API Key"
						variant="filled"
						onChange={e => setApiKey(e.target.value)}
						disabled={loading}
						error={!!error}
						helperText={error}
					/>
				</div>
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
				.marginTop {
					margin-top: 20px;
				}
			`}</style>
		</div>
	)
}
