import { fetcher } from '../src/auth'
import { formatRFC3339, subDays, startOfDay } from 'date-fns'
import React, { useEffect, useMemo, useState } from 'react'
import { Match } from '../src/types'
import { primaryColor } from '../src/theme'
import { Tooltip } from '@material-ui/core'

export default function Statistics() {
    const [loading, setLoading] = useState(true)
    const [data, setData] = useState<Array<Match>>([])
    const [dateRange, setDateRangeFromAndTo] = useState<[Date, Date] | undefined>(undefined)

    const fetchAnalytics = async (from: Date, to: Date) => {
        try {
            setLoading(true)
            const fetchedData = await fetcher.fetch(
                `/api/v1/analytics/matches/period/${formatRFC3339(from)}/${formatRFC3339(to)}`
            )
            setData(
                fetchedData.map((entry: any) => ({
                    ...entry,
                    when: new Date(entry.when),
                }))
            )
        } finally {
            setLoading(false)
        }
    }

    const dayNames = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Su'];

    const formattedData = useMemo(() => {
        const resOffset = dateRange?.[1]?.getDay() || 0;

        const res = data
            .reduce((acc, match) => {
                acc[match.when.getDay()]++
                return acc
            }, [0, 0, 0, 0, 0, 0, 0])
            .map((value, idx) => ({
                dayName: dayNames[idx],
                value,
            }))

        return {
            list: [...res.splice(resOffset + 1), ...res],
            max: res.reduce((acc, item) => item.value > acc ? item.value : acc, 0),
        }
    }, [data, dateRange])

    useEffect(() => {
        const from = startOfDay(subDays(new Date(), 6));
        const to = new Date();
        setDateRangeFromAndTo([from, to])
        fetchAnalytics(from, to)
    }, [])

    return (
        <div className="charts">
            <div className="chartContainer">
                <h3>Matches per day</h3>
                <div className="chart">
                    <div className="info">
                        <div className="max">{formattedData.max}</div>
                    </div>
                    {formattedData.list.map((item, idx) =>
                        <div key={idx} className="day">
                            {!item.value
                                ? <Tooltip title="0 matches" enterDelay={1000}>
                                    <div className="barContainer">
                                        <div className="bar" />
                                    </div>
                                </Tooltip>
                                : <div className="barContainer">
                                    <Tooltip title={item.value + " matches"} enterDelay={1000}>
                                        <div
                                            className="bar"
                                            style={{
                                                height: (100 / formattedData.max * item.value) + '%',
                                                backgroundColor: primaryColor,
                                            }}
                                        />
                                    </Tooltip>
                                </div>
                            }
                            <div className="dayName">{item.dayName}</div>
                        </div>
                    )}
                </div>
            </div>
            <style jsx>{`
                .charts {
                    padding: 10px;
                    width: 700px;
                    box-sizing: border-box;
                    max-width: calc(100vw - 20px);
                    display: flex;
                }
                .chart {
                    height: 200px;
                    display: flex;
                    align-items: stretch;

                    padding: 10px;
                    overflow: hidden;
                    border-radius: 4px;
                    background-color: #424242;
                }
                .info, .day {
                    width: 35px;
                }
                .day .barContainer {
                    height: calc(100% - 30px);
                }
                .day .bar {
                    background-color: white;
                    margin: 1px;
                    border-radius: 4px;
                    box-sizing: border-box;
                }
                .day .dayName, .info .max {
                    height: 30px;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    text-align: center;
                    font-weight: bold;
                }
            `}</style>
        </div>
    )
}
