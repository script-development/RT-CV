import dynamic from 'next/dynamic'
import { fetcher } from '../src/auth'
import { formatRFC3339, subDays, startOfDay } from 'date-fns'
import React, { useEffect, useMemo, useState } from 'react'
import { Match } from '../src/types'
import { BarChartProps } from './chartProps'
import { CircularProgress, Tooltip } from '@mui/material'

const BarChart = dynamic<BarChartProps>(() =>
    import('./chart').then(m => m.BarChart),
    {
        ssr: false,
        loading: () => <ChartLoader />,
    },
)

export default function Statistics() {
    const [loading, setLoading] = useState(true)
    const [matches, setMatches] = useState<Array<Match>>([])
    const [dateRange, setDateRangeFromAndTo] = useState<[Date, Date] | undefined>(undefined)
    const [profilesStats, setProfilesStats] = useState({ total: undefined as undefined | number, usable: undefined as undefined | number })

    const fetchAnalytics = async (from: Date, to: Date) => {
        try {
            setLoading(true)
            const [fetchedData, profilesStats] = await Promise.all([
                fetcher.fetch(`/api/v1/analytics/matches/period/${formatRFC3339(from)}/${formatRFC3339(to)}`),
                fetcher.fetch(`/api/v1/profiles/count`),
            ])
            setProfilesStats(profilesStats)
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
            <div className="char">
                <h3>Matches per day</h3>
                <div className="box">
                    {loading
                        ? <ChartLoader />
                        : <BarChart
                            data={matchesPerDayOfPrev7Days}
                            singleTooltip="match"
                            multipleTooltip="matches"
                        />
                    }
                </div>
            </div>
            <div className="char">
                <h3>Profiles Count</h3>
                <div className="box">
                    <div>
                        <p>Total</p>
                        <h1>{profilesStats.total}</h1>
                    </div>
                    <Tooltip title="The profiles that are active and have a OnMatch property set" placement="top">
                        <div>
                            <p>Used by matcher</p>
                            <h1>{profilesStats.usable}</h1>
                        </div>
                    </Tooltip>
                </div>
            </div>
            <style jsx>{`
                .charts {
                    width: 700px;
                    box-sizing: border-box;
                    max-width: calc(100vw - 20px);
                    display: flex;
                    flex-wrap: wrap;
                    align-items: stretch;
                }
                .char {
                    padding: 10px;
                    flex-grow: 1;
                }
                .box {
                    height: 230px;
                    padding: 20px 20px 10px 20px;
                    overflow: hidden;
                    border-radius: 4px;
                    background-color: #424242;
                    min-width: 300px;
                }
            `}</style>
        </div>
    )
}

function ChartLoader() {
    return (
        <div className="loader">
            <CircularProgress />
            <style jsx>{`
                .loader {
                    width: 285px;
                    height: 200px;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                }
            `}</style>
        </div>
    )
}
