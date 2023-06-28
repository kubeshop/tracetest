import {ExclamationCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Body = styled.div`
  margin-top: 5px;
  position: relative;
`;

export const Connector = styled.div`
  background-color: ${({theme}) => theme.color.error};
  height: 2px;
  left: -26px;
  position: absolute;
  top: 49px;
  width: 26px;
  z-index: 9;
`;

export const Container = styled.div`
  border: ${({theme}) => `2px solid ${theme.color.error}`};
  background-color: ${({theme}) => theme.color.white};
  border-radius: 10px;
  bottom: -40px;
  font-size: ${({theme}) => theme.size.xs};
  height: 100px;
  padding: 10px;
  position: absolute;
  right: -220px;
  width: 200px;
`;

export const Content = styled.div`
  height: 100%;
  overflow-y: scroll;

  ::-webkit-scrollbar {
    display: none;
  }
`;

export const ErrorIcon = styled(ExclamationCircleFilled)<{$isClickable?: boolean}>`
  color: ${({theme}) => theme.color.error};
  cursor: ${({$isClickable}) => $isClickable && 'pointer'};
`;

export const List = styled.ul`
  padding-inline-start: 20px;
`;

export const RuleContainer = styled.div`
  margin-bottom: 8px;
`;

export const Text = styled(Typography.Text)`
  color: inherit;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 0;
  }
`;
