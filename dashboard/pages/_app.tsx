import type { AppProps } from 'next/app'
import Head from 'next/head'
import CssBaseline from '@material-ui/core/CssBaseline'
import { ThemeProvider } from '@material-ui/core/styles'
import { useEffect } from 'react'
import { fetcher } from '../src/auth'
import { useRouter } from 'next/router'
import { theme } from '../src/theme'

function MyApp({ Component, pageProps }: AppProps) {
  const router = useRouter()

  useEffect(() => {
    if ((!fetcher.getApiKey || !fetcher.getApiKeyId) && router.route != '/login') {
      router.push('/login')
    }
  }, [])

  return (<>
    <Head>
      <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700&display=swap" />
      <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons" />
      <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
      <link rel="icon" href="/favicon.ico" />
    </Head>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Component {...pageProps} />
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
    `}</style>
  </>)
}
export default MyApp
