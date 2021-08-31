import { RedocStandalone } from 'redoc';
import Header from '../components/header'
import { secondaryColor } from '../src/theme'

export default function Docs() {
	return (
		<div className="root">
			<Header />
			<RedocStandalone
				specUrl="/api/v1/schema/openAPI"
				options={{
					theme: { colors: { primary: { main: secondaryColor } } },
				}}
			/>
			<style jsx>{`
				.root {
					background-color: white;
					min-height: 100vh;
				}
			`}</style>
		</div>
	)
}
