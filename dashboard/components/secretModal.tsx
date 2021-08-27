import {
    Button,
    DialogContentText,
    FormControl,
    InputLabel,
    ListSubheader,
    Select,
    TextField,
    Tooltip,
    MenuItem,
    FormHelperText,
} from '@material-ui/core'
import FormatIndentIncrease from '@material-ui/icons/FormatIndentIncrease'
import RefreshIcon from '@material-ui/icons/Refresh'
import dynamic from 'next/dynamic'
import React, { useState, useEffect } from 'react'
import { fetcher } from '../src/auth'
import { Secret, ApiKey } from '../src/types'
import { randomString } from '../src/random'
import { Modal, ModalKind } from './modal'

const SyntaxHighlighter = dynamic(
    () => import('react-syntax-highlighter'),
    { ssr: false },
)

interface SecretModalProps {
    kind: ModalKind
    onClose: () => void
    secret?: Secret
    apiKeys?: Array<ApiKey>
}

let syntaxHighlighterStyle: any = undefined

export function SecretModal({ kind, onClose: onCloseArg, secret, apiKeys }: SecretModalProps) {
    const [decryptionKey, setDecryptionKey] = useState('')
    const [decryptionKeyError, setDecryptionKeyError] = useState('encryption key value must have a minimal length of 16 chars')
    const [apiError, setApiError] = useState('')
    const [key, setKey] = useState('')
    // If kind == 'create' this might contains a string value. If kind == 'view' this might contains the decrypted value as json so probably an array or object
    const [secretValue, setSecretValue] = useState(undefined as any)
    const [secretValueError, setSecretValueError] = useState('')
    const [keyId, setKeyId] = useState(undefined as undefined | string)

    const canSubmit = kind == ModalKind.Create ? !secretValueError && !decryptionKeyError && keyId && key && decryptionKey : true

    const onClose = () => {
        setDecryptionKey('')
        setDecryptionKeyError('encryption key value must have a minimal length of 16 chars')
        setApiError('')
        setKey('')
        setSecretValueError('')
        setSecretValue(undefined)
        onCloseArg()
    }

    useEffect(() => {
        if (kind == ModalKind.Create && typeof secretValue != 'string')
            setSecretValue('{}')
        if (kind == ModalKind.Create && keyId === undefined && apiKeys)
            setKeyId(apiKeys.filter(key => key.enabled)[0].id)

        if (kind != ModalKind.Closed && syntaxHighlighterStyle === undefined) {
            // When the modal is opend start pre-loading the highlighter
            import('react-syntax-highlighter')
            // Also loadin the highlighter style
            import('react-syntax-highlighter/dist/esm/styles/hljs').then(v => syntaxHighlighterStyle = v.monokaiSublime)
        }
    }, [kind])

    const submit = async () => {
        try {
            switch (kind) {
                case ModalKind.View:
                    setSecretValue(
                        await fetcher.get(
                            `/api/v1/secrets/otherKey/${secret?.keyId}/${secret?.key}/${decryptionKey}`,
                        ),
                    )
                    break
                case ModalKind.Create:
                    await fetcher.post(
                        `/api/v1/secrets/otherKey/${keyId}/${key}/${decryptionKey}`,
                        JSON.parse(secretValue),
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
        } catch (e) {
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
            else if (kind == ModalKind.Create)
                return (<div>
                    <DialogContentText>
                        Create a new secret value.<br />
                        The Encryption / Decryption is not stored by the server, it's only used by the server to encrypt / decrypt the send value<br />
                        This also means you should store the encryption key used here somewhere safe
                    </DialogContentText>

                    <TextField
                        id='key'
                        label='Key'
                        value={key}
                        onChange={(e) => setKey(e.target.value)}
                        variant="filled"
                        fullWidth
                    />

                    <div className="inputWithButton marginTop">
                        <TextField
                            id="secret"
                            label="Encryption key"
                            value={decryptionKey}
                            error={!!decryptionKeyError}
                            helperText={decryptionKeyError}
                            onChange={(e) => {
                                const { value } = e.target
                                setDecryptionKey(value)
                                setDecryptionKeyError(
                                    value.length < 16
                                        ? 'encryption key value must have a minimal length of 16 chars'
                                        : ''
                                )
                            }}
                            variant="filled"
                            fullWidth
                        />
                        <div className="toggles">
                            <Tooltip title='Generate random value'>
                                <Button onClick={() => setDecryptionKey(randomString(32))}>
                                    <RefreshIcon fontSize="small" />
                                </Button>
                            </Tooltip>

                        </div>
                    </div>
                    <div className="inputWithButton secretValueInputArea marginTop">
                        <TextField
                            className="secretModalSecretValueInput"
                            id="secret-value"
                            label="JSON Value"
                            value={secretValue}
                            helperText={secretValueError || 'json value is valid'}
                            error={!!secretValueError}
                            onChange={(e) => {
                                try {
                                    const { value } = e.target
                                    setSecretValue(value)
                                    JSON.parse(value)
                                    setSecretValueError('')
                                } catch (e) {
                                    setSecretValueError(e.message)
                                }
                            }}
                            variant="filled"
                            multiline
                            fullWidth
                        />
                        <div className="toggles">
                            <Tooltip title='Format json'>
                                <Button
                                    disabled={!!secretValueError}
                                    onClick={() => setSecretValue((prev: string) => JSON.stringify(JSON.parse(prev), null, 2))}
                                >
                                    <FormatIndentIncrease fontSize="small" />
                                </Button>
                            </Tooltip>
                        </div>
                    </div>
                    <div className="marginTop">
                        <FormControl fullWidth variant="filled">
                            <InputLabel htmlFor="secret-key-id">Api key</InputLabel>
                            <Select
                                value={keyId}
                                onChange={(id: any) => setKeyId(id.target.value)}
                                id="secret-key-id"
                            >
                                {apiKeys?.filter(key => key.enabled).reduce((acc: Array<any>, key: ApiKey) => {
                                    return [
                                        ...acc,
                                        <ListSubheader key={key.id + '-header'}>{key.domains.join(', ')}</ListSubheader>,
                                        <MenuItem key={key.id + '-selectable'} value={key.id}>{key.id}</MenuItem>,
                                    ]
                                }, [])}
                            </Select>
                            <FormHelperText>The API key selected will be able to access the secret</FormHelperText>
                        </FormControl>
                    </div>
                    <style jsx>{`
                        .inputWithButton {
                            display: flex;
                            justify-content: space-between;
                            align-items: center;
                        }
                        .inputWithButton .toggles {
                            margin-left: 10px;
                        }
                        .marginTop {
                            margin-top: 10px;
                        }
                        .secretValueInputArea {
                            align-items: flex-start;
                        }
                    `}</style>
                    <style global jsx>{`
                        .secretModalSecretValueInput textarea, .secretModalSecretValueInput input {
                            font-family:Consolas,Monaco,Lucida Console,Liberation Mono,DejaVu Sans Mono,Bitstream Vera Sans Mono,Courier New, monospace;
                        }
                    `}</style>
                </div>)
            else
                return (<></>)
        }}</Modal>
    )
}
