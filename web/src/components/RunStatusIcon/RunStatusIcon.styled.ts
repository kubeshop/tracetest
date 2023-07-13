import {CheckCircleFilled, InfoCircleFilled, LoadingOutlined, MinusCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const IconSuccess = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
`;

export const IconFail = styled(MinusCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const IconInfo = styled(InfoCircleFilled)`
  color: ${({theme}) => theme.color.textLight};
`;

export const SequenceText = styled.div`
  color: ${({theme}) => theme.color.textLight};
  font-weight: 600;
`;

export const IconWrapper = styled.div`
  width: 14px;
  cursor: pointer;
`;

export const LoadingIcon = styled(LoadingOutlined)`
  color: ${({theme}) => theme.color.primary};
`;

export const GateList = styled.ul<{$color: string}>`
  && {
    padding-left: 25px;
    color: ${({$color}) => $color};

    span {
      color: ${({$color}) => $color};
    }
  }
`;

export const Gate = styled.li``;

export const GateName = styled(Typography.Text)``;
