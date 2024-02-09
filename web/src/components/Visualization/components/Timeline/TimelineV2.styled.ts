import {Typography} from 'antd';
import {SemanticGroupNames, SemanticGroupNamesToColor} from 'constants/SemanticGroupNames.constants';
import styled, {css} from 'styled-components';

export const Container = styled.div`
  padding: 50px 24px 0 24px;
`;

export const Row = styled.div<{$isEven: boolean; $isMatched: boolean; $isSelected: boolean}>`
  background-color: ${({theme, $isEven}) => ($isEven ? theme.color.background : theme.color.white)};
  display: grid;
  grid-template-columns: 300px 1fr;
  grid-template-rows: 32px;
  padding: 0px 16px;

  :hover {
    background-color: ${({theme}) => theme.color.backgroundInteractive};
  }

  ${({$isMatched}) =>
    $isMatched &&
    css`
      background-color: ${({theme}) => theme.color.alertYellow};
    `};

  ${({$isSelected}) =>
    $isSelected &&
    css`
      background: rgba(97, 23, 94, 0.1);

      :hover {
        background: rgba(97, 23, 94, 0.1);
      }
    `};
`;

export const Col = styled.div`
  display: grid;
  grid-template-columns: 1fr 8px;
`;

export const ColDuration = styled.div`
  overflow: hidden;
  position: relative;
`;

export const Header = styled.div`
  align-items: center;
  display: flex;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
`;

export const NameContainer = styled.div`
  overflow: hidden;
  text-overflow: ellipsis;
`;

export const Separator = styled.div`
  border-left: 1px solid rgb(222, 227, 236);
  cursor: ew-resize;
  height: 32px;
  padding: 0px 3px;
  width: 1px;
`;

export const Title = styled(Typography.Text)`
  color: ${({theme}) => theme.color.text};
  font-size: ${({theme}) => theme.size.sm};
  font-weight: 400;
`;

export const Connector = styled.svg`
  flex-shrink: 0;
  overflow: hidden;
  overflow-clip-margin: content-box;
`;

export const SpanBar = styled.div<{$type: SemanticGroupNames}>`
  background-color: ${({$type}) => SemanticGroupNamesToColor[$type]};
  border-radius: 3px;
  height: 18px;
  min-width: 2px;
  position: absolute;
  top: 7px;
`;

export const SpanBarLabel = styled.div<{$side: 'left' | 'right'}>`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.xs};
  padding: 1px 4px 0 4px;
  position: absolute;

  ${({$side}) =>
    $side === 'left'
      ? css`
          right: 100%;
        `
      : css`
          left: 100%;
        `};
`;

export const TextConnector = styled.text<{$isActive?: boolean}>`
  fill: ${({theme, $isActive}) => ($isActive ? theme.color.white : theme.color.text)};
  font-size: ${({theme}) => theme.size.xs};
`;

export const CircleDot = styled.circle`
  fill: ${({theme}) => theme.color.textSecondary};
  stroke-width: 2;
  stroke: ${({theme}) => theme.color.white};
`;

export const LineBase = styled.line`
  stroke: ${({theme}) => theme.color.textSecondary};
`;

export const RectBase = styled.rect<{$isActive?: boolean}>`
  fill: ${({theme, $isActive}) => ($isActive ? theme.color.primary : theme.color.white)};
  stroke: ${({theme}) => theme.color.textSecondary};
`;

export const RectBaseTransparent = styled(RectBase)`
  cursor: pointer;
  fill: transparent;
`;

export const HeaderRow = styled.div`
  background-color: ${({theme}) => theme.color.white};
  display: grid;
  grid-template-columns: 300px 1fr;
  grid-template-rows: 32px;
  padding: 0px 16px;
`;

export const HeaderContent = styled.div`
  align-items: center;
  display: flex;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
`;

export const HeaderTitle = styled(Typography.Title)`
  && {
    margin: 0;
  }
`;
