import { Button, DialogContentText, TextField, Tooltip, FormHelperText } from "@material-ui/core"
import FormatIndentIncrease from "@material-ui/icons/FormatIndentIncrease"
import Code from '@material-ui/icons/Code'
import People from '@material-ui/icons/People'
import Person from '@material-ui/icons/Person'
import Delete from '@material-ui/icons/Delete'
import PersonAdd from '@material-ui/icons/PersonAdd'
import React, { useEffect, useMemo, useState } from "react"
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

function StrictUserValueKind({ value, setValue }: ModifyValueProps) {
    const [user, setUser] = useState({ username: '', password: '', edited: true })

    useEffect(() => {
        if (value && !user.username && !user.password) {
            try {
                const { username, password } = JSON.parse(value)

                if (typeof username != 'string' || typeof password != 'string')
                    throw new Error('Invalid JSON')

                setUser(u => ({ username, password, edited: u.edited }))
            } catch (e) { }
        }
    }, [value])

    useEffect(() => {
        if (user.edited)
            setValue(() => JSON.stringify({ username: user.username, password: user.password }))
    }, [user])

    return (
        <div className="root">
            <DialogContentText>User:</DialogContentText>
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
            <ValueToBeStored value={value} />
            <style jsx>{`
                .root {
                    border: 2px solid rgba(255, 255, 255, 0.15);
                    padding: 10px;
                    border-radius: 6px;
                }
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
        </div>
    )
}

function StrictUsersValueKind({ value, setValue }: ModifyValueProps) {
    const [users, setUsers] = useState([{ username: '', password: '' }])
    const [modified, setModified] = useState(false)

    useEffect(() => {
        if (modified)
            setValue(() => JSON.stringify(users))
    }, [users, modified])

    useEffect(() => {
        try {
            const usersFromValue = JSON.parse(value).map((u: any) => {
                const { username, password } = u

                if (typeof username != 'string' || typeof password != 'string')
                    throw 'Invalid JSON'

                return { username, password }
            })
            setUsers(usersFromValue)
        } catch (e) {
            setValue(() => JSON.stringify(users))
        }
    }, [])

    return (
        <div className="root">
            <DialogContentText>Users:</DialogContentText>
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
                                onClick={() => setUsers(users => [...users.slice(0, idx), ...users.slice(idx + 1)])}
                                variant="outlined"
                            ><Delete fontSize="small" /></Button>
                        </div>
                    </div>
                )}
                <div className="addRow">
                    <Button
                        onClick={() => setUsers(users => [...users, { username: '', password: '' }])}
                        variant="outlined"
                        fullWidth
                    ><PersonAdd fontSize="small" /></Button>
                </div>
            </div>
            <ValueToBeStored value={value} />
            <style jsx>{`
                .root {
                    border: 2px solid rgba(255, 255, 255, 0.15);
                    padding: 10px;
                    border-radius: 6px;
                }
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
        </div>
    )
}

function ValueToBeStored({ value }: { value: string }) {
    const formattedValue = useMemo(() => JSON.stringify(JSON.parse(value || 'null'), null, 2), [value])

    return (
        <div className="root">
            <DialogContentText>Data that will be stored:</DialogContentText>
            <pre className="value">{formattedValue}</pre>
            <style jsx>{`
                .root {
                    margin-top: 10px;
                }
                .root .value {
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

function JsonValueKind({ value, setValue, valueError }: ModifyValueProps) {
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
