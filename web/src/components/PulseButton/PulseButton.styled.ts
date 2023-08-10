import styled, {DefaultTheme, css, keyframes} from 'styled-components';

export const getPulseAnimation = (theme: DefaultTheme) =>
  keyframes`
  0% {
    transform: scale(.9);
    box-shadow: 0 0 0 0 ${theme.color.primaryLight};
  }
  70% {
    transform: scale(1);
    box-shadow: 0 0 0 2px ${theme.color.primaryLight};
  }
  100% {
    transform: scale(.9);
    box-shadow: 0 0 0 0 ${theme.color.primaryLight};
  }
`;

export const withPulseAnimation = (theme: DefaultTheme) => css`
  animation-name: ${getPulseAnimation(theme)};
  animation-duration: 1.5s;
  animation-iteration-count: infinite;
`;

export const PulseButton = styled.button`
  width: 9px;
  height: 9px;
  border: none;
  padding: 0px;
  border-radius: 50%;
  cursor: pointer;
  background: ${({theme}) => theme.color.primary};

  ${({theme}) => withPulseAnimation(theme)}
`;
