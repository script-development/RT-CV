import { Button, Dialog, DialogTitle, DialogContentText, DialogContent, TextField, DialogActions } from '@material-ui/core'
import { useEffect } from 'react'
import { useState } from 'react'
import { ApiKey } from '../src/types'

export enum KeyModalKind {
    Closed,
    Create,
    Edit,
    Delete,
}

interface KeyModalArgs {
    kind: KeyModalKind
    onClose: () => void

    // Only required if kind = Edit or Delete
    apiKey?: ApiKey
}

export function KeyModal({ kind, onClose, apiKey = undefined }: KeyModalArgs) {
    // Inner kinds reflects the value of kind only if the kind != KeyModalKind.Closed
    // This makes it so when you close the modal the content won't change while the closing animation is playing
    const [innerKind, setInnerKind] = useState(KeyModalKind.Create)

    useEffect(() => {
        if (kind != KeyModalKind.Closed) setInnerKind(kind)
    }, [kind])

    const titleText = innerKind == KeyModalKind.Create ? 'Create Api key' : innerKind == KeyModalKind.Edit ? 'Edit Api key' : 'Delete Api key'
    const confirmText = innerKind == KeyModalKind.Create ? 'Create' : innerKind == KeyModalKind.Edit ? 'Save' : 'Delete'

    return (
        <Dialog open={kind != KeyModalKind.Closed} onClose={onClose}>
            <DialogTitle>{titleText}</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    To subscribe to this website, please enter your email address here. We will send updates
                    occasionally.
                </DialogContentText>
                <TextField
                    autoFocus
                    margin="dense"
                    id="name"
                    label="Email Address"
                    type="email"
                    fullWidth
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={onClose} color="primary">{confirmText}</Button>
            </DialogActions>
        </Dialog>
    )
}