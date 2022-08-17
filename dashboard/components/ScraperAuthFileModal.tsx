import { Button, DialogContentText } from "@mui/material"
import dynamic from 'next/dynamic'

import { ApiKey } from "../src/types"
import { Modal, ModalKind } from "./modal"

// Use dynamic so we only load the bytes of this component once we really need it
// It'q quite large and loading it this way splits it from the main bundle
// So that means the first load of the page stays quick
const JSONCode = dynamic(() => import('./jsonCode'), { ssr: false })

interface ScraperAuthFileModalProps {
    apiKey: ApiKey | undefined
    onClose: () => void
}

export function ScraperAuthFileModal({ apiKey, onClose }: ScraperAuthFileModalProps) {
    const fileContents = apiKey ? {
        primary_server: {
            server_location: location.origin,
            api_key_id: apiKey.id,
            api_key: apiKey.key,
        },
    } : {}

    const download = () => {
        const content = encodeURIComponent(JSON.stringify(fileContents))

        const el = document.createElement('a')
        el.style.display = 'none'

        el.setAttribute('href', 'data:application/json;charset=utf-8,' + content)
        el.setAttribute('download', 'env.json')

        document.body.appendChild(el)
        el.click()
        document.body.removeChild(el)
    }

    return (
        <Modal
            kind={apiKey === undefined ? ModalKind.Closed : ModalKind.View}
            onClose={onClose}
            onSubmit={onClose}
            title='Scraper authentication file'
            showConfirm={false}
            cancelText='Close'
            fullWidth
        >{
                _ => <div>
                    <JSONCode json={fileContents} />
                    <DialogContentText>
                        Note that this file has more options defined in <a target="_blank" rel="noopener noreferrer" href="https://github.com/script-development/rtcv_scraper_client#2-obtain-a-envjson">rtcv_scraper_client README.md</a>
                    </DialogContentText>
                    <Button onClick={() => download()} >Download env.json</Button>
                </div>
            }</Modal>
    )
}
