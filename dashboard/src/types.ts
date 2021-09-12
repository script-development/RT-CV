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
    debug: boolean
    desiredProfession: boolean
    domains?: string
    driversLicense: boolean
    educationOrCourse: boolean
    id: string
    keyId: string
    professionExperienced: boolean
    profileId: string
    requestId: string
    when: Date
    yearsSinceEducation: boolean
    yearsSinceWork: boolean
    zipCode: null | DutchZipCode
}

export interface DutchZipCode {
    from: number
    to: number
}
