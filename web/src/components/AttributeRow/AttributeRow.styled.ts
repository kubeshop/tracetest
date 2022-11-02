import {Badge, Tag as AntdTag, Typography} from 'antd';
import styled from 'styled-components';

export {default as AttributeTitle} from './AttributeTitle';

export const Container = styled.div`
  align-items: center;
  background-color: ${({theme}) => theme.color.white};
  display: flex;
  margin-bottom: 4px;
  padding: 12px;
  transition: background-color 0.2s ease;

  &:hover {
    background-color: ${({theme}) => theme.color.background};
  }
`;

export const Header = styled.div`
  cursor: pointer;
  flex: 1;
`;

export const AttributeValueRow = styled.div`
  display: flex;
  word-break: break-word;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const Text = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.sm};
`;

export const CustomBadge = styled(Badge)`
  border: ${({theme}) => `1px solid ${theme.color.textSecondary}`};
  border-radius: 9999px;
  cursor: pointer;
  line-height: 19px;
  margin-left: 8px;
  padding: 0 8px;
  white-space: nowrap;

  .ant-badge-status-text {
    color: ${({theme}) => theme.color.textSecondary};
    font-size: ${({theme}) => theme.size.sm};
    margin-left: 3px;
  }
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 0;
  }
`;

export const DetailContainer = styled.div`
  width: 270px;
`;

export const TagsContainer = styled.div`
  margin-top: 8px;
`;

export const Tag = styled(AntdTag)`
  background: #e7e8eb;
`;
