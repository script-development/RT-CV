import {
    DialogContentText,
    TextField,
} from '@material-ui/core'
import React, { useState } from 'react'
import { fetcher } from '../src/auth'
import { Secret } from '../src/types'
import { Modal, ModalKind } from './modal'
import dynamic from 'next/dynamic'
import { useEffect } from 'react'

const SyntaxHighlighter = dynamic(
    () => import('react-syntax-highlighter'),
    { ssr: false },
)
// import { docco } from 'react-syntax-highlighter/dist/esm/styles/hljs';

interface SecretModalProps {
    kind: ModalKind
    onClose: () => void
    secret?: Secret
}

export function SecretModal({ kind, onClose: onCloseArg, secret }: SecretModalProps) {
    const [decryptionKey, setDecryptionKey] = useState('')
    const [apiError, setApiError] = useState('')
    const [secretValue, setSecretValue] = useState(undefined as any)
    const [syntaxHighlighterStyle, setSyntaxHighlighterStyle] = useState(undefined as any)

    const canSubmit = true

    const onClose = () => {
        setSecretValue(undefined)
        setDecryptionKey('')
        onCloseArg()
    }

    useEffect(() => {
        if (kind != ModalKind.Closed && syntaxHighlighterStyle === undefined) {
            // When the modal is opend start pre-loading the highlighter
            import('react-syntax-highlighter')
            // Also loadin the highlighter style
            import('react-syntax-highlighter/dist/esm/styles/hljs').then(v => setSyntaxHighlighterStyle(v.monokaiSublime))
        }
    }, [kind])

    const submit = async () => {
        try {
            switch (kind) {
                case ModalKind.View:
                    setSecretValue(await fetcher.get(`/api/v1/secrets/otherKey/${secret?.keyId}/${secret?.key}/${decryptionKey}`))
                    break
                case ModalKind.Create:
                    // TODO
                    onClose()
                    break
                case ModalKind.Delete:
                    // TODO
                    onClose()
                    break
                case ModalKind.Edit:
                    // TODO
                    onClose()
                    break
            }
        } catch (e) {
            setApiError(e?.message || e)
        }
    }

    const tryDecrypt = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
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
            showConfirm={kind != ModalKind.View || secretValue === undefined}
            submitDisabled={!canSubmit}
            apiError={apiError}
            setApiError={setApiError}
            fullWidth={true}
        >{(kind: ModalKind) => {
            if (kind == ModalKind.Delete)
                return (<DialogContentText>
                    Are you sure you want to delete this secret?
                </DialogContentText>)
            else if (kind == ModalKind.View)
                if (secretValue)
                    if (syntaxHighlighterStyle)
                        return (<div className="code">
                            <SyntaxHighlighter
                                wrapLongLines={true}
                                language="json"
                                style={syntaxHighlighterStyle}
                            >{JSON.stringify(secretValue, null, 4)}</SyntaxHighlighter>

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
                            Fillin the decryption key to continue
                        </DialogContentText>
                        <TextField
                            id="secret"
                            label="Decryption key"
                            value={decryptionKey}
                            onChange={(e) => setDecryptionKey(e.target.value)}
                            variant="filled"
                            fullWidth
                        />
                    </div>)
            else
                return (<div>
                    {/* TODO */}
                </div>)
        }}</Modal>
    )
}
