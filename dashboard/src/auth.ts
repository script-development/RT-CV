import { Roles } from './roles'

function randomString(length: number): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let str = '';
    for (let i = 0; i < length; i++) {
        str += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return str;
};

class AuthenticatedFetcher {
    private triedToRestoreCredentials = false;
    private apiKey = '';
    private apiKeyId = '';
    private salt = '';

    private serverSeed = '';

    private rollingKey = new Uint8Array;

    async login(key: string, keyId: string) {
        if (!this.triedToRestoreCredentials) this.tryRestoreCredentials();

        if (key == '')
            throw 'api key not set'
        if (keyId == '')
            throw 'api key id not set'

        this.apiKey = key
        this.apiKeyId = keyId
        await this.refreshSeedAndRollingKey()
        const keyInfo = await this.fetch('/api/v1/auth/keyinfo')

        // Check if this key has the required roles
        const [informationObtainerRole, controllerRole] = keyInfo.roles.reduce((acc: [boolean, boolean], role: any) => {
            if (role.role == Roles.InformationObtainer) acc[0] = true
            if (role.role == Roles.Controller) acc[1] = true
            return acc
        }, [false, false])
        if (!informationObtainerRole || !controllerRole) {
            this.apiKey = '';
            this.apiKeyId = '';
            throw 'key does not have the required roles';
        }

        this.storeCredentials()
    }

    get getApiKey(): string {
        if (!this.triedToRestoreCredentials) this.tryRestoreCredentials();

        return this.apiKey
    }

    get getApiKeyId(): string {
        if (!this.triedToRestoreCredentials) this.tryRestoreCredentials();

        return this.apiKeyId
    }

    async fetch(path: string) {
        if (!this.triedToRestoreCredentials) this.tryRestoreCredentials();

        let authHeader = await this.authorizationHeader()
        const url = (path[0] != '/' ? '/' : '') + path
        const args = (authHeader: string) => ({
            headers: {
                "Content-Type": "application/json",
                "Authorization": authHeader,
            },
        })

        let r = await fetch(url, args(authHeader))
        if (r.status == 401) {
            // Firstly lets just re-try it
            r = await fetch(url, args(authHeader))
            if (r.status == 401) {
                // Retrying did not work, lets refresh the seed and rolling key and try again
                await this.refreshSeedAndRollingKey()
                authHeader = await this.authorizationHeader()

                r = await fetch(url, args(authHeader))
                if (r.status == 401) {
                    // FIXME redirect to login screen as something with the credentials is going wrong
                    console.log(this)

                    const resData = await r.json()
                    throw resData.error
                }
                this.storeCredentials()
            }
        }

        return await r.json()
    }

    // Returns true if the apiKey and apiKeyId is set
    private tryRestoreCredentials(): boolean {
        this.triedToRestoreCredentials = true;
        this.apiKey = localStorage.getItem('rtcv_api_key') || ''
        this.apiKeyId = localStorage.getItem('rtcv_api_key_id') || ''
        this.salt = localStorage.getItem('rtcv_salt') || ''
        this.serverSeed = localStorage.getItem('rtcv_server_seed') || ''
        const rollingKeyHex = localStorage.getItem('rtcv_rolling_key') || ''
        if (rollingKeyHex != '')
            this.rollingKey = new Uint8Array(Buffer.from(rollingKeyHex, 'hex'))

        return !!(this.apiKey && this.apiKeyId)
    }

    private storeCredentials() {
        localStorage.setItem('rtcv_api_key', this.apiKey)
        localStorage.setItem('rtcv_api_key_id', this.apiKeyId)
        localStorage.setItem('rtcv_salt', this.salt)
        localStorage.setItem('rtcv_server_seed', this.serverSeed)
        this.storeRollingKey()
    }

    private storeRollingKey() {
        localStorage.setItem('rtcv_rolling_key', this.rollingKeyHex)
    }

    private get rollingKeyHex(): string {
        return (this.rollingKey.length == 0)
            ? ''
            : Buffer.from(this.rollingKey).toString('hex')
    }

    private get apiKeyAndSalt(): string {
        return (this.apiKey || '') + (this.salt || '');
    }

    private async fetchServerSeed(): Promise<string> {
        const r = await fetch('/api/v1/auth/seed')
        const seedRes = await r.json()
        return seedRes.seed
    }

    private async newRollingKey(): Promise<Uint8Array> {
        return new Uint8Array(
            await crypto.subtle.digest(
                'SHA-512',
                new TextEncoder().encode((this.serverSeed || '') + this.apiKeyAndSalt),
            ),
        )
    }

    private async refreshSeedAndRollingKey() {
        this.salt = randomString(32)
        this.serverSeed = await this.fetchServerSeed()
        this.rollingKey = await this.newRollingKey()
    }

    private async authorizationHeader(): Promise<string> {
        if (!this.serverSeed || this.rollingKey.length == 0)
            await this.refreshSeedAndRollingKey()

        const genNextKeyAppendValue = new TextEncoder().encode(this.apiKeyAndSalt)
        const inputForNextRollingKey = new Uint8Array(this.rollingKey.length + genNextKeyAppendValue.length)
        inputForNextRollingKey.set(this.rollingKey)
        inputForNextRollingKey.set(genNextKeyAppendValue, this.rollingKey.length)

        const nextRollingKey = await crypto.subtle.digest('SHA-512', inputForNextRollingKey)
        this.rollingKey = new Uint8Array(nextRollingKey)
        this.storeRollingKey()

        // FIXME btoa is deprecated for some reason
        const basicKeyValue = btoa(`sha512:${this.apiKeyId}:${this.salt}:${this.rollingKeyHex}`)
        return "Basic " + basicKeyValue;
    }
}

export const fetcher = new AuthenticatedFetcher()