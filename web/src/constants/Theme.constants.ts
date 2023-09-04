import {DefaultTheme} from 'styled-components';

export const theme: DefaultTheme = {
  color: {
    background: '#FBFBFF',
    backgroundInteractive: 'rgba(56, 101, 246, 0.05)',
    backgroundDark: '#E2E4ED',
    border: '#CDD1DB',
    borderLight: 'rgba(3, 24, 73, 0.1)',
    error: '#FF4D4F',
    interactive: '#2D62FF',
    primary: '#61175E',
    primaryLight: '#61175e33',
    success: '#52C41A',
    text: '#031849',
    textHighlight: '#61175e29',
    textLight: 'rgba(3, 24, 73, 0.4)',
    textSecondary: '#687492',
    white: '#FFFFFF',
    warningYellow: '#FAAD14',
    alertYellow: '#fffbe6',
  },
  size: {
    xs: '10px',
    sm: '12px',
    md: '14px',
    lg: '16px',
    xl: '18px',
    xxxl: '24px',
  },
  font: {
    family: 'SFPro, Inter',
  },
  notification: {
    success: {
      color: '#52C41A',
      style: {
        border: '1px solid #52C41A',
        background: '#F6FFED',
        minWidth: '450px',
      },
    },
    error: {
      color: '#EE1847',
      style: {
        border: '1px solid #EE1847',
        background: '#FFF1F0',
        minWidth: '450px',
      },
    },
    info: {
      color: '#3B61F6',
      style: {
        border: '1px solid #3B61F6',
        background: '#E7EBFE',
        minWidth: '450px',
      },
    },
    warning: {
      color: '#FAAD14',
      style: {
        border: '1px solid #FAAD14',
        background: '#FFFBE6',
        minWidth: '450px',
      },
    },
    open: {
      color: '#3B61F6',
      style: {
        minWidth: '450px',
      },
    },
  },
};
