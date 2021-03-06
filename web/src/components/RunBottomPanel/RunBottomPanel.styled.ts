import {Button, Typography} from 'antd';
import styled from 'styled-components';
import chevron from '../../assets/chevron.svg';

export const Header = styled.div`
  display: flex;
  align-items: center;
  cursor: pointer;
  justify-content: space-between;
  width: 100%;
  background: #f5f5fa;
  height: 64px;
  padding: 0 24px;
`;

export const HeaderText = styled(Typography.Text)``;

export const StartDateText = styled(Typography.Text)`
  && {
    font-size: ${({theme}) => theme.size.sm};
    margin-left: 14px;
    margin-right: 40px;
  }
`;

export const CountNumber = styled.span`
  margin-right: 15px;
`;

export const Container = styled.div`
  background-color: ${({theme}) => theme.color.white};
  height: calc(100% - 64px);
  overflow-y: scroll;
`;

export const Content = styled.div`
  padding: 24px;
`;

export const AddAssertionButton = styled(Button).attrs({
  type: 'primary',
})`
  && {
    font-weight: 600;
    margin-left: 14px;
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

export const ChevronContainer = styled.span`
  margin-left: 16px;
`;

export const Row = styled.div`
  display: flex;
  gap: 8px;
  align-items: center;
`;
