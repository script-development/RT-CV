import { Breadcrumbs, Button, ButtonGroup, Tooltip } from '@material-ui/core'
import Add from '@material-ui/icons/Add'
import Delete from '@material-ui/icons/Delete'
import Visibility from '@material-ui/icons/Visibility'
import React, { useState, useEffect } from 'react'
import { fetcher } from '../src/auth'
import { Secret } from '../src/types'
import Card from './card'
import { ModalKind } from './modal'
import Dynamic from 'next/dynamic'
import { SecretModalProps } from './secretModal/secretModalProps'

const SecretModal = Dynamic<SecretModalProps>(() => import('./secretModal/secretModal').then(m => m.SecretModal))

interface SecretsCardArgs { }

export default function SecretsCard({ }: SecretsCardArgs) {
    const [loading, setLoading] = useState(true)
    const [secrets, setSecrets] = useState(undefined as Array<Secret> | undefined)
    const [modal, setModal] = useState({ kind: ModalKind.Closed, secret: undefined as (undefined | Secret) })

    const fetchData = async () => {
        try {
            setLoading(true)
            const secrets = await fetcher.fetch(`/api/v1/secrets/otherKey`)
            setSecrets(secrets)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchData() }, [])

    return (
        <Card
            title="Secrets"
            loading={loading}
            headerContent={<>
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
                                onClick={() => setModal({ kind: ModalKind.Create, secret: undefined })}
                            >
                                <Add fontSize={'small'} />
                            </Button>
                        </Tooltip>
                    </ButtonGroup>
                </div>
            </>}
        >
            <SecretModal
                kind={modal.kind}
                secret={modal.secret}
                onClose={() => {
                    setModal({ kind: ModalKind.Closed, secret: undefined })
                    fetchData()
                }}
            />

            {secrets?.map((secret, idx) =>
                <div key={secret.id} className={"simpleRow" + (secrets.length == (idx + 1) ? ' last' : '')}>
                    <div className="secretIdentifier">
                        <Breadcrumbs>
                            <p>{secret.keyId}</p>
                            <b style={{ color: "white" }}>{secret.key}</b>
                        </Breadcrumbs>
                        {secret.description ? <p className="description">{secret.description}</p> : undefined}
                    </div>
                    <div>
                        <Tooltip title="View secret contents">
                            <Button onClick={() => setModal({ kind: ModalKind.View, secret: secret })}>
                                <Visibility fontSize="small" />
                            </Button>
                        </Tooltip>
                        <Tooltip title="Delete secret">
                            <Button onClick={() => setModal({ kind: ModalKind.Delete, secret: secret })}>
                                <Delete fontSize="small" />
                            </Button>
                        </Tooltip>
                    </div>
                </div>
            )}

            <style jsx>{`
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
				.secretIdentifier .description {
					color: rgba(255, 255, 255, .8);
				}
            `}</style>
        </Card>
    )
}
