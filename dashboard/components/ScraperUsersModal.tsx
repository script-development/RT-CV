import { Button, DialogContentText, TextField, Typography } from "@mui/material"
import { useEffect, useState } from "react"
import { fetcher } from '../src/auth'
import { ApiKey } from "../src/types"
import { Modal, ModalKind } from "./modal"
import { GridColDef, DataGrid } from '@mui/x-data-grid'
import { Delete, Add, ChangeCircle, FmdBadTwoTone } from '@mui/icons-material'

interface LoginUser {
    username: string
    password?: string
}

interface ScraperUsersModalProps {
    apiKey: ApiKey | undefined
    onClose: () => void
}

export function ScraperUsersModal({ apiKey, onClose }: ScraperUsersModalProps) {
    const [state, setState] = useState<{
        users: Array<LoginUser>
        publicKey: string
    }>()

    const obtainUsersForKey = async () => {
        const data = await fetcher.get('/api/v1/scraperUsers/' + apiKey?.id)
        setState({
            users: data.users,
            publicKey: data.scraperPubKey
        })
    }

    const setLoginUsers = (users: Array<LoginUser>) => setState(s => s ? { ...s, users } : undefined)
    const setPublicKey = (publicKey: string) => setState(s => s ? { ...s, publicKey } : undefined)

    useEffect(() => {
        if (apiKey) obtainUsersForKey()
        else setState(undefined)
    }, [apiKey])

    return (
        <Modal
            kind={apiKey === undefined ? ModalKind.Closed : ModalKind.View}
            onClose={onClose}
            onSubmit={onClose}
            title='Scraper login users'
            showConfirm={false}
            cancelText='Close'
            fullWidth
        >{
                _ => state ? <div>
                    <ManagePublicKey
                        apiKeyId={apiKey?.id ?? ''}
                        currentPublicKey={state.publicKey}
                        usersCount={state.users.length}
                        updatePublicKey={setPublicKey}
                    />
                    <NewUser
                        hasPublicKey={Boolean(state.publicKey)}
                        apiKeyId={apiKey?.id ?? ''}
                        newUsers={setLoginUsers}
                    />
                    <ListLoginUsers
                        users={state.users}
                        setUsers={setLoginUsers}
                        apiKeyId={apiKey?.id ?? ''}
                    />
                </div> : <DialogContentText>Loading...</DialogContentText>
            }</Modal>
    )
}

interface ManagePublicKeyProps {
    apiKeyId: string
    currentPublicKey: string
    usersCount: number
    updatePublicKey: (value: string) => void
}

function ManagePublicKey({ currentPublicKey, usersCount, apiKeyId, updatePublicKey }: ManagePublicKeyProps) {
    const [publicKey, setPublicKey] = useState(currentPublicKey)
    const [formStatus, setFormStatus] = useState({ uploading: false, error: undefined as undefined | string })

    const changePublicKey = async () => {
        setFormStatus({ uploading: true, error: undefined })
        try {
            const resp = await fetcher.fetch(`/api/v1/scraperUsers/${apiKeyId}/setPublicKey`, 'PATCH', { publicKey })
            updatePublicKey(resp.scraperPubKey)
            setFormStatus({ uploading: false, error: undefined })
        } catch (e) {
            setFormStatus({ uploading: false, error: `${e}` })
        }
    }

    const validBase64 = (testStr: string) => {
        try {
            atob(testStr)
            return undefined
        } catch (e) {
            return 'Invalid base64'
        }
    }

    const validation = publicKey.length != 44 ? 'Invalid length public key' : validBase64(publicKey)

    return (
        <div>
            <Typography variant='h6'>Public key</Typography>
            {currentPublicKey && usersCount != 0
                ? <Typography color='orange'>Warning: changing the public key will remove all current users</Typography>
                : undefined
            }
            <div className="input">
                <TextField
                    id="public-key"
                    label="Public key"
                    value={publicKey}
                    fullWidth
                    onChange={e => setPublicKey(e.target.value)}
                    variant="filled"
                    disabled={formStatus.uploading}
                    error={!!formStatus.error || !!validation}
                    helperText={validation}
                />
                <Button
                    disabled={formStatus.uploading}
                    onClick={changePublicKey}
                ><ChangeCircle /> Change</Button>
            </div>
            {formStatus.error}
            <style jsx>{`
                .input {
                    display: grid;
                    grid-template-columns: 1fr 100px;
                    grid-column-gap: 10px;
                }
            `}</style>
        </div>
    )
}

