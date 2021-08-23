import { AccordionSummary, Accordion, AccordionDetails, Button, ButtonGroup, Divider, AccordionActions } from '@material-ui/core'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import Add from '@material-ui/icons/Add'
import Delete from '@material-ui/icons/Delete'
import Edit from '@material-ui/icons/Edit'
import Head from 'next/head'
import React, { useEffect, useState } from 'react'
import { fetcher } from '../src/auth'
import { ApiKey } from '../src/types'

export default function Home() {
	const [keys, setKeys] = useState(undefined as Array<ApiKey> | undefined)

	useEffect(() => {
		fetcher.fetch('/api/v1/keys').then(v => {
			setKeys(v)
		})
	}, [])

	return (
		<div>
			<Head>
				<title>RT-CV home</title>
			</Head>

			<main>
				<h1>RT-CV</h1>
				<div className="cardContainer">
					<h3>Keys</h3>
					<div className="accordionHeader">
						<div>
							{keys?.length} key{keys?.length == 1 ? '' : 's'}
						</div>
						<div>
							<ButtonGroup color="primary" variant="contained">
								<Button><Add fontSize={'small'} /></Button>
							</ButtonGroup>
						</div>
					</div>
					{keys?.map(key => <Accordion key={key.id}>
						<AccordionSummary expandIcon={<ExpandMoreIcon />}>
							<div className="accordionSummary">
								<h4>
									<div
										style={{ backgroundColor: key.enabled ? '#00e676' : '#ff3d00' }}
										className="status"
									>{key.enabled ? 'Enabled' : 'Disabled'}</div>
									{key.id}
								</h4>
								<p>{key.domains.join(', ')}</p>
							</div>
						</AccordionSummary>
						<AccordionDetails>
							<div>
								<p>id: <b>{key.id}</b></p>
								<p>key: <b>{key.key}</b></p>
								<p>domains: <b>{key.domains.join(', ')}</b></p>
								<p>enabled: <b>{key.enabled ? 'Enabled' : 'Disabled'}</b></p>
								<p>roles: <b>{key.roles}</b></p>
							</div>
						</AccordionDetails>
						<Divider />
						<AccordionActions>
							<Button><Delete fontSize="small" style={{ marginRight: 6 }} />Delete</Button>
							<Button><Edit fontSize="small" style={{ marginRight: 6 }} />Edit</Button>
						</AccordionActions>
					</Accordion>)}
				</div>
			</main>

			<style jsx>{`
				main {
					padding: 50px 20px;
					display: flex;
					justify-content: center;
					flex-direction: column;
					align-items: center;
				}
				.cardContainer {
					padding: 10px;
					width: 700px;
					box-sizing: border-box;
					max-width: calc(100vw - 20px);
				}
				.accordionSummary {
					display: flex;
					justify-content: space-between;
					width: 100%;
				}
				.accordionSummary .status {
					display: inline-block;
					padding: 0.15rem 0.5rem 0.05rem 0.5rem;
					margin-right: 10px;
					font-size: 0.8rem;
					color: black;
					border-radius: 10px;
				}
				.accordionHeader {
					padding: 10px;
					border-top-left-radius: 4px;
					border-top-right-radius: 4px;
					background-color: #424242;
					display: flex;
					justify-content: space-between;
					align-items: center;
				}
			`}</style>
		</div >
	)
}