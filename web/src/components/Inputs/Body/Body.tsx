import {Radio, Select, Typography} from 'antd';
import {noop} from 'lodash';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from './Body.styled';
import {useBodyMode} from './useBodyMode';
import useLanguageExtensionsMemo from './useLanguageExtensionsMemo';

interface IProps {
  value?: string;
  onChange?(value: string): void;
}

const Body = ({value = '', onChange = noop}: IProps) => {
  const [bodyMode, setBodyMode] = useBodyMode(value);
  const extensions = useLanguageExtensionsMemo(bodyMode);
  const hasNoBody = bodyMode === 'none';

  return (
    <S.Container>
      <S.RadioContainer>
        <Select value={bodyMode} onChange={setBodyMode} data-cy="bodyMode">
          <Radio value="none" data-cy="bodyMode-none">
            none
          </Radio>
          <Radio value="raw" data-cy="bodyMode-raw">
            raw
          </Radio>
          <Radio value="json" data-cy="bodyMode-json">
            JSON
          </Radio>
          <Radio value="xml" data-cy="bodyMode-xml">
            XML
          </Radio>
        </Select>
      </S.RadioContainer>
      {hasNoBody && (
        <div>
          <Typography.Text type="secondary">This request does not have a body</Typography.Text>
        </div>
      )}
      <S.BodyFieldContainer $isDisplaying={hasNoBody} data-cy="body">
        <Editor
          value={value}
          onChange={onChange}
          type={SupportedEditors.Interpolation}
          basicSetup={{lineNumbers: true, indentOnInput: true}}
          extensions={extensions}
          placeholder="Enter request body text"
        />
      </S.BodyFieldContainer>
    </S.Container>
  );
};

export default Body;
