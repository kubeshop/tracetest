import {Typography} from 'antd';
import styled from 'styled-components';

export const Title = styled(Typography.Title).attrs({
  level: 3,
})`
  && {
    font-size: ${({theme}) => theme.size.md};
    font-weight: 600;
    margin-bottom: 16px;
  }
`;

export const TitleContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const Container = styled.div`
  margin: 16px 0;
`;
