import {Typography} from 'antd';
import styled, {css} from 'styled-components';

export const Ticks = styled.div`
  padding-right: 16px;
  pointer-events: none;
  position: relative;
`;

export const Tick = styled.div`
  align-items: center;
  background: ${({theme}) => theme.color.borderLight};
  display: flex;
  height: 100%;
  position: absolute;
  width: 1px;

  :first-child,
  :last-child {
    width: 0;
  }
`;

export const TickLabel = styled(Typography.Text)<{$isEndAnchor: boolean}>`
  color: ${({theme}) => theme.color.text};
  font-size: ${({theme}) => theme.size.sm};
  font-weight: 400;
  left: 0.25rem;
  position: absolute;
  white-space: nowrap;

  ${({$isEndAnchor}) =>
    $isEndAnchor &&
    css`
      left: initial;
      right: 0.25rem;
    `};
`;
