import { fetcher } from '../src/auth'
import React, { useEffect, useState } from 'react'
import { Tooltip } from '@mui/material'

export default function Statistics() {
    const [profilesStats, setProfilesStats] = useState({ total: undefined as undefined | number, usable: undefined as undefined | number })

    useEffect(() => {
        fetcher.fetch(`/api/v1/profiles/count`).then(setProfilesStats)
    }, [])

    return (
        <div className="charts">
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
                    width: 300px;
                }
            `}</style>
        </div>
    )
}
