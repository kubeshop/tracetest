import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div``;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.md};
  }
`;

export const Description = styled(Typography.Paragraph)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
    font-size: ${({theme}) => theme.size.md};
  }
`;

export const SwitchContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
  margin-bottom: 18px;
`;
