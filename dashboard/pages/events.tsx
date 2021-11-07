import Head from 'next/head'
import { ReactNode, useEffect, useState } from 'react'
import { Icon, Typography } from '@material-ui/core'
import { fetcher } from '../src/auth'
import Header from '../components/header'
import { CV, LanguageLevelToString } from '../src/types'
import Check from '@material-ui/icons/Check'
import Close from '@material-ui/icons/Close'

function getWebsocketUrl() {
    const url = fetcher.getAPIPath(`/api/v1/events/ws/${fetcher.authorizationValue}`, true)
    return url[0] == '/'
        ? `ws${location.protocol == 'https:' ? 's' : ''}://${location.host}${url}`
        : url
}

export default function Events() {
    const [connectionStatus, setConnectionStatus] = useState({
        connected: false,
        msg: '',
    })
    const [events, setEvents] = useState<Array<any>>(JSON.parse(testValue))

    const connectToSocket = () => {
        try {
            const socket = new WebSocket(getWebsocketUrl())

            let open = true
            const close = () => {
                if (!open) { return }
                setConnectionStatus({
                    connected: false,
                    msg: 'trying to reconnect to websocket in 5 seconds',
                })
                open = false
                socket.onmessage = null
                socket.onopen = null
                socket.onerror = null
                socket.onclose = null
                setTimeout(() => {
                    setConnectionStatus({
                        connected: false,
                        msg: 'reconnecting..',
                    })
                    connectToSocket()
                }, 5000)
            }

            socket.onmessage = (ev: MessageEvent<any>) => {
                console.log('received message', ev.data)
                try {
                    const newMsg = JSON.parse(ev.data)
                    if (newMsg.kind == "recived_cv") {
                        setEvents(prev => [newMsg, ...prev])
                    }
                } catch (e) { }
            }
            socket.onopen = () => {
                setConnectionStatus({
                    connected: true,
                    msg: '',
                })
                console.log('connected to websocket')
            }
            socket.onerror = () => {
                console.log('disconnected due to websocket error')
                close()
            }
            socket.onclose = () => {
                console.log('websocket connection closed')
                close()
            }
            return () => {
                open = false;
                socket.close(1000, 'navigating to different route')
            }
        } catch (e) {
            console.error(e)
            setConnectionStatus({
                connected: false,
                msg: 'unable to connect with websocket',
            })
            return () => { }
        }
    }

    useEffect(connectToSocket, [])

    return (
        <div>
            <Header />

            <Head><title>RT-CV events</title></Head>

            <div className="status">
                <div className="dot" style={{ backgroundColor: connectionStatus.connected ? '#8bc34a' : '#ff5722' }} />
                {connectionStatus.connected
                    ? 'connected to server'
                    : connectionStatus.msg
                        ? 'disconnected, ' + connectionStatus.msg
                        : 'disconnected'
                }
            </div>

            <div className="eventsList">
                <h2>Events</h2>
                <Typography variant="body2">
                    All events shown here are available as long as you have the tab open.
                    Once you reload the tab the events are gone
                </Typography>
                {events.map((ev, idx) =>
                    <Event
                        key={idx}
                        event={ev}
                        isLast={idx == events.length - 1}
                    />
                )}
            </div>

            <style jsx>{`
                .eventsList {
                    padding: 20px;
                }
                .status {
                    padding: 10px;
                    text-align: center;
                    color: rgba(255, 255, 255, 0.7);
                }
                .status .dot {
                    display: inline-block;
                    height: 10px;
                    width: 10px;
                    background-color: white;
                    border-radius: 5px;
                    margin-right: 4px;
                }
            `}</style>
        </div>
    )
}

interface EventProps {
    event: {
        data: CV,
        kind: 'recived_cv',
    }
    isLast: boolean
}

