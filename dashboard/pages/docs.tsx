import Header from '../components/header'
import { useEffect, useState } from "react"
import Dynamic from 'next/dynamic'
import { LinearProgress } from '@material-ui/core'
import { fetcher } from '../src/auth'

import '@stoplight/elements/styles.min.css'
const API = Dynamic<any>(
	() => import('@stoplight/elements').then(d => d.API),
	{
		ssr: false,
		loading: loader
	})

export default function Docs() {
	const [currentSize, setCurrentSize] = useState(0)
	const mobile = currentSize < 1000
	const [openApiPath, setOpenApiPath] = useState('')

	useEffect(() => {
		setOpenApiPath(fetcher.getAPIPath("/api/v1/schema/openAPI"))

		let lastSetCurrentSize = 0;
		const onResize = () => {
			// Only set the currentSize attr if the window size has changed by a certain amount
			// This reduces the amount of unnecessary re-renders
			const newSize = Math.round(window.innerWidth / 100) * 100
			if (lastSetCurrentSize == newSize) return;

			lastSetCurrentSize = newSize;
			setCurrentSize(newSize)
		}

		onResize()
		window.addEventListener('resize', onResize)
		return () => window.removeEventListener('resize', onResize)
	}, [])

	return (
		<div className="rtcvDocs" suppressHydrationWarning={true}>
			<Header
				arrowBackStyle={{ top: 0 }}
			/>
			{openApiPath ?
				<API
					apiDescriptionUrl={openApiPath}
					router="hash"
					hideTryIt={true}
					layout={mobile ? "stacked" : "sidebar"}
				/>
				: ''}
			<style global jsx>{`
				.rtcvDocs {
					background-color: white;
					min-height: 100vh;
					color: black;
				}
				.rtcvDocs a {
					color: black;
				}
			`}</style>
		</div>
	)
}

function loader() {
	return (
		<div>
			<LinearProgress />
			<p className="loading">Loading...</p>
			<style jsx>{`
				.loading {
					margin: 20px;
					text-align: center;
					font-size: 20px;
				}
			`}</style>
		</div>
	)
}
