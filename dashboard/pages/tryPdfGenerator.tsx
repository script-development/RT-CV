import Head from "next/head"
import dynamic from 'next/dynamic'
import Header from '../components/header'
import { useEffect, useRef, useState } from 'react'
import { fetcher } from '../src/auth'
import { Checkbox, FormControl, FormLabel, Select, MenuItem, InputLabel, CircularProgress } from "@material-ui/core"

const JSONCode = dynamic(() => import('../components/jsonCode'), { ssr: false })

interface PdfCreationOptions {
    fontHeader?: string,
    fontRegular?: string,
    style?: string,
    headerColor?: string,
    subHeaderColor?: string,
    logoImageUrl?: string,
    companyName?: string,
    companyAddress?: string,
}

export default function TryPdfGenerator() {
    const previewIframe = useRef<HTMLIFrameElement>(null)
    const loadPreviewTimeout = useRef<NodeJS.Timeout | null>()
    const [options, setOptions] = useState<PdfCreationOptions>({})
    const [loading, setLoading] = useState(true)

    const reLoadPreview = async (options: PdfCreationOptions) => {
        loadPreviewTimeout.current = null;
        const r = await fetcher.fetchNoJsonMarshal(
            '/api/v1/exampleAttachmentPdf',
            'POST',
            { options },
        )
        if (r.status != 200) { return }
        const blob = await r.blob()
        const url = URL.createObjectURL(blob)
        if (previewIframe.current) {
            previewIframe.current.src = url
        }
        console.log(loadPreviewTimeout.current)
        if (!loadPreviewTimeout.current) {
            setLoading(false)
        }
    }

    const fontOptions = [
        { value: 'BeVietnamPro', label: 'Be Vietnam Pro' },
        { value: 'IBMPlexMono', label: 'IBM Plex Mono' },
        { value: 'IBMPlexSans', label: 'IBM Plex Sans' },
        { value: 'IBMPlexSerif', label: 'IBM Plex Serif' },
        { value: 'Lobster', label: 'Lobster' },
        { value: 'OpenSans', label: 'Open Sans' },
        { value: 'PlayfairDisplay', label: 'Playfair Display' },
        { value: 'RobotoSlab', label: 'Roboto Slab' },
    ]

    const styleOptions = [
        { value: 'style_1', label: 'Style 1' },
        { value: 'style_2', label: 'Style 2' },
        { value: 'style_3', label: 'Style 3' },
    ]

    useEffect(() => {
        setLoading(true)
        if (loadPreviewTimeout.current) { clearTimeout(loadPreviewTimeout.current) }
        loadPreviewTimeout.current = setTimeout(() => reLoadPreview(options), 2000)
    }, [options])

    return (
        <div className="container">
            <Head><title>RT-CV home</title></Head>
            <Header />

            <div className="editor">
                <div className="inputs">
                    <div>
                        <h2>Options</h2>
                        <FormControl fullWidth>
                            <FormLabel>Header background color</FormLabel>
                            <div className="addPadding">
                                <Checkbox
                                    checked={options.headerColor !== undefined}
                                    onChange={() =>
                                        setOptions(v => {
                                            if (v.headerColor === undefined) {
                                                v.headerColor = '#ff0000'
                                            } else {
                                                v.headerColor = undefined
                                            }
                                            return { ...v }
                                        })
                                    }
                                    color="primary"
                                />
                                {options.headerColor !== undefined ?
                                    <input
                                        value={options.headerColor}
                                        type="color"
                                        onChange={e => setOptions(v => ({ ...v, headerColor: e.target.value }))}
                                    />
                                    : undefined}
                            </div>
                        </FormControl>
                        <FormControl fullWidth>
                            <FormLabel>Sub header background color</FormLabel>
                            <div className="addPadding">
                                <Checkbox
                                    checked={options.subHeaderColor !== undefined}
                                    onChange={() => setOptions(v => ({
                                        ...v,
                                        subHeaderColor: v.subHeaderColor === undefined
                                            ? '#ff0000'
                                            : undefined,
                                    }))}
                                    color="primary"
                                />
                                {options.subHeaderColor !== undefined ?
                                    <input
                                        value={options.subHeaderColor}
                                        type="color"
                                        onChange={e => setOptions(v => ({ ...v, subHeaderColor: e.target.value }))}
                                    />
                                    : undefined}
                            </div>
                        </FormControl>
                        <FormLabel>Fonts</FormLabel>
                        <div>
                            <Checkbox
                                checked={options.fontRegular !== undefined}
                                onChange={() => setOptions(v => ({
                                    ...v,
                                    fontRegular: v.fontRegular === undefined
                                        ? 'OpenSans'
                                        : undefined,
                                    fontHeader: v.fontHeader === undefined
                                        ? 'OpenSans'
                                        : undefined,
                                }))}
                                color="primary"
                            />
                            {options.fontRegular ? <>
                                <FormControl>
                                    <InputLabel id="font-header-label">Header</InputLabel>
                                    <Select labelId="font-header-label" value={options.fontHeader} onChange={e => setOptions(v => ({ ...v, fontHeader: e.target.value as string | undefined }))}>
                                        {fontOptions.map((font, i) => <MenuItem key={i} value={font.value}>{font.label}</MenuItem>)}
                                    </Select>
                                </FormControl>
                                <div style={{ paddingRight: 10, display: 'inline-block' }} />
                                <FormControl>
                                    <InputLabel id="font-regular-label">Other</InputLabel>
                                    <Select labelId="font-regular-label" value={options.fontRegular} onChange={e => setOptions(v => ({ ...v, fontRegular: e.target.value as string | undefined }))}>
                                        {fontOptions.map((font, i) => <MenuItem key={i} value={font.value}>{font.label}</MenuItem>)}
                                    </Select>
                                </FormControl>
                            </> : undefined}
                        </div>
                        <FormLabel>Layout</FormLabel>
                        <div>
                            <Checkbox
                                checked={options.style !== undefined}
                                onChange={() => setOptions(v => ({
                                    ...v,
                                    style: v.style === undefined
                                        ? 'style_1'
                                        : undefined,
                                }))}
                                color="primary"
                            />
                            {options.style ?
                                <FormControl>
                                    <InputLabel id="style-label">Style</InputLabel>
                                    <Select labelId="style-label" value={options.style} onChange={e => setOptions(v => ({ ...v, style: e.target.value as string | undefined }))}>
                                        {styleOptions.map((entry, i) => <MenuItem key={i} value={entry.value}>{entry.label}</MenuItem>)}
                                    </Select>
                                </FormControl>
                                : undefined}
                        </div>
                    </div>
                    <div>
                        <h2>Config send to server</h2>
                        <JSONCode json={options} />
                    </div>
                </div>
                <div className="preview">
                    <div className="loading-indicator-container">
                        <div
                            className="loading-indicator"
                            style={{
                                opacity: loading ? 1 : 0,
                                transform: `translateY(${loading ? 20 : 0}px)`,
                            }}
                        >
                            <CircularProgress size={20} />
                            <span>Loading</span>
                        </div>
                    </div>
                    <iframe className="preview" ref={previewIframe} />
                </div>
            </div>
            <style jsx>{`
                .addPadding {
                    padding: 10px 0;
                }
                .container {
                    max-height: 100vh;
                    min-height: 100vh;
                    display: flex;
                    flex-direction: column;
                }
                .editor {
                    flex-grow: 1;
                    display: flex;
                }
                .editor > * {
                    width: 50%;
                }
                .inputs {
                    padding: 20px;
                }
                .preview iframe {
                    height: 100%;
                    width: 100%;
                    border: 0 solid transparent;
                }
                .loading-indicator-container {
                    width: 100%;
                    max-height: 0;
                }
                .loading-indicator {
                    pointer-events: none;
                    position: relative;
                    margin: 0 auto;
                    min-width: 130px;
                    max-width: 130px;
                    min-height: 40px;
                    max-height: 40px;
                    border-radius: 20px;
                    background-color: #303030;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    z-index: 10;
                    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
                    transition: opacity .5s, transform .5s;
                }
                .loading-indicator span {
                    display: inline-block;
                    padding-left: 10px;
                }
            `}</style>

            <style jsx global>{`
                .appContainer .version {
                   display: none;
                }
            `}</style>
        </div>
    )
}