function Event({ event, isLast }: EventProps) {
    const { data } = event

    const name = (
        (data.personalDetails?.firstName || data.personalDetails?.initials || '')
        + (data.personalDetails?.surNamePrefix ? ' ' + data.personalDetails?.surNamePrefix : '')
        + (data.personalDetails?.surName ? ' ' + data.personalDetails?.surName : '')
    ).trim()

    const { zip, city, streetName, houseNumber, houseNumberSuffix, country } = (data.personalDetails || {})

    return (
        <div className="event">
            <div className="decoration">
                <Icon fontSize="large">contact_page</Icon>
                {isLast ? undefined : <div className="lineToNextEvent" />}
            </div>
            <div className="content">
                {data.referenceNumber ? <p>{data.referenceNumber}</p> : undefined}
                <h1>{name ? name : <span className="undefined">Unknown name</span>}</h1>
                {zip || city || streetName || houseNumber || houseNumberSuffix || country ?
                    <p>
                        {zip
                            ? <a target="_blank" rel="noopener noreferrer" href={"https://www.google.com/maps/search/" + zip}>{zip}</a>
                            : undefined
                        }
                        {zip && city ? ' - ' : undefined}
                        {city}
                        {city && streetName ? ' ' : undefined}
                        {streetName ? streetName + ' ' + houseNumber + (houseNumberSuffix ? ' ' + houseNumberSuffix : '') : undefined}
                        {streetName && country ? ' - ' : undefined}
                        {country}
                    </p>
                    : undefined}

                <div className="detailedInfo">
                    {data.preferredJobs ?
                        <EventDetailsSection
                            icon="work"
                            title={data.preferredJobs.length > 1 ? "Preferred jobs" : "Preferred job"}
                        >
                            {data.preferredJobs.map((job, idx) =>
                                <p key={idx}>{job}</p>
                            )}
                        </EventDetailsSection>
                        : undefined}

                    {data.educations ?
                        <EventDetailsSection
                            icon="school"
                            title={data.educations.length > 1 ? "Educations" : "Education"}
                        >
                            {data.educations.map((education, idx) =>
                                <div key={idx} className="listItem">
                                    {education.name}
                                    <div className="checklist">
                                        <div>
                                            <IsOk>{education.isCompleted}</IsOk>
                                            <span>{education.isCompleted ? 'Compleet' : 'Niet Compleet'}</span>
                                        </div>
                                        <div>
                                            <IsOk>{education.hasDiploma}</IsOk>
                                            <span>{education.hasDiploma ? 'Heeft diploma' : 'Geen diploma'}</span>
                                        </div>
                                    </div>
                                </div>
                            )}
                        </EventDetailsSection>
                        : undefined}

                    {data.courses ?
                        <EventDetailsSection
                            icon="school"
                            title={data.courses.length > 1 ? "Courses" : "Course"}
                        >
                            {data.courses.map((course, idx) =>
                                <div key={idx} className="listItem">
                                    {course.name}
                                    <div className="checklist">
                                        <div>
                                            <IsOk>{course.isCompleted}</IsOk>
                                            <span>{course.isCompleted ? 'Compleet' : 'Niet Compleet'}</span>
                                        </div>
                                    </div>
                                </div>
                            )}
                        </EventDetailsSection>
                        : undefined}

                    {data.driversLicenses ?
                        <EventDetailsSection
                            icon="drive_eta"
                            title={data.driversLicenses.length > 1 ? "Drivers licenses" : "Driver license"}
                        >
                            {data.driversLicenses.map((license, idx) =>
                                <span key={idx}>{license + ' '}</span>
                            )}
                        </EventDetailsSection>
                        : undefined}

                    {data.languages ?
                        <EventDetailsSection
                            icon="translate"
                            title={data.languages.length > 1 ? "Languages" : "Language"}
                        >
                            {data.languages.map((language, idx) =>
                                <div key={idx} className="listItem">
                                    {language.name} - Spoken: <b>{LanguageLevelToString(language.levelSpoken)}</b>, Written: <b>{LanguageLevelToString(language.levelWritten)}</b>
                                </div>
                            )}
                        </EventDetailsSection>
                        : undefined}
                </div>
            </div>
            <style jsx>{`
                .event {
                    display: flex;
                    padding: 5px;
                }
                .decoration {
                    min-width: 40px;
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                }
                .decoration .lineToNextEvent {
                    margin-top: 10px;
                    width: 2px;
                    flex-grow: 1;
                    border-radius: 1px;
                    background-color: rgba(255,255,255,0.6);
                }
                h1 .undefined {
                    color: rgba(255, 255, 255, 0.4)
                }
                .content {
                    flex-grow: 1;
                }
                .detailedInfo {
                    width: 100%;
                        display: grid;
                        grid-template-columns: repeat(var(--detailed-info-rows, 1), 1fr);
                }
                .listItem {
                    margin-bottom: 5px;
                }
                .checklist > * {
                    margin-right: 10px;
                    display: inline-flex;
                    align-items: center;
                    color: rgba(255, 255, 255, 0.7);
                }
                .checklist > * > span {
                    margin-left: 3px;
                }
                @media screen and (min-width: 900px) {
                    .detailedInfo { --detailed-info-rows: 2; }
                }
                @media screen and (min-width: 1300px) {
                    .detailedInfo { --detailed-info-rows: 3; }
                }
                @media screen and (min-width: 1700px) {
                    .detailedInfo { --detailed-info-rows: 4; }
                }
            `}</style>
        </div>
    )
}

function IsOk({ children }: { children: boolean }) {
    return children
        ? <Check style={{ color: 'lightgreen' }} />
        : <Close style={{ color: 'red' }} />
}

interface EventDetailsSectionParams {
    icon: string
    title: string
    children?: ReactNode
}

function EventDetailsSection({ icon, title, children }: EventDetailsSectionParams) {
    return (
        <div className="details">
            <div className="icon"><Icon>{icon}</Icon></div>
            <div className="content">
                <h3>{title}</h3>
                {children}
            </div>
            <style jsx>{`
                .details {
                    margin-top: 10px;
                    display: flex;
                }
                .icon {
                    margin-right: 10px;
                    display: inline-block;
                    min-height: 40px;
                    max-height: 40px;
                    min-width: 40px;
                    max-width: 40px;
                    background-color: white;
                    border-radius: 8px;
                    color: black;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                }
                h3 {
                    margin-bottom: 0;
                    padding-bottom: 0;
                }
            `}</style>
        </div>
    )
}

const testValue = '[]'

