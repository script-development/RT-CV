import {
	AccordionSummary,
	Accordion,
	AccordionDetails,
	Button,
	ButtonGroup,
	Divider,
	AccordionActions,
	Tooltip,
	LinearProgress,
} from '@material-ui/core'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import Add from '@material-ui/icons/Add'
import Delete from '@material-ui/icons/Delete'
import Edit from '@material-ui/icons/Edit'
import Head from 'next/head'
import React, { useEffect, useState } from 'react'
import { fetcher } from '../src/auth'
import { ApiKey } from '../src/types'
import { KeyModal, KeyModalKind } from '../components/keyModal'

export default function Home() {
	const [keys, setKeys] = useState(undefined as Array<ApiKey> | undefined)
	const [loadingKeys, setLoadingKeys] = useState(true)
	const [keyModal, setKeyModal] = useState({ kind: KeyModalKind.Closed, key: undefined as (undefined | ApiKey) })

	const fetchKeys = async () => {
		setLoadingKeys(true)
		setKeys(await fetcher.fetch(`/api/v1/keys`))
		setLoadingKeys(false)
	}

	useEffect(() => { fetchKeys() }, [])

	return (
		<div>
			<Head>
				<title>RT-CV home</title>
			</Head>

			<KeyModal
				kind={keyModal.kind}
				apiKey={keyModal.key}
				onClose={() => {
					setKeyModal({ kind: KeyModalKind.Closed, key: undefined })
					fetchKeys()
				}}
			/>

			<main>
				<h1>RT-CV</h1>
				<div className="cardContainer">
					<h3>Keys</h3>
					<div className="accordionHeader">
						<div className="accordionHeaderContent">
							<div>
								{loadingKeys ? 'loading..' : keys?.length} key{keys?.length == 1 ? '' : 's'}
							</div>
							<div>
								<ButtonGroup color="primary" variant="contained">
									<Tooltip title="Create Api key">
										<Button
											onClick={() => setKeyModal({ kind: KeyModalKind.Create, key: undefined })}
										>
											<Add fontSize={'small'} />
										</Button>
									</Tooltip>
								</ButtonGroup>
							</div>
						</div>
						<div className="accordionHeaderProgress">
							<div>
								{
									loadingKeys ? <LinearProgress /> : ''
								}
							</div>
						</div>
					</div>
					{keys?.map(key =>
						<Accordion key={key.id}>
							<AccordionSummary expandIcon={<ExpandMoreIcon />}>
								<div className="accordionSummary">
									<h4>
										{key.system ?
											<Tooltip title="This is a internal key created by RT-CV">
												<div
													style={{ backgroundColor: '#EF6C00' }}
													className="status"
												>System</div>
											</Tooltip>
											: ''}
										<Tooltip title={key.enabled ? 'Key can be used for authentication' : 'Key can\'t be used for authentication'}>
											<div
												style={{ backgroundColor: key.enabled ? '#00e676' : '#ff3d00' }}
												className="status"
											>{key.enabled ? 'Enabled' : 'Disabled'}</div>
										</Tooltip>
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
								<Button
									onClick={() => setKeyModal({ kind: KeyModalKind.Delete, key: key })}
								><Delete fontSize="small" style={{ marginRight: 6 }} />Delete</Button>
								<Button
									onClick={() => setKeyModal({ kind: KeyModalKind.Edit, key: key })}
								><Edit fontSize="small" style={{ marginRight: 6 }} />Edit</Button>
							</AccordionActions>
						</Accordion>
					)}
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
					overflow: hidden;
					border-top-left-radius: 4px;
					border-top-right-radius: 4px;
				}
				.accordionHeaderContent {
					padding: 10px;
					background-color: #424242;
					display: flex;
					justify-content: space-between;
					align-items: center;
				}
				.accordionHeaderProgress {
					max-height: 0px;
					transform: translate(0, -4px);
				}
			`}</style>
		</div >
	)
}
