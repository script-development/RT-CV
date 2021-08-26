import Head from "next/head";
import Dynamic from "next/dynamic"
import React, { useRef } from "react"

const Editor = Dynamic(
    () => import("@monaco-editor/react"),
    { ssr: false }
);

export default function TryMatcher() {
    const editorRef = useRef(null as any)

    const handleEditorWillMount = (monaco: any) => {
        monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
            validate: true,
            schemas: [{
                // For info about how this works see:
                // https://json-schema.org/learn/getting-started-step-by-step.html
                uri: "http://myserver/foo-schema.json",
                fileMatch: ['*'],
                schema: {
                    type: "object",
                    properties: {
                        name: { type: 'string' },
                        p1: {
                            enum: ["v1", "v2"]
                        },
                    },
                    required: ["productId"]
                }
            }]
        })
    }

    const handleEditorDidMount = (editor: any, monaco: any) =>
        editorRef.current = editor

    const getValue = () =>
        editorRef.current.getValue()

    return (
        <div>
            <Head>
                <title>RT-CV home</title>
            </Head>
            <div>
                <Editor
                    height="90vh"
                    defaultLanguage="json"
                    defaultValue={`{"//": "let's write some broken code 😈"}`}
                    theme="vs-dark"
                    beforeMount={handleEditorWillMount}
                    onMount={handleEditorDidMount}
                />
            </div>
        </div>
    )
}
