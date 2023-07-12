import styled from 'styled-components';

export const SelectAsCurrentContainer = styled.div`
  display: flex;
  justify-content: center;
  left: 20%;
  margin-top: 2px;
  position: absolute;
`;

export const FloatingText = styled.div`
  background-color: ${({theme}) => theme.color.interactive};
  border-radius: 12px;
  color: ${({theme}) => theme.color.white};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.xs};
  padding: 2px 6px;
`;
