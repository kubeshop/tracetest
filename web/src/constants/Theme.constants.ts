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
  },
  font: {
    family: 'SFPro',
  },
  notification: {
    success: {
      style: {
        border: '1px solid #52C41A',
        background: '#F6FFED',
      },
    },
    error: {
      style: {
        border: '1px solid #F5222D',
        background: '#FFF1F0',
      },
    },
  },
};
