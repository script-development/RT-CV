import { useMemo } from 'react'
import SyntaxHighlighter from 'react-syntax-highlighter'
import { monokaiSublime } from 'react-syntax-highlighter/dist/cjs/styles/hljs'

export default function JSONCode({ json }: { json: any }) {
    const jsonString = useMemo(
        () => JSON.stringify(json, null, 4),
        [json],
    )

    return (
        <div className="code">
            <SyntaxHighlighter
                wrapLongLines={true}
                language="json"
                style={monokaiSublime}
            >{jsonString}</SyntaxHighlighter>
        </div>
    )
}
