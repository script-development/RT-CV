import {
    DialogContentText,
    TextField,
    Button,
    FormControl,
    FormLabel,
    RadioGroup,
    FormControlLabel,
    Radio,
    ButtonGroup,
    Checkbox,
} from '@mui/material'
import { DataGrid } from '@mui/x-data-grid'
import React, { useState, useEffect } from 'react'
import { Modal, ModalKind } from '../modal'
import { SecretValueStructure } from '../../src/types'
import { ModalProps } from './props'
import { fetcher } from '../../src/auth'

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

export function SecretModal({ kind, onClose: onCloseArg, hook }: ModalProps) {
    const [apiError, setApiError] = useState('')
    const [headers, setHeaders] = useState<Array<{ key: string, value: string }>>([])
    const [method, setMethod] = useState('POST')
    const [url, setUrl] = useState('https://')

    const canSubmit = url.length != 0
    const viewState = { value: undefined }

    const headersTableColumns = [
        { flex: 1, editable: true, field: 'key', headerName: 'Key' },
        { flex: 2, editable: true, field: 'value', headerName: 'Value' },
    ]
    const [selectedHeaders, setSelectedHeaders] = useState<Array<number>>([])

    const addHeader = () => setHeaders(v => [...v, { key: '', value: '' }])
    const removeHeader = () => {
        const selectedHeadersToRemove = selectedHeaders.sort().reverse()
        setHeaders(v => {
            const newList = [...v]
            selectedHeadersToRemove.map(idx => newList.splice(idx, 1))
            return newList
        })
        setSelectedHeaders([])
    }
    const formatHeadersToApi = () => headers.filter(h => h.key.length > 0).map(h => ({ key: h.key, value: [h.value] }))

    const onClose = () => {
        setApiError('')
        setHeaders([])
        setMethod('POST')
        setUrl('https://')
        onCloseArg()
    }

    const onSubmit = async () => {
        try {
            switch (kind) {
                case ModalKind.Create:
                    await fetcher.post(`/api/v1/onMatchHooks`, {
                        method,
                        url,
                        addHeaders: formatHeadersToApi(),
                    })
                    onClose()
                    break
                case ModalKind.Edit:
                    await fetcher.put(`/api/v1/onMatchHooks/${hook?.id}`, {
                        method,
                        url,
                        addHeaders: formatHeadersToApi(),
                    })
                    onClose()
                    break
                case ModalKind.Delete:
                    await fetcher.delete(`/api/v1/onMatchHooks/${hook?.id}`)
                    onClose()
                    break
                default:
                    throw 'TODO'
            }

        } catch (e: any) {
            setApiError(e?.message || e)
        }
    }

    useEffect(() => {
        if (kind == ModalKind.Edit && hook) {
            setHeaders((hook.addHeaders || []).map(h => ({ key: h.key, value: h.value.join(',') })) || [])
            setMethod(hook.method)
            setUrl(hook.url)
        }
    }, [kind, hook])

    return (
        <Modal
            kind={kind}
            onClose={onClose}
            onSubmit={onSubmit}
            title={{
                create: 'Create a On match hook',
                view: 'View on Match hook',
                edit: 'Edit on match hook',
                delete: 'Delete on Match hook',
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
            maxWidth={'md'}
            fullWidth
        >{(kind: ModalKind) => {
            if (kind == ModalKind.Delete)
                return (
                    <DialogContentText>
                        Are you sure you want to delete this hook?
                    </DialogContentText>
                )
            else if (kind == ModalKind.View)
                if (viewState.value)
                    return (
                        <></>
                    )
                else
                    return (
                        <></>
                    )
            else if (kind == ModalKind.Edit || kind == ModalKind.Create)
                return (
                    <div>
                        <FormControl component="fieldset" fullWidth>
                            <FormLabel component="legend">Method</FormLabel>
                            <RadioGroup aria-label="method" name="method" value={method} onChange={e => setMethod(e.target.value)}>
                                <FormControlLabel value="GET" control={<Radio />} label="Get" />
                                <FormControlLabel value="POST" control={<Radio />} label="Post" />
                                <FormControlLabel value="PUT" control={<Radio />} label="Put" />
                                <FormControlLabel value="PATCH" control={<Radio />} label="Patch" />
                                <FormControlLabel value="DELETE" control={<Radio />} label="Delete" />
                            </RadioGroup>
                        </FormControl>

                        <div className='padTop'>
                            <TextField
                                value={url}
                                autoFocus
                                onChange={e => setUrl(e.target.value)}
                                id="url"
                                label="URL"
                                variant="filled"
                                fullWidth
                            />
                        </div>

                        <div className='padTop'>
                            <DialogContentText>Add custom headers</DialogContentText>
                            <DataGrid
                                rows={headers.map((h, idx) => ({ ...h, id: idx }))}
                                columns={headersTableColumns}
                                pageSize={100}
                                onCellEditCommit={e => setHeaders(current => {
                                    const field = e.field as 'key' | 'value'
                                    const id = e.id as number

                                    const newHeaders = [...current]
                                    newHeaders[id][field] = typeof e.value == 'string' ? e.value : ''
                                    return newHeaders
                                })}
                                checkboxSelection={true}
                                selectionModel={selectedHeaders}
                                onSelectionModelChange={e => setSelectedHeaders(e as Array<number>)}
                                autoHeight
                                disableSelectionOnClick
                            />
                        </div>
                        <div className='padTop'>
                            <ButtonGroup color="primary" variant="contained">
                                <Button onClick={addHeader}>Add header</Button>
                                <Button onClick={removeHeader} disabled={selectedHeaders.length == 0}>Remove selected</Button>
                            </ButtonGroup>
                        </div>

                        <style jsx>{`
                            .padTop {
                                padding-top: 10px;
                            }
                        `}</style>
                    </div>
                )
            else
                return (<></>)
        }}</Modal >
    )
}
