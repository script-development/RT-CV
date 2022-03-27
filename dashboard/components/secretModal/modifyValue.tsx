import { Button, DialogContentText, TextField, Tooltip, FormHelperText, FormControlLabel, Switch } from "@mui/material"
import FormatIndentIncrease from "@mui/icons-material/FormatIndentIncrease"
import Code from '@mui/icons-material/Code'
import People from '@mui/icons-material/People'
import Person from '@mui/icons-material/Person'
import Delete from '@mui/icons-material/Delete'
import PersonAdd from '@mui/icons-material/PersonAdd'
import React, { ReactNode, useEffect, useMemo, useState } from "react"
import { SecretValueStructure } from '../../src/types'

interface ModifyValueProps {
    valueStructure: SecretValueStructure | undefined
    setValueStructure: (kind: SecretValueStructure) => void
    value: string
    valueError: string
    setValue: (setter: (prev: string) => string) => void
}

export default function ModifyValue(props: ModifyValueProps) {
    const optionsToShow: { [key: string]: any } = {
        [SecretValueStructure.Free]: JsonValueKind,
        [SecretValueStructure.StrictUser]: StrictUserValueKind,
        [SecretValueStructure.StrictUsers]: StrictUsersValueKind,
    }
    const Show = optionsToShow[props.valueStructure as string] || ChoseValueStructure

    return (
        <div style={{ margin: '10px 0' }}>
            <Show {...props} />
        </div>
    )
}

interface PrettyValueEditorContainerProps {
    children: ReactNode
    label: string
    value: string
    valueError: string
    setValue: (setter: (prev: string) => string) => void
    onJsonSwitched?: (newJsonValue: boolean) => void
}

function PrettyValueEditorContainer({ children, label, value, setValue, valueError, onJsonSwitched }: PrettyValueEditorContainerProps) {
    const [useJson, setUseJson] = useState(false)

    useEffect(() => {
        onJsonSwitched?.(useJson)
    }, [useJson])

    return (
        <div className="root">
            <div className="labelAndJsonSwitch">
                <DialogContentText>{label}</DialogContentText>

                <FormControlLabel
                    disabled={!!valueError}
                    control={<Switch checked={useJson} onChange={() => setUseJson(v => !v)} />}
                    label="JSON"
                    labelPlacement="start"
                />
            </div>

            {useJson
                ? <JsonValueKind value={value} setValue={setValue} valueError={valueError} />
                : <>
                    {children}

                    <div className="valueToBeStored">
                        <DialogContentText>Data that will be stored:</DialogContentText>
                        <pre className="value">{value}</pre>
                    </div>
                </>
            }
            <style jsx>{`
                .labelAndJsonSwitch {
                    display: flex;
                    align-items: center;
                    justify-content: space-between;
                    padding-bottom: 10px;
                    padding-right: 15px;
                }
                .root {
                    border: 2px solid rgba(255, 255, 255, 0.15);
                    padding: 10px;
                    border-radius: 6px;
                }
                .valueToBeStored {
                    margin-top: 10px;
                }
                .value {
                    background-color: rgba(255, 255, 255, 0.09);
                    color: rgba(255,255,255,0.6);
                    padding: 10px;
                    border-radius: 4px;
                    font-size: 16px;
                    font-family:Consolas,Monaco,Lucida Console,Liberation Mono,DejaVu Sans Mono,Bitstream Vera Sans Mono,Courier New, monospace;
                }
            `}</style>
        </div>
    )
}

