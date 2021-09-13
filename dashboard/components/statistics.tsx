import dynamic from 'next/dynamic'
import { fetcher } from '../src/auth'
import { formatRFC3339, subDays, startOfDay } from 'date-fns'
import React, { useEffect, useMemo, useState } from 'react'
import { Match } from '../src/types'
import { BarChartProps } from './chartProps'

const BarChart = dynamic<BarChartProps>(() =>
    import('./chart').then(m => m.BarChart),
    { ssr: false },
)

export default function Statistics() {
    const [loading, setLoading] = useState(true)
    const [matches, setMatches] = useState<Array<Match>>([])
    const [dateRange, setDateRangeFromAndTo] = useState<[Date, Date] | undefined>(undefined)

    const fetchAnalytics = async (from: Date, to: Date) => {
        try {
            setLoading(true)
            const fetchedData = await fetcher.fetch(
                `/api/v1/analytics/matches/period/${formatRFC3339(from)}/${formatRFC3339(to)}`
            )
            setMatches(
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

    const matchesPerDayOfPrev7Days = useMemo(() => {
        const resOffset = dateRange?.[1]?.getDay() || 0;

        const res = matches
            .reduce((acc, match) => {
                acc[match.when.getDay()]++
                return acc
            }, [0, 0, 0, 0, 0, 0, 0])
            .map((value, idx) => ({
                label: dayNames[idx],
                value,
            }))

        return [...res.splice(resOffset + 1), ...res]
    }, [matches, dateRange])

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
                    <BarChart
                        data={matchesPerDayOfPrev7Days}
                        singleTooltip="match"
                        multipleTooltip="matches"
                    />
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
                    padding: 20px 20px 10px 20px;
                    overflow: hidden;
                    border-radius: 4px;
                    background-color: #424242;
                }
            `}</style>
        </div>
    )
}
