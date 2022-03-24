import { Secret } from '../../src/types'
import { ModalKind } from '../modal'

export interface SecretModalProps {
    kind: ModalKind
    setKind: (newKind: ModalKind) => void
    onClose: () => void
    secret?: Secret
}
