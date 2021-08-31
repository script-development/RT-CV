import { createTheme } from '@material-ui/core/styles'
import { red } from '@material-ui/core/colors'

export const primaryColor = '#ff9800'
export const secondaryColor = '#ff6e40'

export const theme = createTheme({
    palette: {
        type: 'dark',
        // https://material-ui.com/customization/color/#playground
        primary: { main: primaryColor }, // orange
        secondary: { main: secondaryColor }, // red
        error: {
            main: red.A400,
        },
    },
});
