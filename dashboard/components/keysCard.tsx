import {
    AccordionSummary,
    Accordion,
    AccordionActions,
    AccordionDetails,
    Button,
    ButtonGroup,
    Divider,
    Tooltip,
} from '@mui/material'
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'
import Edit from '@mui/icons-material/Edit'
import Add from '@mui/icons-material/Add'
import Delete from '@mui/icons-material/Delete'
import Password from '@mui/icons-material/Password'
import VpnKey from '@mui/icons-material/VpnKey'
import React, { useEffect, useState } from 'react'
import { Roles } from '../src/roles'
import Card from './card'
import { getKeys } from '../src/auth'
import { ApiKey } from '../src/types'
import { ModalKind } from './modal'
import { KeyModal } from './keyModal'
import { ScraperAuthFileModal } from './ScraperAuthFileModal'
import { ScraperUsersModal } from './ScraperUsersModal'

export default function KeysCard() {
    const [loading, setLoading] = useState(true)
    const [keys, setKeys] = useState(undefined as Array<ApiKey> | undefined)
    const [modal, setModal] = useState({ kind: ModalKind.Closed, key: undefined as (undefined | ApiKey) })
    const [scraperAuthFileModal, setScraperAuthFileModal] = useState<undefined | ApiKey>(undefined)
    const [scraperAuthUsersModal, setScraperAuthUsersModal] = useState<undefined | ApiKey>(undefined)

    const fetchData = async () => {
        try {
            setLoading(true)
            const keys = await getKeys()
            setKeys(keys)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchData() }, [])

    return (
        <Card
            title="Keys"
            loading={loading}
            headerContent={<>
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
                                onClick={() => setModal({ kind: ModalKind.Create, key: undefined })}
                            >
                                <Add fontSize={'small'} />
                            </Button>
                        </Tooltip>
                    </ButtonGroup>
                </div>
            </>}
        >
            <KeyModal
                kind={modal.kind}
                apiKey={modal.key}
                onClose={() => {
                    setModal({ kind: ModalKind.Closed, key: undefined })
                    fetchData()
                }}
            />

            <ScraperAuthFileModal
                apiKey={scraperAuthFileModal}
                onClose={() => setScraperAuthFileModal(undefined)}
            />

            <ScraperUsersModal
                apiKey={scraperAuthUsersModal}
                onClose={() => setScraperAuthUsersModal(undefined)}
            />

            {keys?.map(key => <KeyAccordionEntry
                key={key.id}
                apiKey={key}
                onEdit={() => setModal({ kind: ModalKind.Edit, key: key })}
                onDelete={() => setModal({ kind: ModalKind.Delete, key: key })}
                openAuthFileModal={() => setScraperAuthFileModal(key)}
                openScraperUsersModal={() => setScraperAuthUsersModal(key)}
            />)}
        </Card >
    )
}

interface KeyAccordionEntryParams {
    apiKey: ApiKey
    onEdit: () => void
    onDelete: () => void
    openAuthFileModal: () => void
    openScraperUsersModal: () => void
}

function KeyAccordionEntry(params: KeyAccordionEntryParams) {
    const { apiKey: key, onEdit, onDelete, openAuthFileModal, openScraperUsersModal } = params
    const isScraper = (key.roles & Roles.Scraper) === 0
    const noScraperRoleErr = "This key does not have the scraper role"

    return (
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
                        {key.name}
                    </h4>
                    <p className="domains">{key.domains.join(', ')}</p>
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
                <Tooltip title={isScraper ? noScraperRoleErr : "Manage a scraper's site login users"}>
                    <span>
                        <Button
                            disabled={isScraper}
                            onClick={openScraperUsersModal}
                        ><Password fontSize="small" style={{ marginRight: 6 }} />Login users</Button>
                    </span>
                </Tooltip>
                <Tooltip title={isScraper ? noScraperRoleErr : "Get a authentication file for this scraper"}>
                    <span>
                        <Button
                            disabled={isScraper}
                            onClick={openAuthFileModal}
                        ><VpnKey fontSize="small" style={{ marginRight: 6 }} />Auth File</Button>
                    </span>
                </Tooltip>
                <Button
                    onClick={onDelete}
                ><Delete fontSize="small" style={{ marginRight: 6 }} />Delete</Button>
                <Button
                    onClick={onEdit}
                ><Edit fontSize="small" style={{ marginRight: 6 }} />Edit</Button>
            </AccordionActions>
            <style jsx>{`
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
                .accordionSummary .domains {
                    flex-shrink: 1;
                    max-width: 230px;
                    overflow: hidden;
                    white-space: nowrap;
                    text-overflow: ellipsis;
                }
            `}</style>
        </Accordion>
    )
}
