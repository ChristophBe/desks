// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Vuetify
import {createVuetify} from 'vuetify'


export default createVuetify(
    {
        theme: {
            defaultTheme: window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? "dark" : "light",
            themes: {
                dark: {
                    dark: true,
                    colors: {
                        primary: "#104f94",
                        secondary: "#009993",
                    }
                },
                light: {
                    dark: false,
                    colors: {
                        primary: "#104f94",
                        secondary: "#009993",
                    }
                },
            },
        }
    }
)
