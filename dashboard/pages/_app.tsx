import type { AppProps } from 'next/app'
import Head from 'next/head'
import CssBaseline from '@material-ui/core/CssBaseline'
import { createTheme, ThemeProvider } from '@material-ui/core/styles'
import { red } from '@material-ui/core/colors'
import { useEffect } from 'react'
import { fetcher } from '../src/auth'
import { useRouter } from 'next/router'

const theme = createTheme({
  palette: {
    type: 'dark',
    // https://material-ui.com/customization/color/#playground
    primary: { main: '#ff9800' }, // orange
    secondary: { main: '#ff6e40' }, // red
    error: {
      main: red.A400,
    },
  },
});

function MyApp({ Component, pageProps }: AppProps) {
  const router = useRouter()

  useEffect(() => {
    if (!fetcher.tryRestoreCredentials() && router.route != '/login') {
      router.push('/login')
    }
  }, [])

  return (<>
    <Head>
      <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700&display=swap" />
      <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons" />
      <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
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
