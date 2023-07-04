import styled, {DefaultTheme, keyframes} from 'styled-components';

const getPulseAnimation = (theme: DefaultTheme) =>
  keyframes`
  0% {
    transform: scale(.7);
    box-shadow: 0 0 0 0 ${theme.color.primaryLight};
  }
  70% {
    transform: scale(1);
    box-shadow: 0 0 0 2px ${theme.color.primaryLight};
  }
  100% {
    transform: scale(.7);
    box-shadow: 0 0 0 0 ${theme.color.primaryLight};
  }
`;

export const PulseButton = styled.button`
  width: 9px;
  height: 9px;
  border: none;
  padding: 0px;
  border-radius: 50%;
  cursor: pointer;
  background: ${({theme}) => theme.color.primary};

  animation-name: ${({theme}) => getPulseAnimation(theme)};
  animation-duration: 1.5s;
  animation-iteration-count: infinite;
`;
