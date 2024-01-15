import {Typography} from 'antd';
import {COMMUNITY_SLACK_URL} from 'constants/Common.constants';
import {SupportedDataStores} from 'types/DataStore.types';
import * as S from './TracingBackend.styled';
import BackendCard from './BackendCard';

const backends = Object.values(SupportedDataStores);

interface IProps {
  onSelect(backend: SupportedDataStores): void;
}

const Selector = ({onSelect}: IProps) => {
  return (
    <S.Container>
      <Typography.Title level={1}>Select Tracing Backend you would like to use</Typography.Title>

      <Typography.Paragraph type="secondary">
        Tracetest requires configuration details to retrieve traces from your distributed tracing solution.
      </Typography.Paragraph>

      <S.BackendSelector>
        {backends.map(backend => (
          <BackendCard backend={backend} key={backend} onSelect={() => onSelect(backend)} />
        ))}
      </S.BackendSelector>

      <Typography.Paragraph type="secondary">
        Don&apos;t see the Data Stores you need?{' '}
        <a href={COMMUNITY_SLACK_URL} target="__blank">
          Submit a request
        </a>
      </Typography.Paragraph>
    </S.Container>
  );
};

export default Selector;