function StrictUserValueKind({ value, valueError, setValue }: ModifyValueProps) {
    const [user, setUser] = useState({ username: '', password: '', edited: true })

    const setUserFromValue = (value: string) => {
        try {
            const { username, password } = JSON.parse(value)

            if (typeof username != 'string' || typeof password != 'string')
                throw new Error('Invalid JSON')

            setUser(u => ({ username, password, edited: u.edited }))
        } catch (e) { }
    }

    useEffect(() => {
        if (value && !user.username && !user.password) setUserFromValue(value)
    }, [value])

    useEffect(() => {
        if (user.edited)
            setValue(() => JSON.stringify({ username: user.username, password: user.password }, null, 2))
    }, [user])

    return (
        <PrettyValueEditorContainer
            label="User:"
            value={value}
            setValue={setValue}
            valueError={valueError}
            onJsonSwitched={jsonView => {
                // Update the user object once we switch back from the raw json input
                if (!jsonView) setUserFromValue(value)
            }}
        >
            <div className="inputs">
                <div>
                    <TextField
                        value={user.username}
                        onChange={e =>
                            setUser(s => ({ ...s, username: e.target.value, edited: true }))
                        }
                        id="username"
                        label="Username / Email"
                        variant="filled"
                        fullWidth
                    />
                </div>
                <div>
                    <TextField
                        value={user.password}
                        onChange={e =>
                            setUser(s => ({ ...s, password: e.target.value, edited: true }))
                        }
                        id="password"
                        label="Password"
                        variant="filled"
                        fullWidth
                    />
                </div>
            </div>
            <style jsx>{`
                .inputs {
                    display: flex;
                }
                .inputs > div {
                    flex-grow: 1;
                }
                .inputs > div:first-child {
                    margin-right: 10px;
                }
            `}</style>
        </PrettyValueEditorContainer>
    )
}

function StrictUsersValueKind({ value, valueError, setValue }: ModifyValueProps) {
    const [users, setUsers] = useState([{ username: '', password: '' }])
    const [modified, setModified] = useState(false)

    useEffect(() => {
        if (modified)
            setValue(() => JSON.stringify(users, null, 2))
    }, [users, modified])

    const setUsersFromValue = (value: string) => {
        try {
            const usersFromValue = JSON.parse(value).map((u: any) => {
                const { username, password } = u

                if (typeof username != 'string' || typeof password != 'string')
                    throw 'Invalid JSON'

                return { username, password }
            })
            setUsers(usersFromValue)
        } catch (e) {
            setValue(() => JSON.stringify(users, null, 2))
        }
    }

    useEffect(() => {
        setUsersFromValue(value)
    }, [])

    return (
        <PrettyValueEditorContainer
            label="Users:"
            value={value}
            setValue={setValue}
            valueError={valueError}
            onJsonSwitched={jsonView => {
                // Update the users array once we switch back from the raw json input
                if (!jsonView) setUsersFromValue(value)
            }}
        >
            <div className="inputs">
                {users.map((user, idx) =>
                    <div className="row" key={idx}>
                        <div className="input">
                            <TextField
                                value={user.username}
                                onChange={e => {
                                    setUsers(users => {
                                        users[idx].username = e.target.value
                                        return [...users]
                                    })
                                    setModified(true)
                                }}
                                id="username"
                                label="Username / Email"
                                variant="filled"
                                fullWidth
                            />
                        </div>
                        <div className="input">
                            <TextField
                                value={user.password}
                                onChange={e => {
                                    setUsers(users => {
                                        users[idx].password = e.target.value
                                        return [...users]
                                    })
                                    setModified(true)
                                }}
                                id="password"
                                label="Password"
                                variant="filled"
                                fullWidth
                            />
                        </div>
                        <div className="removeRow">
                            <Button
                                onClick={() => {
                                    setUsers(users => [...users.slice(0, idx), ...users.slice(idx + 1)])
                                    setModified(true)
                                }}
                                variant="outlined"
                            ><Delete fontSize="small" /></Button>
                        </div>
                    </div>
                )}
                <div className="addRow">
                    <Button
                        onClick={() => {
                            setUsers(users => [...users, { username: '', password: '' }])
                            setModified(true)
                        }}
                        variant="outlined"
                        fullWidth
                    ><PersonAdd fontSize="small" /></Button>
                </div>
            </div>
            <style jsx>{`
                .inputs .row {
                    display: flex;
                    align-items: center;
                    margin-bottom: 10px;
                }
                .inputs .row > .input {
                    flex-grow: 1;
                }
                .inputs .row > .input {
                    margin-right: 10px;
                }
                .addRow {
                    margin-top: 5px;
                    display: flex;
                    justify-content: flex-end;
                }
            `}</style>
        </PrettyValueEditorContainer>
    )
}

