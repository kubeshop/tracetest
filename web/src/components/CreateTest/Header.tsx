import {LoadingOutlined} from '@ant-design/icons';
import {Form} from 'antd';
import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import {TriggerTypes} from 'constants/Test.constants';
import EntryPointFactory from 'components/TestPlugins/EntryPointFactory';
import Test from 'models/Test.model';
import * as S from './CreateTest.styled';

interface IProps {
  isLoading: boolean;
  isRunStateFinished?: boolean;
  isValid: boolean;
  onRunTest?(): void;
  onStopTest?(): void;
  triggerType: TriggerTypes;
}

const Header = ({isLoading, isRunStateFinished, isValid, onRunTest, onStopTest, triggerType}: IProps) => {
  const form = Form.useFormInstance();

  const handleOnRunClick = () => {
    if (isRunStateFinished) {
      onRunTest?.();
      form.submit();
      return;
    }
    onStopTest?.();
  };

  return (
    <S.Header>
      <S.HeaderLeft>
        <EntryPointFactory type={triggerType} />
      </S.HeaderLeft>

      <S.HeaderRight>
        {Test.shouldAllowRun(triggerType) && (
          <AllowButton
            block
            ButtonComponent={CreateButton}
            data-cy="run-test-submit"
            disabled={!isValid}
            icon={isLoading && <LoadingOutlined />}
            onClick={handleOnRunClick}
            operation={Operation.Edit}
            type="primary"
          >
            {isRunStateFinished ? 'Run' : 'Stop'}
          </AllowButton>
        )}
      </S.HeaderRight>
    </S.Header>
  );
};

export default Header;
