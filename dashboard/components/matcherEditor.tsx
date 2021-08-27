import Monokai from 'monaco-themes/themes/Monokai.json'
import MonacoEditor from "@monaco-editor/react"
import React, { useState, useEffect, useRef, CSSProperties } from 'react';
import { CircularProgress, IconButton } from '@material-ui/core';
import PlayArrow from '@material-ui/icons/PlayArrow'
import { fetcher } from '../src/auth';
import { parse } from 'jsonc-parser'

interface MatcherEditorProps {
    style?: CSSProperties
}

const schemaUrl = () => location.origin + '/api/v1/schema/cv'

export default function MatcherEditor({ style }: MatcherEditorProps) {
    const [cvSchema, setCvSchema] = useState(undefined)
    const [inputEditorRef, outputEditorRef] = [useRef(null as any), useRef(null as any)]
    const [loading, setLoading] = useState(false)

    const handleInputEditorWillMount = (monaco: any) => {
        monaco.editor.defineTheme('monokai', Monokai)
        monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
            validate: true,
            schemas: [{
                // For info about how this works see:
                // https://json-schema.org/learn/getting-started-step-by-step.html
                uri: schemaUrl(),
                fileMatch: ['*'],
                schema: cvSchema,
            }],
        })
    }

    const handleOutputEditorWillMount = (monaco: any) =>
        monaco.editor.defineTheme('monokai', Monokai)

    const handleInputEditorDidMount = (editor: any, monaco: any) =>
        inputEditorRef.current = editor

    const handleOutputEditorDidMount = (editor: any, monaco: any) =>
        outputEditorRef.current = editor


    const fetchSchema = async () => {
        const r = await fetch('/api/v1/schema/cv')
        const jsonData = await r.json()
        jsonData.$id = schemaUrl()
        setCvSchema(jsonData)
    }

    const execute = async () => {
        try {
            setLoading(true)

            // We use jsonc here to parse the JSON and convert it back to json
            // This way we can use comments
            const parsedCv = parse(inputEditorRef.current.getValue())
            const requestValue = {
                cv: parsedCv,
                debug: true,
            }

            const startTime = performance.now()
            const res = await fetcher.post('/api/v1/scraper/scanCV', requestValue)
            const endTime = performance.now()
            outputEditorRef.current.setValue(`// api call took ${Math.round(endTime - startTime)} milliseconds with ${res.length} results\n${JSON.stringify(res, null, '\t')}`)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchSchema() }, [])

    return (
        <div className="root" style={style}>
            <div className="editor input">
                {cvSchema
                    ? <MonacoEditor
                        height="100%"
                        defaultLanguage="json"
                        defaultValue={`{\n\t// Press ctrl + space to start hacking\n\t\n}`}
                        theme="monokai"
                        beforeMount={handleInputEditorWillMount}
                        onMount={handleInputEditorDidMount}
                        options={{ fontSize: 18, minimap: { enabled: false } }}
                    />
                    : ''
                }
            </div>
            <div className="separator">
                <div className="playButtonContainer">
                    <IconButton
                        onClick={execute}
                        disabled={loading}
                        color="primary"
                        style={{ backgroundColor: '#ff6e40', color: 'black' }}
                    >
                        <PlayArrow />
                    </IconButton>
                </div>
            </div>
            <div className="editor output">
                <MonacoEditor
                    height="100%"
                    defaultLanguage="json"
                    defaultValue={`// press the play button to see the api result`}
                    theme="monokai"
                    beforeMount={handleOutputEditorWillMount}
                    onMount={handleOutputEditorDidMount}
                    options={{ fontSize: 18, minimap: { enabled: false }, readOnly: true }}
                />
                <div className={"loader" + (loading ? ' loading' : '')}>
                    {loading ? <CircularProgress /> : ''}
                </div>
            </div>
            <style jsx>{`
                .root {
                    display: flex;
                    flex-wrap: nowrap;
                    justify-content: space-between;
                    align-items: stretch;
                }
                .editor {
                    width: calc(50% - 10px);
                }
                .editor .loader {
                    pointer-events: none;
                    position: relative;
                    height: 100%;
                    width: 100%;
                    top: -100%;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    opacity: 0;
                    transition: background-color 0.2s, opacity 0.2s;
                }
                .editor .loader.loading {
                    background-color: rgba(0,0,0,0.5);
                    opacity: 1;
                }
                .separator {
                    z-index: 5;
                    width: 20px;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                }
            `}</style>
        </div>
    )
}
