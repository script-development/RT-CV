import { OnMatchHook } from '../../src/types'
import { ModalKind } from '../modal'

export interface ModalProps {
    kind: ModalKind
    onClose: () => void
    hook?: OnMatchHook
}
