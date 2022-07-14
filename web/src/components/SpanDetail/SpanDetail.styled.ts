import {Typography} from 'antd';
import styled from 'styled-components';

export const Header = styled.div`
  margin-bottom: 8px;
`;

export const HeaderItem = styled.div`
  align-items: center;
  color: ${({theme}) => theme.color.text};
  display: flex;
  font-size: ${({theme}) => theme.size.md};
  margin-right: 8px;
`;

export const HeaderItemText = styled(Typography.Text)`
  color: inherit;
  margin-left: 5px;
`;

export const HeaderTitle = styled(Typography.Title)`
  && {
    margin: 0 0 0 8px;
  }
`;

export const Row = styled.div`
  align-items: center;
  display: flex;
  margin-bottom: 4px;
`;

export const SpanDetail = styled.div`
  display: flex;
  flex-direction: column;
  padding: 24px;
`;

export const Dot = styled.div`
  background-color: ${({theme}) => theme.color.textHighlight};
  border-radius: 50%;
  display: inline-block;
  height: 10px;
  margin-left: 5px;
  width: 10px;
`;
