import { LinearProgress } from '@mui/material'
import { ReactNode } from 'react'

interface CardArgs {
    title?: string
    loading?: boolean
    headerContent?: ReactNode
    children?: ReactNode
}

export default function Card({ title, loading, headerContent, children }: CardArgs) {
    return (
        <div className="container">
            {title ? <h3>{title}</h3> : undefined}
            {headerContent ? <div className="cardHeader">
                <div className="cardHeaderContent">
                    {headerContent}
                </div>
                <div className="cardHeaderProgress">
                    <div>{loading ? <LinearProgress /> : undefined}</div>
                </div>
                {children}
            </div> : undefined}
            <style jsx>{`
                .container {
					padding: 10px;
					width: 700px;
					box-sizing: border-box;
					max-width: calc(100vw - 20px);
				}
                .cardHeader {
					overflow: hidden;
					border-top-left-radius: 4px;
					border-top-right-radius: 4px;
					background-color: #424242;
				}
				.cardHeaderContent {
					padding: 10px;
					display: flex;
					justify-content: space-between;
					align-items: center;
				}
				.cardHeaderProgress {
					max-height: 0px;
					transform: translate(0, -4px);
				}
            `}</style>
        </div>
    )
}
