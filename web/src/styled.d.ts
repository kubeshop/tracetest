// import original module declarations
import 'styled-components';

// Extend them
declare module 'styled-components' {
  export interface DefaultTheme {
    /** Colors */
    color: {
      background: string;
      backgroundInteractive: string;
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
      warningYellow: string;
    };
    /** Font size */
    size: {
      xs: string;
      sm: string;
      md: string;
      lg: string;
      xl: string;
    };
    /** Font defaults */
    font: {
      family: string;
    };

    notification: {
      success: {
        style: {
          border: string;
          background: string;
        };
      };
      error: {
        style: {
          border: string;
          background: string;
        };
      };
    };
  }
}
