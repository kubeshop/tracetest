import {ExclamationCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled, {css} from 'styled-components';

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
  border-radius: 10px;
  font-size: ${({theme}) => theme.size.xs};
  height: 100px;
  padding: 10px;
  position: absolute;
  right: -210px;
  top: -32px;
  width: 200px;
`;

export const Content = styled.div`
  height: 100%;
  overflow-y: scroll;

  ::-webkit-scrollbar {
    display: none;
  }
`;

export const ErrorIcon = styled(ExclamationCircleFilled)<{$isAbsolute?: boolean}>`
  color: ${({theme}) => theme.color.error};

  ${({$isAbsolute}) =>
    $isAbsolute &&
    css`
      position: absolute;
      right: 10px;
      top: 5px;
    `}
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
