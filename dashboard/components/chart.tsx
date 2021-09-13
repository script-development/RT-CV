import React, { useEffect, useState } from "react"
import { primaryColor, secondaryColor } from '../src/theme'
import { Tooltip } from '@material-ui/core'
import { BarChartProps } from './chartProps'

export function BarChart({ data, singleTooltip, multipleTooltip }: BarChartProps) {
    const [max, setMax] = useState(1)

    useEffect(() => {
        setMax(Math.max(...data.map(b => b.value)) || 1)
    }, [data])

    return (
        <div className="barChart">
            <div className="info">
                <div className="values">
                    <div>{max}</div>
                    {max > 50 ? <div>{Math.floor(max / 4 * 3)}</div> : undefined}
                    {max > 1 ? <div>{max > 50 ? Math.floor(max / 2) : max / 2}</div> : undefined}
                    {max > 50 ? <div>{Math.floor(max / 4)}</div> : undefined}
                    <div>0</div>
                </div>
            </div>
            {data.map((bar, idx) => {
                return <div key={idx} className="bar">
                    <Tooltip title={bar.value + ' ' + (bar.value == 1 ? singleTooltip : multipleTooltip)}>
                        <div className="valueContainer">
                            <div
                                className="value"
                                style={{
                                    height: (100 / max * bar.value) + '%',
                                    backgroundColor: bar.value == max ? secondaryColor : primaryColor,
                                }}
                            />
                        </div>
                    </Tooltip>
                    <div className="label">{bar.label}</div>
                </div>
            })}
            <style jsx>{`
                .barChart {
                    height: 200px;
                    display: flex;
                    align-items: stretch;
                }
                .info {
                    width: 40px;
                    padding-right: 5px;
                }
                .bar {
                    width: 35px;
                }
                .bar .valueContainer, .info .values {
                    height: calc(100% - 40px);
                    margin: 5px 0;
                }
                .info .values {
                    display: flex;
                    flex-direction: column;
                    justify-content: space-between;
                }
                .info .values div {
                    color: rgba(255, 255, 255, 0.7);
                    text-align: right;
                    font-weight: bold;
                }
                .bar .valueContainer {
                    display: flex;
                    justify-content: center;
                    align-items: flex-end;                }
                .bar .value {
                    width: 100%;
                    background-color: white;
                    box-sizing: border-box;
                    margin: 1px;
                    border-radius: 4px;
                }
                .bar .label {
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    text-align: center;
                    font-weight: bold;
                    color: rgba(255, 255, 255, 0.7);
                    height: 30px;
                }
            `}</style>
        </div>
    )
}
