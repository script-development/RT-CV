import {
    Button,
    DialogContentText,
    FormControl,
    InputLabel,
    Select,
    TextField,
    Tooltip,
    MenuItem,
    FormHelperText,
} from '@mui/material'
import { randomString } from '../../src/random'
import RefreshIcon from '@mui/icons-material/Refresh'
import { ModifyState } from './modal'
import React, { Dispatch, SetStateAction, useEffect, useState } from 'react'
import { ApiKey } from '../../src/types'
import ModifyValue from './modifyValue'

interface ModifyProps {
    state: ModifyState
    setState: Dispatch<SetStateAction<ModifyState>>
    create?: boolean
    apiKeys: Array<ApiKey>
}

export default function Modify({
    state,
    setState,
    create,
    apiKeys,
}: ModifyProps) {

    useEffect(() => {
        if (!state.selectedKeyId) {
            setState({ ...state, selectedKeyId: apiKeys[0].id })
        }
    }, [state])

    return (
        <div>
            <DialogContentText>
                Create a new secret value.<br />
                The Encryption / Decryption is not stored by the server, it's only used by the server to encrypt / decrypt the send value<br />
                This also means <b>you should store the encryption key used here</b> somewhere safe
            </DialogContentText>

            <div className="apiKeyAndKey">
                <FormControl fullWidth variant="filled">
                    <InputLabel htmlFor="secret-key-id">Api key</InputLabel>

                    {/* Show placeholder select while we're still loading the keys */}
                    {!state.selectedKeyId
                        ? <Select
                            id="secret-key-id"
                            disabled
                            value=""
                        >
                            <MenuItem value="">
                                <em>None</em>
                            </MenuItem>
                        </Select>
                        : <Select
                            disabled={!create}
                            value={state.selectedKeyId}
                            onChange={(id: any) => setState(s => ({ ...s, selectedKeyId: id.target.value }))}
                            id="secret-key-id"
                        >
                            {apiKeys?.map((key: ApiKey) =>
                                <MenuItem key={key.id} value={key.id}>{
                                    key.name
                                        ? key.name
                                        : <>{key.id} <i>{key.domains.join(', ')}</i></>
                                }</MenuItem>
                            )}
                        </Select>
                    }
                    <FormHelperText>The API key selected will be able to access the secret</FormHelperText>
                </FormControl>

                <div className="divider">/</div>

                <TextField
                    id='key'
                    label='Key'
                    error={!state.key}
                    value={state.key}
                    onChange={(e) => setState(s => ({ ...s, key: e.target.value }))}
                    variant="filled"
                    helperText="the secret identifier used to access the key"
                    fullWidth
                    disabled={!create}
                />
            </div>

            <div className="marginTop">
                <TextField
                    id='description'
                    label='Description'
                    value={state.description}
                    onChange={(e) => setState(s => ({ ...s, description: e.target.value }))}
                    variant="filled"
                    helperText="Additional information that describes the value"
                    multiline
                    fullWidth
                />
            </div>

            <div className="inputWithButton marginTop">
                <TextField
                    id="secret"
                    label="Encryption key"
                    value={state.decryptionKey}
                    error={!!state.decryptionKeyError}
                    helperText={state.decryptionKeyError}
                    onChange={(e) => {
                        const { value } = e.target
                        setState(s => ({
                            ...s,
                            decryptionKey:
                                value,
                            decryptionKeyError: value.length < 16
                                ? 'encryption key value must have a minimal length of 16 chars'
                                : ''
                        }))
                    }}
                    variant="filled"
                    fullWidth
                />
                <div className="toggles">
                    <Tooltip title='Generate random value'>
                        <Button
                            onClick={() => setState(s => ({
                                ...s,
                                decryptionKey: randomString(32),
                                decryptionKeyError: '',
                            }))}
                        >
                            <RefreshIcon fontSize="small" />
                        </Button>
                    </Tooltip>

                </div>
            </div>

            <ModifyValue
                valueStructure={state.valueStructure}
                setValueStructure={ValueStructure => setState(s => ({ ...s, valueStructure: ValueStructure }))}
                value={state.value}
                setValue={(setter) => setState(currentState => {
                    const value = setter(currentState.value)

                    let valueError = ''
                    try {
                        JSON.parse(value)
                    } catch (e: any) {
                        valueError = e.message
                    }

                    return {
                        ...currentState,
                        value,
                        valueError,
                    }
                })}
                valueError={state.valueError}
            />

            <style jsx>{`
                .apiKeyAndKey {
                    display: flex;
                }
                .apiKeyAndKey .divider {
                    margin: 0 10px;
                    font-size: 20px;
                    font-weight: bold;
                    margin-top: 15px;
                }
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
            `}</style>
            <style global jsx>{`
                .secretModalSecretValueInput textarea, .secretModalSecretValueInput input {
                    font-family:Consolas,Monaco,Lucida Console,Liberation Mono,DejaVu Sans Mono,Bitstream Vera Sans Mono,Courier New, monospace;
                }
            `}</style>
        </div>
    )
}
