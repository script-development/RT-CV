export interface ApiKey {
    domains: Array<string>
    enabled: boolean
    id: string
    key: string
    roles: number
    system: boolean
}

export interface Secret {
    id: string
    keyId: string
    key: string
    description: string
}

export interface Match {
    id: string
    keyId: string
    profileId: string
    requestId: string
    when: Date

    debug?: boolean

    domains?: string
    yearsSinceWork?: number
    yearsSinceEducation?: number
    education?: string
    course?: string
    desiredProfession?: string
    professionExperienced?: boolean
    driversLicense?: boolean
    zipCode?: null | DutchZipCode
}

export interface DutchZipCode {
    from: number
    to: number
}