interface JsonValueKindProps {
    value: string
    valueError: string
    setValue: (setter: (prev: string) => string) => void
}

function JsonValueKind({ value, setValue, valueError }: JsonValueKindProps) {
    return (
        <div className="root">
            <TextField
                className="secretModalSecretValueInput"
                id="secret-value"
                label="JSON Value"
                value={value}
                helperText={valueError || 'json value is valid'}
                error={!!valueError}
                onChange={(e) => setValue(() => e.target.value)}
                variant="filled"
                multiline
                fullWidth
            />
            <div className="toggles">
                <Tooltip title='Format json'>
                    <Button
                        disabled={!!valueError}
                        onClick={() => {
                            setValue(prevValue => {
                                try {
                                    return JSON.stringify(JSON.parse(prevValue), null, 2)
                                } catch (e: any) { }
                                return prevValue
                            })
                        }}
                    >
                        <FormatIndentIncrease fontSize="small" />
                    </Button>
                </Tooltip>
            </div>
            <style jsx>{`
                .root {
                    display: flex;
                    justify-content: space-between;
                    align-items: flex-start;
                }
                .toggles {
                    margin-left: 10px;
                }
            `}</style>
        </div>
    )
}

function ChoseValueStructure({ setValueStructure }: ModifyValueProps) {
    const sides = [
        {
            title: 'Strict',
            info: 'Chose from a set of pre defined layouts so the programs using RT-CV have predictable json responses.',
            actions: [
                {
                    label: 'User',
                    tooltip: 'Set the value type to be a single user',
                    icon: <Person />,
                    structure: SecretValueStructure.StrictUser,
                },
                {
                    label: 'Users',
                    tooltip: 'Set the value type to be a list of users',
                    icon: <People />,
                    structure: SecretValueStructure.StrictUsers,
                },
            ],
        },
        {
            title: 'Free input',
            info: `Free json inputs without any requirements, as long as it's valid json it's ok. This will mean tough that programs using RT-CV also needs to understand the value you put in.`,
            actions: [
                {
                    label: 'Use free input',
                    tooltip: 'Set the value type to be free json input',
                    icon: <Code />,
                    structure: SecretValueStructure.Free,
                },
            ],
        },
    ]

    return (
        <div className="root">
            <DialogContentText>Chose value type</DialogContentText>
            <FormHelperText error>The secret needs a value type to be created</FormHelperText>
            <div className="options">
                {sides.map((side, idx) =>
                    <div key={idx}>
                        <h3>{side.title}</h3>
                        <div className="info">
                            <DialogContentText>{side.info}</DialogContentText>
                        </div>
                        <div className="actions">
                            {side.actions.map((action, actionIdx) =>
                                <span key={actionIdx}>
                                    <Tooltip title={action.tooltip}>
                                        <Button
                                            variant="outlined"
                                            startIcon={action.icon}
                                            onClick={() => setValueStructure(action.structure)}
                                        >{action.label}</Button>
                                    </Tooltip>
                                </span>
                            )}
                        </div>
                    </div>
                )}
            </div>
            <style jsx>{`
                .options {
                    display: flex;
                }
                .options > div {
                    border: 2px solid rgba(255, 255, 255, 0.15);
                    margin: 5px;
                    padding: 10px;
                    border-radius: 6px;
                    display: flex;
                    flex-direction: column;
                }
                .info {
                    flex-grow: 1;
                }
                .actions {
                    display: flex;
                }
                .actions > span {
                    margin: 2px;
                }
            `}</style>
        </div>
    )
}
