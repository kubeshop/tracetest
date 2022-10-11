import styled from 'styled-components';

export const ActionsContainer = styled.div`
  display: flex;
  justify-content: space-between;
  margin: 32px 0 24px;
  width: 100%;
`;

export const Container = styled.div<{$isWhite?: boolean}>`
  background: ${({$isWhite, theme}) => $isWhite && theme.color.white};
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  padding: 0 24px;
`;
