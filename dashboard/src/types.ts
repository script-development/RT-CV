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

export interface CV {
    referenceNumber: string
    educations?: Array<Education>
    courses?: Array<Course>
    preferredJobs?: Array<string>
    languages?: Array<Language>
    personalDetails?: PersonalDetails
    driversLicenses?: Array<string>
}

export interface PersonalDetails {
    initials?: string
    firstName?: string
    surNamePrefix?: string
    surName?: string
    dob?: string
    gender?: string
    streetName?: string
    houseNumber?: string
    houseNumberSuffix?: string
    zip?: string
    city?: string
    country?: string
    phoneNumber?: string
    email?: string
}

export interface Education {
    name: string
    description: string
    institute: string
    isCompleted: boolean
    hasDiploma: boolean
    startDate: string
    endDate: string
}

export interface Course {
    name: string
    institute: string
    startDate: string
    endDate: string
    isCompleted: boolean
    description: string
}

export interface Language {
    name: string
    levelSpoken: LanguageLevel
    levelWritten: LanguageLevel
}

export enum LanguageLevel {
    LanguageLevelUnknown = 0,
    LanguageLevelReasonable = 1,
    LanguageLevelGood = 2,
    LanguageLevelExcellent = 3,
}

export function LanguageLevelToString(lev: LanguageLevel): string {
    switch (lev) {
        case LanguageLevel.LanguageLevelUnknown:
            return 'Unknown'
        case LanguageLevel.LanguageLevelReasonable:
            return 'Reasonable'
        case LanguageLevel.LanguageLevelGood:
            return 'Good'
        case LanguageLevel.LanguageLevelExcellent:
            return 'Excellent'
        default:
            return 'Unknown'
    }
}
