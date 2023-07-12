import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  background-color: ${({theme}) => theme.color.white};
  display: flex;
  gap: 8px;
  padding: 0 12px;
`;
