import styled, {css} from 'styled-components';
import {Modal as AntModal, Typography} from 'antd';

export const CardContainer = styled.div<{$isActive: boolean; $isSelected: boolean}>`
  align-items: center;
  background: ${({theme}) => theme.color.white};
  border: 1px solid ${({theme}) => theme.color.borderLight};
  border-radius: 4px;
  cursor: ${({$isActive}) => ($isActive ? 'pointer' : 'default')};
  display: flex;
  gap: 12px;
  padding: 4px;
  padding-left: 16px;
  transition: 0.3s;
  width: 48%;

  &:hover {
    background: ${({theme}) => theme.color.background};
    border: 1px solid ${({theme}) => theme.color.primary};

    .check {
      opacity: 1;
    }
  }

  ${({$isSelected}) =>
    $isSelected &&
    css`
      background: ${({theme}) => theme.color.background};
      border: 1px solid ${({theme}) => theme.color.primary};

      .check {
        opacity: 1;
      }
    `}
`;

export const IntegrationCardContainer = styled(CardContainer)`
  width: auto;
  padding-right: 16px;
`;

export const CardContent = styled.div`
  display: flex;
  flex-direction: column;
`;

export const CardDescription = styled(Typography.Text)<{$isActive: boolean}>`
  font-size: ${({theme}) => theme.size.xs};
  opacity: ${({$isActive}) => ($isActive ? 1 : 0.5)};
`;

export const CardList = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 24px;
`;

export const IntegrationCardList = styled(CardList)`
  margin-bottom: 16px;
`;

export const CardTitle = styled(Typography.Text).attrs({
  strong: true,
})<{$isActive: boolean}>`
  display: inline-block;
  font-size: ${({theme}) => theme.size.sm};
  opacity: ${({$isActive}) => ($isActive ? 1 : 0.5)};

  a {
    opacity: 1;
  }
`;

export const Check = styled.div`
  background: ${({theme}) => theme.color.primary};
  border-radius: 50%;
  display: inline-block;
  height: 8px;
  opacity: 0;
  width: 8px;
`;

export const Circle = styled.div<{$isActive: boolean}>`
  align-items: center;
  border-radius: 50%;
  border: ${({theme}) => `1px solid ${theme.color.primary}`};
  display: flex;
  justify-content: center;
  max-height: 16px;
  max-width: 16px;
  min-height: 16px;
  min-width: 16px;
  opacity: ${({$isActive}) => ($isActive ? 1 : 0.5)};
`;

export const Modal = styled(AntModal)`
  .ant-modal-body {
    background: ${({theme}) => theme.color.background};
  }
`;

export const Title = styled(Typography.Title)<{$marginBottom?: number}>`
  && {
    margin-bottom: ${({$marginBottom}) => $marginBottom || 0}px;
  }
`;

export const Text = styled(Typography.Text)``;

export const Divider = styled.div`
  height: 1px;
  border-top: 1px dashed ${({theme}) => theme.color.borderLight};
  margin-bottom: 24px;
`;
