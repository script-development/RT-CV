import {
    DialogContentText,
    TextField,
} from '@material-ui/core'
import dynamic from 'next/dynamic'
import React, { useState, useEffect } from 'react'
import { fetcher } from '../../src/auth'
import { Modal, ModalKind } from '../modal'
import Create from './create'
import { ValueKind } from './modifyValue'
import { SecretModalProps } from './secretModalProps'

const SyntaxHighlighter = dynamic(
    () => import('react-syntax-highlighter'),
    {
        ssr: false,
    },
)

let syntaxHighlighterStyle: any = undefined

export interface ModifyState {
    key: string
    description: string
    valueKind: ValueKind
    value: string
    valueError: string
    selectedKeyId: string
    decryptionKey: string
    decryptionKeyError: string
}

export function SecretModal({ kind, onClose: onCloseArg, secret }: SecretModalProps) {
    const [apiError, setApiError] = useState('')

    const [modifyState, setModifyState] = useState<ModifyState>({
        key: '',
        description: '',
        valueKind: undefined,
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

    const canSubmit = kind == ModalKind.Create ?
        !modifyState.valueError && !modifyState.decryptionKeyError && modifyState.selectedKeyId && modifyState.key && modifyState.decryptionKey
        : true

    const onClose = () => {
        setApiError('')
        setModifyState({
            key: '',
            description: '',
            valueKind: undefined,
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
            syntaxHighlighterStyle = styles.monokaiSublime
        } catch (e: any) {
            console.log(e)
        }
    }

    useEffect(() => {
        if (kind != ModalKind.Closed && syntaxHighlighterStyle === undefined)
            // When the modal is opened start pre-loading the highlighter
            loadReactSyntaxHighlighter()
    }, [kind])

    const submit = async () => {
        try {
            switch (kind) {
                case ModalKind.View:
                    const value = await fetcher.get(
                        `/api/v1/secrets/otherKey/${secret?.keyId}/${secret?.key}/${viewState.decryptionKey}`,
                    )
                    setViewState(s => ({ ...s, value }))
                    break
                case ModalKind.Create:
                    await fetcher.post(
                        `/api/v1/secrets/otherKey/${modifyState.selectedKeyId}/${modifyState.key}/${modifyState.decryptionKey}`,
                        {
                            value: JSON.parse(modifyState.value),
                            description: modifyState.description,
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
                case ModalKind.Edit:
                    // TODO
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
                    if (syntaxHighlighterStyle)
                        return (<div className="code">
                            <SyntaxHighlighter
                                wrapLongLines={true}
                                language="json"
                                style={syntaxHighlighterStyle}
                            >{JSON.stringify(viewState.value, null, 4)}</SyntaxHighlighter>

                            <style jsx>{`
                                .code {
                                    overflow: hidden;
                                    border-radius: 4px;
                                }
                            `}</style>
                        </div>)
                    else
                        return (<>Loading...</>)
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
            else if (kind == ModalKind.Create)
                return (
                    <Create
                        state={modifyState}
                        setState={setModifyState}
                    />
                )
            else
                return (<></>)
        }}</Modal >
    )
}
