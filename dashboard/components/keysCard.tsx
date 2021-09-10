import {
    AccordionSummary,
    Accordion,
    AccordionActions,
    AccordionDetails,
    Button,
    ButtonGroup,
    Divider,
    Tooltip,
} from '@material-ui/core'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import Edit from '@material-ui/icons/Edit'
import Add from '@material-ui/icons/Add'
import Delete from '@material-ui/icons/Delete'
import React, { useEffect, useState } from 'react'
import Card from './card'
import { getKeys } from '../src/auth'
import { ApiKey } from '../src/types'
import { ModalKind } from './modal'
import { KeyModal } from './keyModal'

export default function KeysCard() {
    const [loading, setLoading] = useState(true)
    const [keys, setKeys] = useState(undefined as Array<ApiKey> | undefined)
    const [modal, setModal] = useState({ kind: ModalKind.Closed, key: undefined as (undefined | ApiKey) })

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
                            onClick={() => setModal({ kind: ModalKind.Delete, key: key })}
                        ><Delete fontSize="small" style={{ marginRight: 6 }} />Delete</Button>
                        <Button
                            onClick={() => setModal({ kind: ModalKind.Edit, key: key })}
                        ><Edit fontSize="small" style={{ marginRight: 6 }} />Edit</Button>
                    </AccordionActions>
                </Accordion>
            )}

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
            `}</style>
        </Card>
    )
}
