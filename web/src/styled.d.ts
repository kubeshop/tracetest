// import original module declarations
import 'styled-components';

// Extend them
declare module 'styled-components' {
  export interface DefaultTheme {
    /** Colors */
    color: {
      background: string;
      border: string;
      borderLight: string;
      error: string;
      interactive: string;
      primary: string;
      success: string;
      text: string;
      textHighlight: string;
      textLight: string;
      textSecondary: string;
      white: string;
    };
    /** Font size */
    size: {
      xs: string;
      sm: string;
      md: string;
      lg: string;
      xl: string;
    };
  }
}
