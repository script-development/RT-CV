import {
    DialogContentText,
    TextField,
    Button,
    Breadcrumbs,
} from '@mui/material'
import Edit from '@mui/icons-material/Edit'
import dynamic from 'next/dynamic'
import React, { useState, useEffect, useMemo } from 'react'
import { fetcher, getKeys } from '../../src/auth'
import { Modal, ModalKind } from '../modal'
import ModifyOrCreateSecret from './modify'
import { SecretValueStructure, ApiKey } from '../../src/types'
import { SecretModalProps } from './secretModalProps'

// Use dynamic so we only load the bytes of this component once we really need it
// It'q quite large and loading it this way splits it from the main bundle
// So that means the first load of the page stays quick
const JSONCode = dynamic(() => import('../jsonCode'), { ssr: false })

export interface ModifyState {
    id: string
    key: string
    description: string
    valueStructure: SecretValueStructure | undefined
    value: string
    valueError: string
    selectedKeyId: string
    decryptionKey: string
    decryptionKeyError: string
}

export function SecretModal({ kind, setKind, onClose: onCloseArg, secret }: SecretModalProps) {
    const [apiError, setApiError] = useState('')
    const [apiKeys, setApiKeys] = useState<undefined | Array<ApiKey>>(undefined)
    const [modifyState, setModifyState] = useState<ModifyState>({
        id: '',
        key: '',
        description: '',
        valueStructure: undefined,
        value: '',
        valueError: '',
        selectedKeyId: '',
        decryptionKey: '',
        decryptionKeyError: 'encryption key value must have a minimal length of 16 chars',
    })
    const [viewState, setViewState] = useState({ value: undefined as any, decryptionKey: '' })

    const canSubmit = (kind == ModalKind.Create || kind == ModalKind.Edit) ?
        !modifyState.valueError
        && !modifyState.decryptionKeyError
        && modifyState.selectedKeyId
        && modifyState.key
        && modifyState.decryptionKey
        && modifyState.valueStructure !== undefined
        : true

    const onClose = () => {
        setApiError('')
        setModifyState({
            id: '',
            key: '',
            description: '',
            valueStructure: undefined,
            value: '',
            valueError: '',
            selectedKeyId: '',
            decryptionKey: '',
            decryptionKeyError: 'encryption key value must have a minimal length of 16 chars',
        })
        setViewState({ value: undefined, decryptionKey: '' })
        onCloseArg()
    }

    const secretApiKey = useMemo(() => apiKeys?.find(key => key.id == secret?.keyId), [apiKeys, secret])

    useEffect(() => {
        getKeys().then(keys => setApiKeys(keys))
    }, [])

    useEffect(() => {
        if (kind == ModalKind.View)
            // When the modal is opened start pre-loading the highlighter
            import('../jsonCode')

        if (kind == ModalKind.Edit && secret?.id != modifyState.id)
            setModifyState(s => ({
                ...s,
                id: secret?.id || '',
                key: secret?.key || '',
                selectedKeyId: secret?.keyId || '',
                description: secret?.description || '',
                valueStructure: secret?.valueStructure,
                value: JSON.stringify(viewState.value),
                decryptionKey: viewState.decryptionKey,
                decryptionKeyError: '',
            }))
    }, [kind, secret])

    const submit = async () => {
        try {
            switch (kind) {
                case ModalKind.View:
                    const value = await fetcher.get(
                        `/api/v1/secrets/otherKey/${secret?.keyId}/${secret?.key}/${viewState.decryptionKey}`,
                    )
                    setViewState(s => ({ ...s, value }))
                    break
                case ModalKind.Edit:
                case ModalKind.Create:
                    await fetcher.put(
                        `/api/v1/secrets/otherKey/${modifyState.selectedKeyId}/${modifyState.key}`,
                        {
                            value: JSON.parse(modifyState.value),
                            valueStructure: modifyState.valueStructure,
                            description: modifyState.description,
                            encryptionKey: modifyState.decryptionKey,
                        },
                    )
                    onClose()
                    break
                case ModalKind.Delete:
                    await fetcher.delete(
                        `/api/v1/secrets/otherKey/${secret?.keyId}/${secret?.key}`,
                    )
                    onClose()
                    break
            }
        } catch (e: any) {
            setApiError(e?.message || e)
        }
    }

    return (
        <Modal
            kind={kind}
            onClose={onClose}
            onSubmit={submit}
            title={{
                create: 'Create Secret',
                view: 'View Secret',
                edit: 'Edit Secret',
                delete: 'Delete Secret',
            }}
            confirmText={{
                create: 'Create',
                delete: 'Delete',
                edit: 'Save',
                view: 'Decrypt',
            }}
            showConfirm={kind != ModalKind.View || viewState.value === undefined}
            submitDisabled={!canSubmit}
            apiError={apiError}
            setApiError={setApiError}
            fullWidth={true}
            maxWidth={'md'}
        >{(kind: ModalKind) => {
            if (kind == ModalKind.Delete)
                return (
                    <DialogContentText>
                        Are you sure you want to delete this secret?
                    </DialogContentText>
                )
            else if (kind == ModalKind.View)
                if (viewState.value)
                    return (
                        <div>
                            <div className="info">
                                <Breadcrumbs>
                                    <p>{secretApiKey?.name || secret?.keyId}</p>
                                    <b style={{ color: "white" }}>{secret?.key}</b>
                                </Breadcrumbs>
                                {secret?.description ? <DialogContentText>{secret?.description}</DialogContentText> : undefined}
                            </div>

                            <Button
                                variant="outlined"
                                startIcon={<Edit />}
                                onClick={() => setKind(ModalKind.Edit)}
                            >Edit</Button>
                            <JSONCode json={viewState.value} />
                            <style jsx>{`
                                .info {
                                    margin-bottom: 10px;
                                }
                                .code {
                                    margin-top: 10px;
                                    overflow: hidden;
                                    border-radius: 4px;
                                }
                            `}</style>
                        </div>
                    )
                else
                    return (<div>
                        <DialogContentText>
                            Fill in the decryption key to continue
                        </DialogContentText>
                        <TextField
                            id="secret"
                            label="Decryption key"
                            value={viewState.decryptionKey}
                            onChange={e => setViewState(s => ({ ...s, decryptionKey: e.target.value }))}
                            variant="filled"
                            fullWidth
                        />
                    </div>)
            else if ((kind == ModalKind.Create || kind == ModalKind.Edit) && apiKeys !== undefined)
                return (
                    <ModifyOrCreateSecret
                        state={modifyState}
                        setState={setModifyState}
                        create={kind == ModalKind.Create}
                        apiKeys={apiKeys}
                    />
                )
            else
                return (<></>)
        }}</Modal >
    )
}
