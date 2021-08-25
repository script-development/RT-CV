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
}
