export function randomString(length: number): string {
    const cryptoRandomValues = new Uint32Array(length);
    crypto.getRandomValues(cryptoRandomValues);

    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let str = '';
    for (let i = 0; i < length; i++) {
        str += chars.charAt(Math.floor(cryptoRandomValues[i] / 4294967295 * chars.length))
    }
    return str;
}
