import Head from "next/head"
import dynamic from 'next/dynamic'
import Header from '../components/header'
import { useEffect, useRef, useState } from 'react'
import { fetcher } from '../src/auth'
import { Checkbox, FormControl, FormLabel, Select, MenuItem, InputLabel, CircularProgress, Tooltip, TextField } from "@material-ui/core"

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
                        <SwitchableInput
                            title="Header background color"
                            valueToCheck={options.headerColor}
                            setValue={enabled => setOptions(v => ({ ...v, headerColor: enabled ? '#ff0000' : undefined }))}
                        >
                            <input
                                value={options.headerColor}
                                type="color"
                                onChange={e => setOptions(v => ({ ...v, headerColor: e.target.value }))}
                            />
                        </SwitchableInput>
                        <SwitchableInput
                            title="Sub header background color"
                            valueToCheck={options.subHeaderColor}
                            setValue={enabled => setOptions(v => ({
                                ...v,
                                subHeaderColor: enabled ? '#ff0000' : undefined,
                            }))}
                        >
                            <input
                                value={options.subHeaderColor}
                                type="color"
                                onChange={e => setOptions(v => ({ ...v, subHeaderColor: e.target.value }))}
                            />
                        </SwitchableInput>
                        <SwitchableInput
                            title="Fonts"
                            valueToCheck={options.fontRegular}
                            setValue={enabled => setOptions(v => ({
                                ...v,
                                fontRegular: enabled ? 'OpenSans' : undefined,
                                fontHeader: enabled ? 'OpenSans' : undefined,
                            }))}
                        >
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
                        </SwitchableInput>
                        <SwitchableInput
                            title="Layout"
                            valueToCheck={options.style}
                            setValue={enabled => setOptions(v => ({ ...v, style: enabled ? 'style_1' : undefined }))}
                        >
                            <FormControl>
                                <InputLabel id="style-label">Style</InputLabel>
                                <Select labelId="style-label" value={options.style} onChange={e => setOptions(v => ({ ...v, style: e.target.value as string | undefined }))}>
                                    {styleOptions.map((entry, i) => <MenuItem key={i} value={entry.value}>{entry.label}</MenuItem>)}
                                </Select>
                            </FormControl>
                        </SwitchableInput>
                        <SwitchableInput
                            title="Logo image url"
                            valueToCheck={options.logoImageUrl}
                            setValue={enabled => setOptions(v => ({ ...v, logoImageUrl: enabled ? '' : undefined }))}
                        >
                            <TextField
                                id="logo-image-url"
                                label="Logo image url"
                                variant="filled"
                                value={options.logoImageUrl}
                                onChange={e => setOptions(v => ({ ...v, logoImageUrl: e.target.value }))}
                            />
                        </SwitchableInput>
                        <SwitchableInput
                            title="Company name"
                            valueToCheck={options.companyName}
                            setValue={enabled => setOptions(v => ({ ...v, companyName: enabled ? 'a company b.v.' : undefined }))}
                        >
                            <TextField
                                id="company-name"
                                label="Company name"
                                variant="filled"
                                value={options.companyName}
                                onChange={e => setOptions(v => ({ ...v, companyName: e.target.value }))}
                            />
                        </SwitchableInput>
                        <SwitchableInput
                            title="Company address"
                            valueToCheck={options.companyAddress}
                            setValue={enabled => setOptions(v => ({ ...v, companyAddress: enabled ? '9977AB\nSome street 15, a city' : undefined }))}
                        >
                            <TextField
                                id="company-address"
                                label="Company name"
                                variant="filled"
                                multiline
                                rows={4}
                                value={options.companyAddress}
                                onChange={e => setOptions(v => ({ ...v, companyAddress: e.target.value }))}
                            />
                        </SwitchableInput>
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
                .container {
                    max-height: 100vh;
                    min-height: 100vh;
                    display: flex;
                    flex-direction: column;
                }
                .editor {
                    flex-grow: 1;
                    display: flex;
                    overflow: hidden;
                }
                .editor > * {
                    width: 50%;
                    flex-grow: 1;
                }
                .inputs {
                    overflow-y: scroll;
                    overflow-x: hidden;
                    box-sizing: border-box;
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

interface SwitchableInputProps {
    title: string
    valueToCheck: any
    setValue: (enabled: boolean) => void
    children?: React.ReactNode
}

function SwitchableInput({ title, valueToCheck, setValue, children }: SwitchableInputProps) {
    return (
        <div>
            <FormLabel>{title}</FormLabel>
            <div>
                <Tooltip title="Enable / Disable this configuration field">
                    <Checkbox
                        checked={valueToCheck !== undefined}
                        onChange={() => setValue(valueToCheck === undefined ? true : false)}
                        color="primary"
                    />
                </Tooltip>
                {valueToCheck !== undefined ? children : undefined}
            </div>
        </div>
    )
}
