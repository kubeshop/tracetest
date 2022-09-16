import CodeMirror from '@uiw/react-codemirror';
import {Form, Radio} from 'antd';
import {BodyFieldContainer} from './BodyFieldContainer';
import {SingleLabel} from './SingleLabel';
import {useBodyMode} from './useBodyMode';
import {useLanguageExtensionsMemo} from './useLanguageExtensionsMemo';

interface IProps {
  isEditing?: boolean;
  body?: string;
}

export const BodyField = ({isEditing = false, body}: IProps): React.ReactElement => {
  const [bodyMode, setBodyMode] = useBodyMode(isEditing, body);
  const extensions = useLanguageExtensionsMemo(bodyMode);
  const hasNoBody = bodyMode === 'none';
  return (
    <>
      <span>
        <SingleLabel label="Request body">{null}</SingleLabel>
        <Radio.Group value={bodyMode} onChange={e => setBodyMode(e.target.value)}>
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
      </span>
      {hasNoBody && (
        <div>
          <h4>This request has no body {bodyMode}</h4>
        </div>
      )}
      <BodyFieldContainer $isDisplaying={hasNoBody}>
        <Form.Item name="body">
          <CodeMirror
            data-cy="body"
            basicSetup={{lineNumbers: true, indentOnInput: true}}
            extensions={extensions}
            spellCheck={false}
            placeholder={`Enter request body text `}
          />
        </Form.Item>
      </BodyFieldContainer>
    </>
  );
};
