import Monokai from 'monaco-themes/themes/Monokai.json'
import MonacoEditor from "@monaco-editor/react"
import React, { MutableRefObject, useState, useEffect, useRef } from 'react';
import { CircularProgress, IconButton } from '@material-ui/core';
import PlayArrow from '@material-ui/icons/PlayArrow'
import { fetcher } from '../src/auth';
import { parse } from 'jsonc-parser'
import { Monaco } from '@monaco-editor/react'
import { editor } from 'monaco-editor';

const schemaUrl = () => location.origin + '/api/v1/schema/cv'

interface MatcherEditorExposedValues {
    inputEditorRef: MutableRefObject<any>,
    outputEditorRef: MutableRefObject<any>,
}

interface MatcherEditorProps {
    expose?: (values: MatcherEditorExposedValues) => void
}

export default function MatcherEditor({ expose }: MatcherEditorProps) {
    const [cvSchema, setCvSchema] = useState(undefined)
    const [inputValue, setInputValue] = useState(`{\n\t// Press ctrl + space to start hacking\n\t\n}`)
    const [outputValue, setOutputValue] = useState(`// press the play button to see the api result`)
    const [loading, setLoading] = useState(false)
    const [inputEditorRef, outputEditorRef] = [useRef(null as any), useRef(null as any)]

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

    const handleInputEditorDidMount = (editor: editor.IStandaloneCodeEditor, _: Monaco) => {
        inputEditorRef.current = editor
        editor.setPosition({ lineNumber: 0, column: 0 })
    }

    const handleOutputEditorDidMount = (editor: editor.IStandaloneCodeEditor, _: Monaco) => {
        outputEditorRef.current = editor
        editor.setPosition({ lineNumber: 0, column: 0 })
    }

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

    useEffect(() => {
        fetchSchema()
        expose?.({
            inputEditorRef,
            outputEditorRef
        })
    }, [])

    return (
        <div className="root">
            <div className="editor input">
                {cvSchema
                    ? <MonacoEditor
                        height="100%"
                        defaultLanguage="json"
                        value={inputValue}
                        onChange={v => setInputValue(v || '')}
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
                    value={outputValue}
                    onChange={v => setOutputValue(v || '')}
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
                    height: 100vh;
                    width: 100vw;
                    padding-bottom: 50px;
                    box-sizing: border-box;

                }
                .editor {
                    width: calc(50% - 10px);
                    height: 100%;
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
