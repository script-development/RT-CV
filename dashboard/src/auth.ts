function randomString(length: number): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let str = '';
    for (let i = 0; i < length; i++) {
        str += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return str;

};

class AuthenticatedFetcher {
    apiKey?: string;
    apiKeyId?: string;
    salt?: string;

    serverSeed?: string;

    rollingKey?: Uint8Array;

    async login(key: string, keyId: string) {
        this.apiKey = key
        this.apiKeyId = keyId
        this.salt = randomString(32)
    }

    get urlPrefix(): string {
        return location.port == '3000'
            ? `https://${location.hostname}:4000`
            : ''
    }

    get apiKeyAndSalt(): string {
        return (this.apiKey || '') + (this.salt || '');
    }

    async fetch(path: string) {
        const authHeader = await this.authorizationHeader()

        const r = await fetch(this.urlPrefix + (path[0] != '/' ? '/' : '') + path, {
            headers: {
                "Content-Type": "application/json",
                "Authorization": authHeader,
            },
        })
        return await r.json()
    }

    async fetchServerSeed() {
        const r = await fetch(this.urlPrefix + '/api/v1/auth/seed')
        const seedRes = await r.json()
        this.serverSeed = seedRes.seed
    }

    private async authorizationHeader(): Promise<string> {
        if (!this.serverSeed)
            // Get the server seed as we don't have it yet
            await this.fetchServerSeed()

        if (!this.rollingKey)
            // Gen a new rolling key
            this.rollingKey = new Uint8Array(
                await crypto.subtle.digest(
                    'SHA-512',
                    new TextEncoder().encode((this.serverSeed || '') + this.apiKeyAndSalt),
                ),
            )

        const genNextKeyAppendValue = new TextEncoder().encode(this.apiKeyAndSalt)
        const inputForNextRollingKey = new Uint8Array(this.rollingKey.length + genNextKeyAppendValue.length)
        inputForNextRollingKey.set(this.rollingKey)
        inputForNextRollingKey.set(genNextKeyAppendValue, this.rollingKey.length)

        const nextRollingKey = await crypto.subtle.digest( 'SHA-512', inputForNextRollingKey)
        this.rollingKey = new Uint8Array(nextRollingKey)


        const basicKeyValue = Buffer.from(`sha512:${this.apiKeyId}:${this.salt}:${Buffer.from(this.rollingKey).toString('hex')}`, 'base64').toString()
        return "Basic " + basicKeyValue;
    }
}

export const fetcher = new AuthenticatedFetcher()
