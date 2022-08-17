import { Button, DialogContentText, TextField, Typography } from "@mui/material"
import { useEffect, useState } from "react"
import { fetcher } from '../src/auth'
import { ApiKey } from "../src/types"
import { Modal, ModalKind } from "./modal"
import { GridColDef, DataGrid } from '@mui/x-data-grid'
import { Delete, Add } from '@mui/icons-material'

interface LoginUser {
    username: string
    password?: string
}

interface ScraperUsersModalProps {
    apiKey: ApiKey | undefined
    onClose: () => void
}

export function ScraperUsersModal({ apiKey, onClose }: ScraperUsersModalProps) {
    const [loginUsers, setLoginUsers] = useState<Array<LoginUser>>()

    const obtainUsersForKey = async () => {
        const data = await fetcher.get('/api/v1/scraperUsers/' + apiKey?.id)
        setLoginUsers(data.users)
    }

    useEffect(() => {
        if (apiKey) {
            obtainUsersForKey()
        } else {
            setLoginUsers(undefined)
        }
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
                _ => <div>
                    {loginUsers
                        ? <ListLoginUsers
                            users={loginUsers}
                            setUsers={setLoginUsers}
                            apiKeyId={apiKey?.id ?? ''} />
                        : <DialogContentText>Loading...</DialogContentText>}
                </div>
            }</Modal>
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
        { field: 'username', headerName: 'Username', width: 250 },
        { field: 'password', headerName: 'Password', valueGetter: (params) => params.row.password || '***' },
        { field: 'actions', headerName: 'Actions', renderCell: (params) => <Button onClick={() => deleteUser(params.row.username)}><Delete /></Button> },
    ]

    return (
        <div>
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
            <DialogContentText>Add a new user</DialogContentText>
            <NewUser
                apiKeyId={apiKeyId}
                newUsers={setUsers}
            />
        </div>
    )
}

interface NewUserParams {
    apiKeyId: string
    newUsers: (users: Array<LoginUser>) => void
}

function NewUser({ apiKeyId, newUsers }: NewUserParams) {
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
            <div className="input">
                <TextField
                    id="new-user-username"
                    label="Username"
                    value={username}
                    fullWidth
                    onChange={e => setNewLoginUser(curr => ({ ...curr, username: e.target.value }))}
                    variant="filled"
                    disabled={formStatus.uploading}
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
                    disabled={formStatus.uploading}
                    error={!!formStatus.error}
                />
                <Button
                    disabled={formStatus.uploading}
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
