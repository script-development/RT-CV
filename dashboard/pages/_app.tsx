import type { AppProps } from 'next/app'
import Head from 'next/head'
import CssBaseline from '@material-ui/core/CssBaseline'
import { ThemeProvider } from '@material-ui/core/styles'
import { useEffect, useState } from 'react'
import { fetcher } from '../src/auth'
import { useRouter } from 'next/router'
import { theme } from '../src/theme'

function MyApp(args: AppProps) {
  return (<>
    <Head>
      <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700&display=swap" />
      <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons" />
      <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
      <link rel="icon" href="/favicon.ico" />
    </Head>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AppContent {...args} />
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
export default MyApp

function AppContent({ Component, pageProps }: AppProps) {
  const router = useRouter()
  const [version, setVersion] = useState({
    version: '',
    githubCommitURL: '',
  })

  useEffect(() => {
    if ((!fetcher.getApiKey || !fetcher.getApiKeyId) && router.route != '/login') {
      router.push('/login')
      return
    }

    let mounted = true
    fetch('/api/v1/health', {
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
