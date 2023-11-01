import {Form, Radio} from 'antd';
import {useState} from 'react';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import {BodyFieldContainer} from './BodyFieldContainer';
import * as S from './BodyField.styled';
import {useBodyMode} from './useBodyMode';
import {useLanguageExtensionsMemo} from './useLanguageExtensionsMemo';

interface IProps {
  body?: string;
  setBody: (body?: string) => void;
}

const BodyField = ({body, setBody}: IProps) => {
  const [bodyMode, setBodyMode] = useBodyMode(body);
  const [buffer, setBuffer] = useState<undefined | string>(undefined);
  const extensions = useLanguageExtensionsMemo(bodyMode);
  const hasNoBody = bodyMode === 'none';

  return (
    <S.Container>
      <S.RadioContainer>
        <S.SingleLabel label="Request body">{buffer}</S.SingleLabel>
        <Radio.Group
          value={bodyMode}
          onChange={e => {
            if (e.target.value === 'none') {
              setBuffer(body);
              setBody(undefined);
            } else if (buffer) {
              setBody(buffer);
              setBuffer(undefined);
            }
            setBodyMode(e.target.value);
          }}
        >
          <Radio value="none" data-cy="bodyMode-none">
            None
          </Radio>
          <Radio value="raw" data-cy="bodyMode-raw">
            Raw
          </Radio>
          <Radio value="json" data-cy="bodyMode-json">
            JSON
          </Radio>
          <Radio value="xml" data-cy="bodyMode-xml">
            XML
          </Radio>
        </Radio.Group>
      </S.RadioContainer>
      {hasNoBody && (
        <div>
          <h4>This request has no body {bodyMode}</h4>
        </div>
      )}
      <BodyFieldContainer $isDisplaying={hasNoBody}>
        <Form.Item name="body" data-cy="body">
          <Editor
            type={SupportedEditors.Interpolation}
            basicSetup={{lineNumbers: true, indentOnInput: true}}
            extensions={extensions}
            placeholder="Enter request body text"
          />
        </Form.Item>
      </BodyFieldContainer>
    </S.Container>
  );
};

export default BodyField;
