import { Roles } from './roles'
import { ApiKey } from './types';

class AuthenticatedFetcher {
    private triedToRestoreCredentials = false;
    private apiKeyHashed = '';
    private apiKeyId = '';

    async login(key: string, keyId: string) {
        if (!this.triedToRestoreCredentials) this.tryRestoreCredentials();

        if (key == '')
            throw 'api key not set'
        if (keyId == '')
            throw 'api key id not set'

        this.apiKeyHashed = await this.hash(key)
        this.apiKeyId = keyId
        try {
            const keyInfo = await this.fetch('/api/v1/auth/keyinfo')

            // Check if this key has the required roles
            const hasRequiredRole = keyInfo.roles.some((role: any) => role.role == Roles.Dashboard)
            if (!hasRequiredRole) {
                throw 'key does not have the required role'
            }

            localStorage.setItem('rtcv_api_key', key)
            this.storeCredentials()
        } catch (e) {
            this.apiKeyHashed = ''
            this.apiKeyId = ''
            throw e
        }
    }

    get getApiKeyHashed(): string {
        if (!this.triedToRestoreCredentials) this.tryRestoreCredentials()
        return this.apiKeyHashed
    }

    get getApiKey(): string {
        return localStorage.getItem('rtcv_api_key') || ''
    }

    get getApiKeyId(): string {
        if (!this.triedToRestoreCredentials) this.tryRestoreCredentials()
        return this.apiKeyId
    }

    async post(path: string, data?: any) {
        return await this.fetch(path, "POST", data)
    }

    async put(path: string, data?: any) {
        return await this.fetch(path, "PUT", data)
    }

    async delete(path: string, data?: any) {
        return await this.fetch(path, "DELETE", data)
    }

    async get(path: string) {
        return await this.fetch(path)
    }

    getAPIPath(path: string, isWebsocketURL?: boolean): string {
        if (window.location.host != 'localhost:3000' && window.location.host != '127.0.0.1:3000')
            return path

        else if (isWebsocketURL ? (path.indexOf('wss://') == 0 || path.indexOf('ws://') == 0) : (path.indexOf('https://') == 0 || path.indexOf('http://') == 0))
            return path

        const protocol = !isWebsocketURL ?
            window.location.protocol
            : window.location.protocol == 'https:'
                ? 'wss:'
                : 'ws:'

        return `${protocol}//${window.location.hostname}:4000${path}`
    }

    async fetch(path: string, method: 'POST' | 'GET' | 'PUT' | 'DELETE' = 'GET', data?: any) {
        try {
            const r = await this.fetchNoJsonMarshal(path, method, data)
            const resJsonData = await r.json()
            const error = resJsonData?.error

            if (r.status == 401) {
                if (error && error != "you do not have auth roles required to access this route")
                    // redirect to login screen as something with the credentials is going wrong
                    location.pathname = '/login'

                throw error
            }

            if (error || r.status >= 400) {
                throw error
            }

            return resJsonData
        } catch (e: any) {
            throw e
        }
    }

    async fetchNoJsonMarshal(path: string, method: 'POST' | 'GET' | 'PUT' | 'DELETE' = 'GET', data?: any) {
        try {
            if (path.replace(/http(s?):\/\//, '').indexOf('//') != -1)
                throw 'invalid path, path cannot contains empty parts, path: ' + path

            if (!this.triedToRestoreCredentials) this.tryRestoreCredentials();

            const url = this.getAPIPath((path[0] != '/' ? '/' : '') + path)
            const args: RequestInit = {
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": this.authorizationHeader,
                },
                method: method,
                body: data ? JSON.stringify(data) : undefined,
            }

            return await fetch(url, args)
        } catch (e: any) {
            throw e
        }
    }

    // Returns true if the apiKey and apiKeyId is set
    private tryRestoreCredentials(): boolean {
        this.triedToRestoreCredentials = true;
        this.apiKeyId = localStorage.getItem('rtcv_api_key_id') || ''
        this.apiKeyHashed = localStorage.getItem('rtcv_api_key_hashed') || ''
        return !!(this.apiKeyHashed && this.apiKeyId)
    }

    private storeCredentials() {
        localStorage.setItem('rtcv_api_key_id', this.apiKeyId)
        localStorage.setItem('rtcv_api_key_hashed', this.apiKeyHashed)
    }

    get authorizationHeader(): string {
        return `Basic ${this.authorizationValue}`
    }

    get authorizationValue(): string {
        if (!this.triedToRestoreCredentials) this.tryRestoreCredentials()
        return `${this.apiKeyId}:${this.apiKeyHashed}`
    }

    private async hash(data: string): Promise<string> {
        return Buffer.from(
            new Uint8Array(
                await crypto.subtle.digest(
                    'SHA-512',
                    new TextEncoder().encode(data),
                ),
            ),
        ).toString('hex')
    }
}

export const fetcher = new AuthenticatedFetcher()

export const getKeys = (): Promise<Array<ApiKey>> => fetcher.get('/api/v1/keys')
