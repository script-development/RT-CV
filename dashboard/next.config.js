module.exports = {
    webpack(config, _) {
        // Required for monaco editor to work
        if (config.resolve.fallback) {
            config.resolve.fallback.fs = false
        }
        return config
    },
    swcMinify: true, // Much faster way of minifying the javascript when running docker build
}
