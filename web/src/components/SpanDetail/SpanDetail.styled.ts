import {Divider, Typography} from 'antd';
import styled from 'styled-components';

export const Header = styled.div`
  padding: 16px 12px 0;
  min-width: 270px;
`;

export const HeaderCheck = styled(Typography.Text)`
  color: ${({theme}) => theme.color.text};
  font-size: ${({theme}) => theme.size.sm};
  margin-right: 8px;
`;

export const HeaderDot = styled.span<{$passed: boolean}>`
  background-color: ${({$passed, theme}) => ($passed ? theme.color.success : theme.color.error)};
  height: 10px;
  width: 10px;
  display: inline-block;
  margin-right: 4px;
  line-height: 0;
  vertical-align: -0.1em;
  border-radius: 50%;
`;

export const HeaderDivider = styled(Divider)`
  margin: 16px 0;
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
    margin: 0;
  }
`;

export const Column = styled.div`
  align-items: flex-start;
  display: flex;
  flex-direction: column;
  margin-bottom: 4px;
`;

export const Row = styled.div`
  align-items: center;
  display: flex;
  margin-bottom: 4px;
`;

export const AttributesContainer = styled.div<{$top: number}>`
  height: calc(100% - ${({$top}) => `${$top}px`});
  overflow-y: scroll;
`;

export const SearchContainer = styled.div`
  padding: 0 12px;
`;
