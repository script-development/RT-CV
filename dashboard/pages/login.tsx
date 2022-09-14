import {
	Button,
	LinearProgress,
	TextField,
} from '@mui/material'
import Head from 'next/head'
import { useRouter } from 'next/router'
import React, { useEffect, useState } from 'react'
import { fetcher } from '../src/auth'

export default function Home() {
	const [form, setForm] = useState({ id: '', key: '' })
	const [loading, setLoading] = useState(false)
	const [error, setError] = useState('')
	const router = useRouter()

	const submit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		try {
			setLoading(true)
			setError('')
			await fetcher.login(form)

			const url = new URL(location.href)
			const redirectTo = url.searchParams.get('redirectTo')
			if (typeof redirectTo == 'string' && redirectTo?.length > 1) {
				router.push(redirectTo)
			} else {
				router.push('/')
			}
		} catch (e: any) {
			setError(e?.message || e)
		} finally {
			setLoading(false)
		}
	}

	useEffect(() => {
		setForm({
			id: fetcher.getApiKeyId,
			key: fetcher.getApiKey,
		})
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
					value={form.id}
					fullWidth
					id="id"
					label="API Key ID"
					variant="filled"
					onChange={e => setForm(f => ({ ...f, id: e.target.value }))}
					disabled={loading}
					error={!!error}
				/>
				<div className="marginTop" >
					<TextField
						value={form.key}
						fullWidth
						id="key"
						label="API Key"
						variant="filled"
						onChange={e => setForm(f => ({ ...f, key: e.target.value }))}
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
				{loading ? <LinearProgress /> : undefined}
			</form>
			<style jsx>{`
				.container {
					flex-grow: 1;
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
