export interface ApiKey {
    domains: Array<string>
    name: string,
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
    valueStructure: SecretValueStructure
}

export enum SecretValueStructure {
    Free = 'free',
    StrictUser = 'strict-user',
    StrictUsers = 'strict-users',
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

export interface OnMatchHook {
    id: string
    keyId: string
    url: string
    method: string
    addHeaders: Array<{ key: string, value: Array<string> }>
    stopRemainingActions: boolean
}
