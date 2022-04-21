import styled, {keyframes} from 'styled-components';

export const Header = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  width: 100%;
  height: 56px;
  padding: 0 32px;
  border-bottom: 1px solid rgb(213, 215, 224);
`;

export const TestDetailsHeader = styled.div`
  display: flex;
  width: 100%;
  justify-content: space-between;
  margin: 32px 0px 24px;
`;

export const Wrapper = styled.div`
  padding: 0px 24px;
`;

const IndeterminateAnimation = keyframes`
 0% {
    transform:  translateX(0) scaleX(0);
  }
  40% {
    transform:  translateX(0) scaleX(0.4);
  }
  100% {
    transform:  translateX(100%) scaleX(0.5);
  }
`;

export const ProgressBarContainer = styled.div`
  height: 4px;
  background-color: rgba(5, 114, 206, 0.2);
  width: 100%;
  overflow: hidden;
`;

export const Progress = styled.div`
  width: 100%;
  height: 100%;
  background-color: rgb(5, 114, 206);
  animation: ${IndeterminateAnimation} 1s infinite linear;
  transform-origin: 0% 50%;
`;
