import {InfoCircleOutlined} from '@ant-design/icons';
import {Tag as AntdTag, Typography} from 'antd';
import styled from 'styled-components';
import TestOutputMark from 'components/TestOutputMark';
import moreIcon from 'assets/more.svg';

export {default as AttributeTitle} from './AttributeTitle';

export const Container = styled.div`
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

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 0;
  }
`;

export const DetailContainer = styled.div`
  overflow-wrap: break-word;
  width: 270px;
`;

export const TagsContainer = styled.div`
  margin-top: 8px;
`;

export const Tag = styled(AntdTag)`
  background: #e7e8eb;
  margin-bottom: 8px;
`;

export const SectionTitle = styled.div`
  align-items: center;
  display: flex;
`;

export const InfoIcon = styled(InfoCircleOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  margin: 4px;
`;

export const OutputsMark = styled(TestOutputMark)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
    margin: 4px;
  }
`;

export const MoreIcon = styled.img.attrs({
  src: moreIcon,
})``;
