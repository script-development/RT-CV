import { Breadcrumbs, Button, ButtonGroup, Tooltip } from '@mui/material'
import Add from '@mui/icons-material/Add'
import Delete from '@mui/icons-material/Delete'
import Edit from '@mui/icons-material/Edit'
import React, { useState, useEffect } from 'react'

import { fetcher } from '../src/auth'
import { OnMatchHook } from '../src/types'
import Card from './card'
import { ModalKind } from './modal'
import Dynamic from 'next/dynamic'
import { ModalProps } from './onMatchHookModal/props'
import { PlayArrow } from '@mui/icons-material'

const Modal = Dynamic<ModalProps>(() => import('./onMatchHookModal/modal').then(m => m.SecretModal))

interface OnMatchHooksCardArgs { }

export default function OnMatchHooksCard({ }: OnMatchHooksCardArgs) {
    const [loading, setLoading] = useState(true)
    const [onMatchHooks, setOnMatchHooks] = useState(undefined as Array<OnMatchHook> | undefined)
    const [modal, setModal] = useState({ kind: ModalKind.Closed, value: undefined as (undefined | OnMatchHook) })

    const fetchData = async () => {
        try {
            setLoading(true)
            const hooks = await fetcher.fetch(`/api/v1/onMatchHooks`)
            setOnMatchHooks(hooks)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchData() }, [])

    return (
        <Card
            title="On Match Hooks"
            loading={loading}
            headerContent={<>
                <div>
                    {loading
                        ? 'loading..'
                        : onMatchHooks?.length
                            ? `${onMatchHooks?.length} on match hooks`
                            : 'No on match hooks'
                    }
                </div>
                <div>
                    <ButtonGroup color="primary" variant="contained">
                        <Tooltip title="Create secret">
                            <Button
                                onClick={() => setModal({ kind: ModalKind.Create, value: undefined })}
                            >
                                <Add fontSize={'small'} />
                            </Button>
                        </Tooltip>
                    </ButtonGroup>
                </div>
            </>}
        >
            <Modal
                kind={modal.kind}
                hook={modal.value}
                onClose={() => {
                    setModal({ kind: ModalKind.Closed, value: undefined })
                    fetchData()
                }}
            />

            {onMatchHooks?.map((value, idx) => {
                return (<Line
                    key={value.id}
                    value={value}
                    isLastRow={onMatchHooks.length - 1 == idx}
                    openModal={kind => setModal({ kind, value })}
                />)
            })}
        </Card>
    )
}

interface LineProps {
    value: OnMatchHook
    isLastRow: boolean
    openModal: (kind: ModalKind) => void
}

function Line({ value, isLastRow, openModal }: LineProps) {
    const testOnMatchHook = () => fetcher.post(`/api/v1/onMatchHooks/${value.id}/test`)

    return (
        <div key={value.id} className={"simpleRow" + (isLastRow ? ' last' : '')}>
            <div className="side">
                <b style={{ color: "white" }}>{value.method}</b>
                <span style={{ paddingLeft: 5 }}>{value.url}</span>
            </div>
            <div className="side">
                <Tooltip title="Test on match hook">
                    <Button onClick={testOnMatchHook}>
                        <PlayArrow fontSize="small" />
                    </Button>
                </Tooltip>
                <Tooltip title="Edit on match hook">
                    <Button onClick={() => openModal(ModalKind.Edit)}>
                        <Edit fontSize="small" />
                    </Button>
                </Tooltip>
                <Tooltip title="Delete on match hook">
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
            `}</style>
        </div>
    )
}
