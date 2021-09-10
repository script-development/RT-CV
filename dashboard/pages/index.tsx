import { Button, Icon } from '@material-ui/core'
import Head from 'next/head'
import Link from 'next/link'
import React, { useState } from 'react'
// import { fetcher } from '../src/auth'
// import { formatRFC3339, subWeeks } from 'date-fns'
import KeysCard from '../components/keysCard'
import SecretsCard from '../components/secretsCard'

export default function Home() {
	const [loading, setLoading] = useState(true)
	// const [analyticsPeriod, setAnalyticsPeriod] = useState<'week' | 'day' | 'hour'>('week')

	// const fetchAnalytics = async () => {
	// 	const from = formatRFC3339(subWeeks(new Date(), 1));
	// 	const to = formatRFC3339(new Date());
	// 	return await fetcher.fetch(`/api/v1/analytics/matches/period/${from}/${to}`)
	// }

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
