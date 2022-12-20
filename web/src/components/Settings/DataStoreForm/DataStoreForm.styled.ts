import styled from 'styled-components';
import {SupportedDataStores} from 'types/Config.types';
import Jaeger from 'assets/jaeger.svg';
import OpenSearch from 'assets/openSearch.svg';
import SignalFx from 'assets/signalFx.svg';
import Tempo from 'assets/tempo.svg';
import OtelCollector from 'assets/otlp.svg';
import {Typography} from 'antd';

export const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 24px;

  .ant-form-item {
    margin: 0;
  }
`;

export const DataStoreListContainer = styled.div`
  display: flex;
  gap: 16px;
  align-items: center;
`;

export const DataStoreItemContainer = styled.div<{$isSelected: boolean}>`
  display: flex;
  align-items: center;
  gap: 12px;

  background: ${({$isSelected, theme}) => ($isSelected ? theme.color.background : theme.color.white)};
  border: ${({$isSelected, theme}) => `1px solid ${$isSelected ? theme.color.primary : theme.color.textSecondary}`};
  border-radius: 4px;
  padding: 12px 22px;
  cursor: pointer;
`;

export const Circle = styled.div`
  min-height: 16px;
  min-width: 16px;
  max-height: 16px;
  max-width: 16px;
  border: ${({theme}) => `1px solid ${theme.color.primary}`};
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

export const Check = styled.div`
  height: 8px;
  width: 8px;
  background: ${({theme}) => theme.color.primary};
  border-radius: 50%;
  display: inline-block;
`;

const dataStoreIconMap = {
  [SupportedDataStores.JAEGER]: Jaeger,
  [SupportedDataStores.OpenSearch]: OpenSearch,
  [SupportedDataStores.SignalFX]: SignalFx,
  [SupportedDataStores.TEMPO]: Tempo,
  [SupportedDataStores.OtelCollector]: OtelCollector,
};

export const DataStoreIcon = styled.img.attrs<{$dataStore: SupportedDataStores}>(({$dataStore}) => ({
  src: dataStoreIconMap[$dataStore],
}))<{$dataStore: SupportedDataStores}>``;

export const DataStoreName = styled(Typography.Text)`
  && {
    font-size: ${({theme}) => theme.size.sm};
    font-weight: 700;
  }
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.md};
    font-weight: 700;
    margin: 0;
  }
`;

export const Explanation = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
    font-size: ${({theme}) => theme.size.md};
  }
`;
