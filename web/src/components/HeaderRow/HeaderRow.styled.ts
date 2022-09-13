import styled from 'styled-components';

export const HeaderContainer = styled.div`
  align-items: center;
  background: ${({theme}) => theme.color.background};
  border: 1px solid ${({theme}) => theme.color.borderLight};
  display: flex;
  padding: 7px 16px;
  transition: background-color 0.2s ease;
`;

export const Header = styled.div`
  flex: 1;
`;

export const HeaderValue = styled.div`
  display: flex;
  word-break: break-word;
  font-size: ${({theme}) => theme.size.sm};
  font-weight: 600;
`;

export const HeaderKey = styled.div`
  display: flex;
  word-break: break-word;
  color: ${({theme}) => theme.color.textLight};
  font-size: ${({theme}) => theme.size.sm};
  font-weight: 400;
`;