interface ListLoginUsersProps {
    users: Array<LoginUser>
    apiKeyId: string
    setUsers: (users: Array<LoginUser>) => void
}

function ListLoginUsers({ users, apiKeyId, setUsers }: ListLoginUsersProps) {
    const deleteUser = async (username: string) => {
        const resp = await fetcher.fetch('/api/v1/scraperUsers/' + apiKeyId, 'DELETE', { username })
        setUsers(resp.users)
    }

    const columns: Array<GridColDef<LoginUser>> = [
        { field: 'username', headerName: 'Username', width: 400 },
        { field: 'actions', headerName: 'Actions', renderCell: (params) => <Button onClick={() => deleteUser(params.row.username)}><Delete /></Button> },
    ]

    return (
        <div>
            <Typography variant='h6'>Login users</Typography>
            <DataGrid
                getRowId={(row) => row.username}
                rows={users}
                columns={columns}
                autoHeight
                disableColumnFilter
                disableColumnMenu
                disableDensitySelector
                disableColumnSelector
                disableSelectionOnClick
                disableVirtualization
            />
        </div>
    )
}

interface NewUserParams {
    apiKeyId: string
    hasPublicKey: boolean
    newUsers: (users: Array<LoginUser>) => void
}

function NewUser({ apiKeyId, newUsers, hasPublicKey }: NewUserParams) {
    const [{ username, password }, setNewLoginUser] = useState<LoginUser>({ username: '', password: '' })
    const [formStatus, setFormStatus] = useState({ uploading: false, error: undefined as undefined | string })

    const addNewUser = async () => {
        setFormStatus({ uploading: true, error: undefined })
        try {
            const resp = await fetcher.fetch('/api/v1/scraperUsers/' + apiKeyId, 'PATCH', { username, password })
            newUsers(resp.users)
            setNewLoginUser({ username: '', password: '' })
            setFormStatus({ uploading: false, error: undefined })
        } catch (e) {
            setFormStatus({ uploading: false, error: `${e}` })
        }
    }

    return (
        <div>
            <Typography variant='h6'>Add a new user</Typography>
            {!hasPublicKey
                ? <Typography color='red'>No public key set!, You cannot add a new user without a public key set. A scraper users key pair can be generated using <a target="_blank" rel="noopener noreferrer" href="https://github.com/script-development/rtcv_scraper_client/tree/v2/gen_key">the gen key program from the rtcv scraper client</a></Typography>
                : undefined
            }
            <div className="input">
                <TextField
                    id="new-user-username"
                    label="Username"
                    value={username}
                    fullWidth
                    onChange={e => setNewLoginUser(curr => ({ ...curr, username: e.target.value }))}
                    variant="filled"
                    disabled={formStatus.uploading || !hasPublicKey}
                    error={!!formStatus.error}
                />
                <TextField
                    id="new-user-password"
                    label="Password"
                    value={password}
                    fullWidth
                    onChange={e => setNewLoginUser(curr => ({ ...curr, password: e.target.value }))}
                    variant="filled"
                    type="password"
                    disabled={formStatus.uploading || !hasPublicKey}
                    error={!!formStatus.error}
                />
                <Button
                    disabled={formStatus.uploading || !hasPublicKey}
                    onClick={addNewUser}
                ><Add /> Add</Button>
            </div>
            {formStatus.error && <Typography>{formStatus.error}</Typography>}
            <style jsx>{`
                .input {
                    display: grid;
                    grid-template-columns: 1fr 1fr 100px;
                    grid-column-gap: 10px;
                }
            `}</style>
        </div>
    )
}
