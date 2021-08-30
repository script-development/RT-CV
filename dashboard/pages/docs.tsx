import { RedocStandalone } from 'redoc';

export default function Docs() {
	return (
		<div className="root">
			<RedocStandalone
				specUrl="/api/v1/schema/openAPI"
			/>
			<style jsx>{`
				.root {
					background-color: white;
				}
			`}</style>
		</div>
	)
}
