import {useCallback, useMemo} from 'react';
import {Typography} from 'antd';
import {noop} from 'lodash';
import {getIsValidUrl, getParsedURL} from 'utils/Common';
import * as S from './UrlDockerTipInput.styled';

interface IProps {
  onChange?(value: string): void;
  value?: string;
}

const localhostDns = ['localhost', '127.0.0.1'];
const dockerHostDns = 'host.docker.internal';

const UrlDockerTipInput = ({onChange = noop, value = ''}: IProps) => {
  const hostname = useMemo(() => {
    try {
      return getParsedURL(value).hostname;
    } catch (e) {
      return '';
    }
  }, [value]);

  const handleReplaceUrl = useCallback(() => {
    if (getIsValidUrl(value)) {
      const url = getParsedURL(value);
      url.hostname = dockerHostDns;
      onChange(url?.toString());
    }
  }, [onChange, value]);

  return !!hostname && localhostDns.includes(hostname) ? (
    <S.Paragraph>
      Are you running Tracetest on Docker?{' '}
      <Typography.Text type="secondary">
        Try replacing <code>{hostname}</code> with <code>host.docker.internal</code> by clicking{'  '}
        <a onClick={() => handleReplaceUrl()}>here.</a>
      </Typography.Text>
    </S.Paragraph>
  ) : null;
};

export default UrlDockerTipInput;
