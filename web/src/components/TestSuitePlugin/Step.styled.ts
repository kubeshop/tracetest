import {Typography} from 'antd';
import styled from 'styled-components';

export const Step = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 426px;
  max-height: 426px;
  overflow-y: scroll;
  overflow-x: hidden;
`;

export const FormContainer = styled.div`
  height: 100%;
  display: flex;
  flex-direction: column;
  margin-bottom: 24px;

  .ant-form-item {
    margin: 0;
  }
`;

export const Title = styled(Typography.Title)<{$withSubtitle?: boolean}>`
  && {
    font-size: ${({theme}) => theme.size.md};
    margin-bottom: ${({$withSubtitle}) => ($withSubtitle ? '0' : '16px')};
  }
`;

export const Subtitle = styled(Typography.Text)`
  && {
    font-size: ${({theme}) => theme.size.md};
    color: ${({theme}) => theme.color.textSecondary};
    margin-bottom: 16px;
  }
`;
