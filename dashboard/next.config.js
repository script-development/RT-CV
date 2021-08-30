module.exports = {
    async rewrites() {
        return [
            {
                source: '/api/:path*',
                destination: `http://localhost:4000/api/:path*`,
            },
        ]
    },
    webpack(config, _) {
        if (config.resolve.fallback) {
            config.resolve.fallback.fs = false
        }
        return config
    },
}
