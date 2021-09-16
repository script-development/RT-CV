import { Secret } from '../../src/types'
import { ModalKind } from '../modal'

export interface SecretModalProps {
    kind: ModalKind
    onClose: () => void
    secret?: Secret
}
