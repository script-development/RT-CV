import { Breadcrumbs, Button, ButtonGroup, Tooltip } from '@material-ui/core'
import Add from '@material-ui/icons/Add'
import Delete from '@material-ui/icons/Delete'
import Visibility from '@material-ui/icons/Visibility'
import Person from '@material-ui/icons/Person'
import People from '@material-ui/icons/People'
import Code from '@material-ui/icons/Code'
import React, { useState, useEffect, Dispatch, SetStateAction } from 'react'

import { fetcher, getKeys } from '../src/auth'
import { SecretValueStructure, Secret, ApiKey } from '../src/types'
import Card from './card'
import { ModalKind } from './modal'
import Dynamic from 'next/dynamic'
import { SecretModalProps } from './secretModal/secretModalProps'

const SecretModal = Dynamic<SecretModalProps>(() => import('./secretModal/secretModal').then(m => m.SecretModal))

interface SecretsCardArgs { }

export default function SecretsCard({ }: SecretsCardArgs) {
    const [loading, setLoading] = useState(true)
    const [secrets, setSecrets] = useState(undefined as Array<Secret> | undefined)
    const [keys, setKeys] = useState(undefined as Array<ApiKey> | undefined)
    const [modal, setModal] = useState({ kind: ModalKind.Closed, secret: undefined as (undefined | Secret) })

    const fetchData = async () => {
        try {
            setLoading(true)
            const [secrets, keys] = await Promise.all([
                fetcher.fetch(`/api/v1/secrets/otherKey`),
                getKeys(),
            ])
            setSecrets(secrets)
            setKeys(keys)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchData() }, [])

    const GetSecretApiKey = (secret: Secret) =>
        keys?.find(k => k.id == secret.keyId)

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
                setKind={kind => setModal(s => ({ ...s, kind }))}
                secret={modal.secret}
                onClose={() => {
                    setModal({ kind: ModalKind.Closed, secret: undefined })
                    fetchData()
                }}
            />

            {secrets?.map((secret, idx) => {
                const apiKey = GetSecretApiKey(secret)
                return (<Secret
                    key={idx}
                    secret={secret}
                    apiKey={apiKey}
                    isLastRow={secrets.length == (idx + 1)}
                    openModal={kind => setModal({ kind, secret })}
                />)
            })}
        </Card>
    )
}

interface SecretProps {
    secret: Secret
    apiKey: ApiKey | undefined
    isLastRow: boolean
    openModal: (kind: ModalKind) => void
}

function Secret({ secret, apiKey, isLastRow, openModal }: SecretProps) {
    return (
        <div key={secret.id} className={"simpleRow" + (isLastRow ? ' last' : '')}>
            <div className="side">
                <div>
                    {
                        secret.valueStructure == SecretValueStructure.StrictUser
                            ? <Tooltip title="Contains a single user"><Person fontSize='small' /></Tooltip>
                            : secret.valueStructure == SecretValueStructure.StrictUsers
                                ? <Tooltip title="Contains multiple user"><People fontSize='small' /></Tooltip>
                                : <Tooltip title="Contains unknown json value"><Code fontSize='small' /></Tooltip>
                    }
                </div>
                <div className="secretIdentifier">
                    <Breadcrumbs>
                        <p>{apiKey?.name || secret.keyId}</p>
                        <b style={{ color: "white" }}>{secret.key}</b>
                    </Breadcrumbs>
                    {secret.description ? <p className="description">{secret.description}</p> : undefined}
                </div>
            </div>
            <div className="side">
                <Tooltip title="View secret contents">
                    <Button onClick={() => openModal(ModalKind.View)}>
                        <Visibility fontSize="small" />
                    </Button>
                </Tooltip>
                <Tooltip title="Delete secret">
                    <Button onClick={() => openModal(ModalKind.Delete)}>
                        <Delete fontSize="small" />
                    </Button>
                </Tooltip>
            </div>

            <style jsx>{`
                .simpleRow {
					border-top: 1px solid rgba(255, 255, 255, 0.12);
					background-color: #424242;
					padding: 10px;
                    display: flex;
					justify-content: space-between;
					align-items: center;
                }
                .simpleRow .side:first-child {
					display: flex;
                    align-items: center;
				}
				.simpleRow.last {
					border-bottom-left-radius: 4px;
					border-bottom-right-radius: 4px;
				}
                .secretIdentifier {
                    padding-left: 15px;
                }
				.secretIdentifier .description {
					color: rgba(255, 255, 255, .8);
				}
            `}</style>
        </div>
    )
}
