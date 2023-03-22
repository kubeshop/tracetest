import {DeleteOutlined} from '@ant-design/icons';
import {Button, Typography} from 'antd';
import styled from 'styled-components';
import RequestDetailsHeadersInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsHeadersInput';

export const FormContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  width: 70%;
`;

export const AddButton = styled(Button).attrs({
  type: 'link',
})`
  padding: 0;
  font-weight: 600;
`;

export const DeleteCheckIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  margin-top: 9px;
`;

export const ListItem = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

export const ItemListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

export const ItemActionContainer = styled.div`
  display: flex;
  justify-content: center;
  align-self: flex-start;
  flex-basis: 5%;
`;

export const ItemListLabel = styled(Typography.Text)`
  padding: 0 0 8px;
  display: block;
`;

export const FormColumn = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.md};
    font-weight: 700;
    margin: 0;
  }
`;

export const HeadersInput = styled(RequestDetailsHeadersInput)`
  max-width: 1000px;
`;
