import type { AppProps } from 'next/app'
import Head from 'next/head'
import CssBaseline from '@material-ui/core/CssBaseline'
import { ThemeProvider } from '@material-ui/core/styles'
import { useEffect, useState } from 'react'
import { fetcher } from '../src/auth'
import { useRouter } from 'next/router'
import { theme } from '../src/theme'

export default function MyApp(props: AppProps) {
  useEffect(() => {
    // Remove the server-side injected CSS.
    const jssStyles = document.querySelector('#jss-server-side');
    if (jssStyles) {
      jssStyles.parentElement?.removeChild(jssStyles);
    }
  }, []);

  return (<>
    <Head>
      <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
      <meta name="theme-color" content={theme.palette.primary.main} />
    </Head>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AppContent {...props} />
    </ThemeProvider>
    <style jsx global>{`
      * {
        padding: 0;
        margin: 0;
      }
      body {
        font-size: 17px;
      }
      h1 {
        margin-bottom: 10px;
      }
      h2 {
        margin-bottom: 6px;
      }
      h3 {
        margin-bottom: 4px;
      }
      a {
        color: white;
      }
    `}</style>
  </>)
}

function AppContent({ Component, pageProps }: AppProps) {
  const router = useRouter()
  const [version, setVersion] = useState({
    version: '',
    githubCommitURL: '',
  })

  useEffect(() => {
    if ((!fetcher.getApiKeyHashed || !fetcher.getApiKeyId) && router.route != '/login') {
      let url = new URL('http://localhost')
      url.searchParams.set('redirectTo', location.pathname + location.search + location.hash)

      location.href = "/login" + url.search
      return
    }

    let mounted = true
    fetch(fetcher.getAPIPath('/api/v1/health'), {
      headers: { 'Content-Type': 'application/json' }
    })
      .then(r => r.json())
      .then(({ appVersion }) => {
        if (mounted) {
          if (appVersion.length == 40) {
            // App version is a git commit hash
            setVersion({ version: appVersion, githubCommitURL: 'https://github.com/script-development/RT-CV/commit/' + appVersion })
          } else {
            setVersion({ version: appVersion, githubCommitURL: '' })
          }
        }
      })
    return () => { mounted = false }

  }, [])

  return (
    <div className="appContainer">
      <Component {...pageProps} />
      <div className="version">
        version: <b>{version.githubCommitURL ? <a href={version.githubCommitURL}>{version.version}</a> : version.version}</b>
      </div>
      <style jsx>{`
        .appContainer {
          min-height: 100vh;
          display: flex;
          flex-direction: column;
        }
        .version {
          padding: 10px;
          text-align: center;
          color: rgba(255, 255, 255, 0.7);
        }
      `}</style>
    </div>
  )
}
