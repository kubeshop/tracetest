import styled from 'styled-components';

export const Container = styled.div`
  background-color: ${({theme}) => theme.color.white};
  padding: 12px;
  transition: background-color 0.2s ease;

  &:hover {
    background-color: ${({theme}) => theme.color.background};
  }
`;

export const Footer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
`;
