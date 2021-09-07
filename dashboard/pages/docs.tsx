import { API } from "@stoplight/elements"
import Header from '../components/header'
import { useEffect, useState } from "react"

import '@stoplight/elements/styles.min.css'

export default function Docs() {
	const [currentSize, setCurrentSize] = useState(0)
	const mobile = currentSize < 1000

	useEffect(() => {
		let lastSetCurrentSize = 0;
		const onResize = () => {
			// Only set the currentSize attr if the window size has changed by a certain amount
			// This reduces the amount of unnecessary re-renders
			const newSize = Math.round(window.innerWidth / 100) * 100
			if (lastSetCurrentSize == newSize) return;

			lastSetCurrentSize = newSize;
			setCurrentSize(newSize)
		}

		window.addEventListener('resize', onResize)
		return () => window.removeEventListener('resize', onResize)
	}, [])

	return (
		<div className="root" suppressHydrationWarning={true}>
			<Header />
			{process.browser && <API
				apiDescriptionUrl="/api/v1/schema/openAPI"
				router="hash"
				hideTryIt={true}
				layout={mobile ? "stacked" : "sidebar"}
			/>}
			<style jsx>{`
				.root {
					background-color: white;
					min-height: 100vh;
				}
			`}</style>
		</div>
	)
}
