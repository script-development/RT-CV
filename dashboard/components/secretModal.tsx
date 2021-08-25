import {
    Button,
    Dialog,
    DialogTitle,
    DialogContentText,
    DialogContent,
    TextField,
    DialogActions,
    Checkbox,
    FormControlLabel,
    Tooltip,
    FormControl,
    FormHelperText,
    FormLabel,
    Snackbar,
} from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { Modal, ModalKind } from './modal'

interface SecretModalProps {
    kind: ModalKind
    onClose: () => void
}

export function SecretModal({ kind, onClose }: SecretModalProps) {
    const [apiError, setApiError] = useState('')

    const canSubmit = true

    const submit = async () => {
        try {
            // TODO
            onClose()
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
                edit: 'Edit Secret',
                delete: 'Delete Secret',
            }}
            submitDisabled={!canSubmit}
            apiError={apiError}
            setApiError={setApiError}
        >{(kind: ModalKind) => {
            if (kind == ModalKind.Delete)
                return (<DialogContentText>
                    Are you sure you want to delete this secret?'
                </DialogContentText>)
            else
                return (<div>
                    {/* TODO */}
                </div>)
        }}</Modal>
    )
}
