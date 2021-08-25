import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Snackbar } from '@material-ui/core'
import React, { useState, useEffect } from 'react'

export enum ModalKind {
    Closed = 'closed',
    Create = 'create',
    View = 'view',
    Edit = 'edit',
    Delete = 'delete',
}

export type ModalContentOptions = {
    create?: string,
    edit?: string,
    view?: string,
    delete?: string,
} | string

interface ModalProps {
    kind: ModalKind
    onClose: () => void
    onSubmit: () => void
    title: ModalContentOptions

    confirmText?: ModalContentOptions
    cancelText?: ModalContentOptions
    children?: (kind: ModalKind) => React.ReactNode
    submitDisabled?: boolean
    apiError?: string
    setApiError?: (error: string) => void
    fullWidth?: boolean
    showConfirm?: boolean
}

export function Modal({
    kind,
    onClose,
    onSubmit,
    title,
    confirmText = { create: 'Create', edit: 'Save', delete: 'Delete', view: 'View' },
    cancelText = 'Cancel',
    children,
    apiError = '',
    setApiError = () => { },
    submitDisabled = false,
    fullWidth = false,
    showConfirm = true,
}: ModalProps) {
    // Inner kinds reflects the value of kind only if the kind != KeyModalKind.Closed
    // This makes it so when you close the modal the content won't change while the closing animation is playing
    const [innerKind, setInnerKind] = useState(ModalKind.Create)

    useEffect(() => {
        if (kind != ModalKind.Closed)
            setInnerKind(kind)
    }, [kind])

    const resolveModalContentOpts = (value: ModalContentOptions): string =>
        typeof value == 'string' ? value : (value as { [key: string]: string })[innerKind]

    return (
        <Dialog open={kind != ModalKind.Closed} onClose={onClose} fullWidth={fullWidth}>
            <Snackbar
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                open={!!apiError}
                onClose={() => setApiError('')}
                autoHideDuration={6000}
                message={apiError}
                key="key-api-error"
            />

            <DialogTitle>{resolveModalContentOpts(title)}</DialogTitle>


            <DialogContent>
                <form
                    onSubmit={e => {
                        e.preventDefault()
                        onSubmit()
                    }}
                >
                    {children?.(innerKind)}
                </form>
            </DialogContent>

            <DialogActions>
                <Button onClick={onClose}>
                    {resolveModalContentOpts(cancelText) || 'Cancel'}
                </Button>
                {showConfirm ? <Button onClick={onSubmit} color="primary" disabled={submitDisabled}>
                    {resolveModalContentOpts(confirmText)}
                </Button> : ''}
            </DialogActions>
        </Dialog >
    )
}
