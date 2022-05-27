import {Button, Typography} from 'antd';
import styled, {css} from 'styled-components';
import chevron from '../../assets/chevron.svg';

export const Header = styled.div<{visiblePortion: number}>`
  display: flex;
  align-items: center;
  cursor: pointer;
  justify-content: space-between;
  width: 100%;
  background: #f5f5fa;
  ${props =>
    css`
      height: ${props.visiblePortion}px;
    `}
  padding: 0 24px;
  color: rgb(213, 215, 224);
`;

export const HeaderText = styled(Typography.Text)``;

export const StartDateText = styled(Typography.Text)`
  && {
    margin-left: 14px;
    margin-right: 40px;
    font-size: 12px;
  }
`;

export const CountNumber = styled.span`
  margin-right: 15px;
`;

export const Container = styled.div`
  border: 1px solid rgba(0, 0, 0, 0.06);
  margin-bottom: 16px;
  overflow-y: hidden;
`;

export const Content = styled.div`
  padding: 24px;
  overflow-y: scroll;
  height: 330px;
  position: relative;
`;

export const AddAssertionButton = styled(Button).attrs({
  type: 'primary',
})`
  && {
    font-weight: 600;
  }
`;

export const LoadingSpinnerContainer = styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const Chevron = styled.img.attrs({
  src: chevron,
})<{$isCollapsed: boolean}>`
  transform: ${({$isCollapsed}) => ($isCollapsed ? 'rotate(0deg)' : 'rotate(180deg)')};
`;
