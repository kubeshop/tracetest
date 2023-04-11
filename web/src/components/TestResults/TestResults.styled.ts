import {Button, Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  height: 100%;
  position: relative;
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  justify-content: space-between;
  margin-bottom: 38px;
`;

export const HeaderDetail = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.md};
  margin-right: 8px;
`;

export const HeaderDot = styled.span<{$passed: boolean}>`
  background-color: ${({$passed, theme}) => ($passed ? theme.color.success : theme.color.error)};
  border-radius: 50%;
  display: inline-block;
  height: 10px;
  line-height: 0;
  margin-right: 4px;
  vertical-align: -0.1em;
  width: 10px;
`;

export const HeaderText = styled(Typography.Title).attrs({level: 2})`
  && {
    margin-bottom: 0;
  }
`;

export const LoadingContainer = styled.div`
  text-align: center;
`;

export const Row = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
`;

export const CaretDropdownButton = styled(Button)`
  font-weight: 600;
  opacity: 0.7;
  width: 32px;
  padding: 0px;
`;
