import {
    DialogContentText,
    TextField,
    Button,
} from '@material-ui/core'
import Edit from '@material-ui/icons/Edit'
import dynamic from 'next/dynamic'
import React, { useState, useEffect } from 'react'
import { fetcher } from '../../src/auth'
import { Modal, ModalKind } from '../modal'
import Create from './modify'
import { SecretValueStructure } from '../../src/types'
import { SecretModalProps } from './secretModalProps'

const SyntaxHighlighter = dynamic(
    () => import('react-syntax-highlighter'),
    { ssr: false },
)

let syntaxHighlighterStyleCache: any = undefined

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

    const [syntaxHighlighterStyle, setSyntaxHighlighterStyle] = useState<any>(syntaxHighlighterStyleCache)
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

    const [viewState, setViewState] = useState({
        value: undefined as any,
        decryptionKey: '',
    })

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

    const loadReactSyntaxHighlighter = async () => {
        try {
            const [_, styles] = await Promise.all([
                // Pre load the highlighter component
                import('react-syntax-highlighter'),
                // Load the styles
                import('react-syntax-highlighter/dist/esm/styles/hljs'),
            ])

            syntaxHighlighterStyleCache = styles.monokaiSublime
            setSyntaxHighlighterStyle(styles.monokaiSublime)
        } catch (e: any) {
            console.log(e)
        }
    }

    useEffect(() => {
        if (kind == ModalKind.View && syntaxHighlighterStyle === undefined)
            // When the modal is opened start pre-loading the highlighter
            loadReactSyntaxHighlighter()

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
                            <Button
                                variant="outlined"
                                startIcon={<Edit />}
                                onClick={() => setKind(ModalKind.Edit)}
                            >Edit</Button>
                            <div className="code">
                                {syntaxHighlighterStyle
                                    ? <SyntaxHighlighter
                                        wrapLongLines={true}
                                        language="json"
                                        style={syntaxHighlighterStyle}
                                    >{JSON.stringify(viewState.value, null, 4)}</SyntaxHighlighter>
                                    : 'Loading...'
                                }
                            </div>
                            <style jsx>{`
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
            else if (kind == ModalKind.Create || kind == ModalKind.Edit)
                return (
                    <Create
                        state={modifyState}
                        setState={setModifyState}
                        create={kind == ModalKind.Create}
                    />
                )
            else
                return (<></>)
        }}</Modal >
    )
}
