import styled from 'styled-components';

export const Container = styled.div`
  position: absolute;
  left: 20%;
  display: flex;
  justify-content: center;
  margin-top: 2px;
`;

export const FloatingText = styled.div`
  background-color: ${({theme}) => theme.color.interactive};
  color: ${({theme}) => theme.color.white};
  border-radius: 12px;
  padding: 2px 6px;
  font-size: ${({theme}) => theme.size.xs};
  cursor: pointer;
`;
