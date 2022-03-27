import { createTheme } from '@mui/material/styles'
import { red, orange, deepOrange } from '@mui/material/colors';
import createCache from '@emotion/cache';

export const theme = createTheme({
    palette: {
        mode: 'dark',
        // https://material-ui.com/customization/color/#playground
        primary: orange, // orange
        secondary: { main: deepOrange.A200 }, // red
        error: {
            main: red.A400,
        },
    },
});

export const createEmotionCache = () => createCache({ key: 'css' });
export const primaryColor = orange[500];
export const secondaryColor = deepOrange.A200;
