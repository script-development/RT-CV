import {
    Button,
    DialogContentText,
    TextField,
    Checkbox,
    FormControlLabel,
    Tooltip,
    FormControl,
    FormHelperText,
    FormLabel,
} from '@material-ui/core'
import React, { useState, useEffect } from 'react'
import { ApiKey } from '../src/types'
import { allRoles, Roles, roleInfo } from '../src/roles'
import { randomString } from '../src/random'
import RefreshIcon from '@material-ui/icons/Refresh'
import { fetcher } from '../src/auth'
import { Modal, ModalKind } from './modal'

export enum KeyModalKind {
    Closed,
    Create,
    Edit,
    Delete,
}

interface KeyModalProps {
    kind: ModalKind
    onClose: () => void

    // Only required if kind = Edit or Delete
    apiKey?: ApiKey
}

export function KeyModal({ kind, onClose, apiKey = undefined }: KeyModalProps) {
    const [state, setState] = useState<ApiKey>({
        domains: ['*'],
        enabled: true,
        id: '',
        key: randomString(32),
        roles: 0,
        system: false,
    })
    const [apiError, setApiError] = useState('')

    const formControlStyle = { marginTop: 10 }
    const rolesError = state.roles == 0 ? 'You need at least one role' : undefined
    const disabled = state.system
    const canSubmit = !rolesError && !disabled

    const submit = async () => {
        try {
            if (kind == ModalKind.Create)
                await fetcher.post(`/api/v1/keys`, state)
            else if (kind == ModalKind.Edit)
                await fetcher.put(`/api/v1/keys/${state.id}`, state)
            else
                await fetcher.delete(`/api/v1/keys/${state.id}`)

            onClose()
        } catch (e) {
            setApiError(e?.message || e)
        }
    }
    const refreshKey = () => setState(v => ({ ...v, key: randomString(32) }))

    useEffect(() => {
        if (apiKey != undefined && apiKey.id != state.id)
            setState(apiKey)
        else if (apiKey == undefined && state.id)
            setState({
                domains: ['*'],
                enabled: true,
                id: '',
                key: randomString(32),
                roles: 0,
                system: false,
            })
    }, [kind, apiKey])

    return (
        <Modal
            kind={kind}
            onClose={onClose}
            onSubmit={submit}
            title={{
                create: 'Create Api key',
                edit: 'Edit Api key',
                delete: 'Delete Api key',
            }}
            submitDisabled={!canSubmit}
            apiError={apiError}
            setApiError={setApiError}
        >{(kind: ModalKind) => {
            if (kind == ModalKind.Delete)
                return (<DialogContentText>
                    {state.system
                        ? 'System keys cannot be deleted via this UI'
                        : 'Are you sure you want to delete this api key?'
                    }
                </DialogContentText>)
            else
                return (
                    <div>
                        <DialogContentText>
                            {
                                kind == ModalKind.Create
                                    ? 'create a new api key to authenticate with RT-CV'
                                    : 'Edit this api key'
                            }
                        </DialogContentText>

                        <TextField
                            id="domains"
                            label="Domains"
                            multiline
                            value={state.domains.join('\n')}
                            onChange={(e) => setState(v => ({ ...v, domains: e.target.value.split('\n') }))}
                            variant="filled"
                            disabled={disabled}
                            helperText="every new line is a new domain, use * to wildcard"
                        />

                        <div className="checkboxWithFormControl">
                            <Checkbox
                                disabled={disabled}
                                checked={state.enabled}
                                onChange={() => setState((v) => ({ ...v, enabled: !v.enabled }))}
                                color="primary"
                            />
                            <FormControl disabled={disabled} fullWidth>
                                <FormLabel>Enabled</FormLabel>
                                <FormHelperText>Determines if this key can be used to authenticate</FormHelperText>
                            </FormControl>
                        </div>

                        <RolesSelector
                            error={rolesError}
                            value={state.roles}
                            setValue={newValue => setState((v) => ({ ...v, roles: newValue(v.roles) }))}
                            disabled={disabled}
                        />

                        <div className="checkboxWithFormControl">
                            <Checkbox
                                disabled
                                checked={state.system}
                                color="primary"
                            />
                            <FormControl disabled style={formControlStyle} fullWidth>
                                <FormLabel>System key</FormLabel>
                                <FormHelperText>This field can only be set by the system</FormHelperText>
                            </FormControl>
                        </div>

                        <FormControl disabled={disabled} style={formControlStyle} fullWidth>
                            <FormLabel>Api Key</FormLabel>
                            <div className="apiKeyForm">
                                <div className="apiKeyControls">
                                    <Tooltip title='Refresh key'>
                                        <Button onClick={refreshKey} disabled={disabled}>
                                            <RefreshIcon fontSize="small" />
                                        </Button>
                                    </Tooltip>
                                </div>
                                <div className="apiKey" style={{ color: disabled ? 'gray' : 'white' }}>{state.key}</div>
                            </div>
                        </FormControl>
                        <style jsx>{`
                        .checkboxWithFormControl {
                            display: flex;
                            align-items: center;
                            margin: 10px 0;
                        }
                        .apiKeyForm {
                            display: flex;
                            align-items: center;
                            margin: 10px 0;
                        }
                        .apiKeyForm .apiKeyControls {
                            width: 70px;
                        }
                        .apiKeyForm .apiKey {
                            padding: 2px 10px 5px 10px;
                            border-radius: 4px;
                            background-color: #455a64;
                            font-family: monospace;
                            display: block;
                            flex-grow: 1;
                            word-break: break-all;
                        }
                    `}</style>
                    </div>
                )
        }}</Modal>
    )
}

interface RolesSelectorArgs {
    value: number
    setValue: (newValue: (prev: number) => number) => void
    disabled?: boolean
    error?: string
}

function RolesSelector({ value, setValue, disabled = false, error = undefined }: RolesSelectorArgs) {
    const toggleRole = (role: Roles) =>
        setValue((prev) => (prev ^ role))

    return (
        <div className="root">
            <FormControl required error={!!error} component="fieldset" disabled={disabled} fullWidth>
                <FormLabel>Roles</FormLabel>
                <FormHelperText>
                    Roles number: <span className="roleNr">{value}</span>
                </FormHelperText>
                <div className="checkboxes">
                    {allRoles.map(role => {
                        const info = roleInfo(role)

                        return (<div key={role}>
                            <Tooltip title={info?.description || 'unknown'}>
                                <FormControlLabel
                                    disabled={disabled}
                                    control={
                                        <Checkbox
                                            disabled={disabled}
                                            checked={(value & role) == role}
                                            onChange={() => toggleRole(role)}
                                            color="primary"
                                        />
                                    }
                                    label={info?.title || 'unknown'}
                                />
                            </Tooltip>
                        </div>)
                    })}
                </div>
                <FormHelperText>{error}</FormHelperText>
            </FormControl>

            <style jsx>{`
                .roleNr {
                    font-weight: bold;
                }
                .checkboxes {
                    display: flex;
                    flex-wrap: wrap;
                }
            `}</style>
        </div>
    )
}
