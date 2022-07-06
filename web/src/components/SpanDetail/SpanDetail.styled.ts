import {Typography} from 'antd';
import styled from 'styled-components';

export const SpanHeader = styled.div`
  width: 100%;
  align-items: center;
  margin-bottom: 24px;
`;

export const SpanHeaderTitle = styled(Typography.Title).attrs({level: 2})`
  && {
    margin-bottom: 5px;
  }
`;

export const SpanDetail = styled.div`
  padding: 24px;
  display: flex;
  flex-direction: column;
`;

export const Dot = styled.div`
  height: 10px;
  width: 10px;
  margin-left: 5px;
  background-color: ${({theme}) => theme.color.textHighlight};
  border-radius: 50%;
  display: inline-block;
`;
