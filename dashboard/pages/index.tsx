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
	Breadcrumbs,
} from '@material-ui/core'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import Add from '@material-ui/icons/Add'
import Delete from '@material-ui/icons/Delete'
import Edit from '@material-ui/icons/Edit'
import Visibility from '@material-ui/icons/Visibility'
import Head from 'next/head'
import React, { useEffect, useState } from 'react'
import { fetcher } from '../src/auth'
import { ApiKey, Secret } from '../src/types'
import { KeyModal } from '../components/keyModal'
import { ModalKind } from '../components/modal'
import { SecretModal } from '../components/secretModal'

export default function Home() {
	const [keys, setKeys] = useState(undefined as Array<ApiKey> | undefined)
	const [secrets, setSecrets] = useState(undefined as Array<Secret> | undefined)
	const [loading, setLoading] = useState(true)
	const [keyModal, setKeyModal] = useState({ kind: ModalKind.Closed, key: undefined as (undefined | ApiKey) })
	const [secretModal, setSecretModal] = useState({ kind: ModalKind.Closed, secret: undefined as (undefined | Secret) })

	const fetchData = async () => {
		try {
			setLoading(true)
			const [keys, secrets] = await Promise.all([
				fetcher.fetch(`/api/v1/keys`),
				fetcher.fetch(`/api/v1/secrets/otherKey`),
			]);
			setKeys(keys)
			setSecrets(secrets)
		} finally {
			setLoading(false)
		}
	}

	useEffect(() => { fetchData() }, [])

	return (
		<div>
			<Head>
				<title>RT-CV home</title>
			</Head>

			<KeyModal
				kind={keyModal.kind}
				apiKey={keyModal.key}
				onClose={() => {
					setKeyModal({ kind: ModalKind.Closed, key: undefined })
					fetchData()
				}}
			/>

			<SecretModal
				kind={secretModal.kind}
				secret={secretModal.secret}
				onClose={() => {
					setSecretModal({ kind: ModalKind.Closed, secret: undefined })
					fetchData()
				}}
			/>

			<main>
				<h1>RT-CV</h1>
				<div className="cardContainer">
					<h3>Keys</h3>
					<div className="accordionHeader">
						<div className="accordionHeaderContent">
							<div>
								{loading
									? 'loading..'
									: keys?.length
										? `${keys?.length} key${keys?.length == 1 ? '' : 's'}`
										: 'No keys'
								}
							</div>
							<div>
								<ButtonGroup color="primary" variant="contained">
									<Tooltip title="Create Api key">
										<Button
											onClick={() => setKeyModal({ kind: ModalKind.Create, key: undefined })}
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
									loading ? <LinearProgress /> : ''
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
									onClick={() => setKeyModal({ kind: ModalKind.Delete, key: key })}
								><Delete fontSize="small" style={{ marginRight: 6 }} />Delete</Button>
								<Button
									onClick={() => setKeyModal({ kind: ModalKind.Edit, key: key })}
								><Edit fontSize="small" style={{ marginRight: 6 }} />Edit</Button>
							</AccordionActions>
						</Accordion>
					)}
				</div>

				<div className="cardContainer">
					<h3>Secrets</h3>
					<div className="accordionHeader">
						<div className="accordionHeaderContent">
							<div>
								{loading
									? 'loading..'
									: secrets?.length
										? `${secrets?.length} secret${secrets?.length == 1 ? '' : 's'}`
										: 'No secrets'
								}
							</div>
							<div>
								<ButtonGroup color="primary" variant="contained">
									<Tooltip title="Create secret">
										<Button
											onClick={() => setSecretModal({ kind: ModalKind.Create, secret: undefined })}
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
									loading ? <LinearProgress /> : ''
								}
							</div>
						</div>
					</div>
					{secrets?.map((secret, idx) =>
						<div key={secret.id} className={"simpleRow" + (secrets.length == (idx + 1) ? ' last' : '')}>
							<Breadcrumbs>
								<p>{secret.keyId}</p>
								<b style={{ color: "white" }}>{secret.key}</b>
							</Breadcrumbs>
							<div>
								<Tooltip title="View secret contents">
									<Button onClick={() => setSecretModal({ kind: ModalKind.View, secret: secret })}>
										<Visibility fontSize="small" />
									</Button>
								</Tooltip>
								<Tooltip title="Delete secret">
									<Button onClick={() => setSecretModal({ kind: ModalKind.Delete, secret: secret })}>
										<Delete fontSize="small" />
									</Button>
								</Tooltip>
							</div>
						</div>
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
					background-color: #424242;
				}
				.accordionHeaderContent {
					padding: 10px;
					display: flex;
					justify-content: space-between;
					align-items: center;
				}
				.accordionHeaderProgress {
					max-height: 0px;
					transform: translate(0, -4px);
				}
				.simpleRow {
					border-top: 1px solid rgba(255, 255, 255, 0.12);
					background-color: #424242;
					padding: 10px;
					display: flex;
					justify-content: space-between;
					align-items: center;
				}
				.simpleRow.last {
					border-bottom-left-radius: 4px;
					border-bottom-right-radius: 4px;
				}
			`}</style>
		</div >
	)
}
