import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Snackbar } from '@material-ui/core'
import React, { useState, useEffect } from 'react'

export enum ModalKind {
    Closed = 'closed',
    Create = 'create',
    Edit = 'edit',
    Delete = 'delete',
}

export type ModalContentOptions = {
    create: string,
    edit: string,
    delete: string,
} | string

interface ModalProps {
    kind: ModalKind
    onClose: () => void
    onSubmit: () => void
    title: ModalContentOptions,

    confirmText?: ModalContentOptions,
    children?: (kind: ModalKind) => React.ReactNode
    submitDisabled?: boolean,
    apiError?: string
    setApiError?: (error: string) => void

}

export function Modal({
    kind,
    onClose,
    onSubmit,
    title,
    confirmText = { create: 'Create', edit: 'Save', delete: 'Delete' },
    children,
    apiError = '',
    setApiError = () => { },
    submitDisabled = false,
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
        <Dialog open={kind != ModalKind.Closed} onClose={onClose}>
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
                {children?.(innerKind)}
            </DialogContent>

            <DialogActions>
                <Button onClick={onClose}>
                    Cancel
                </Button>
                <Button onClick={onSubmit} color="primary" disabled={submitDisabled}>
                    {resolveModalContentOpts(confirmText)}
                </Button>
            </DialogActions>
        </Dialog>
    )
}
