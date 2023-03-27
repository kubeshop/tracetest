// import original module declarations
import {CSSProperties} from 'react';
import 'styled-components';

type TNotification = {
  color: string;
  style: CSSProperties;
};

// Extend them
declare module 'styled-components' {
  export interface DefaultTheme {
    /** Colors */
    color: {
      background: string;
      backgroundInteractive: string;
      backgroundDark: string;
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
      alertYellow: string;
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
      success: TNotification;
      error: TNotification;
      info: TNotification;
      warning: TNotification;
      open: TNotification;
    };
  }
}
