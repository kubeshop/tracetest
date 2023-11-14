import {Form} from 'antd';
import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import {TriggerTypes} from 'constants/Test.constants';
import EntryPointFactory from 'components/TestPlugins/EntryPointFactory';
import * as S from './CreateTest.styled';

interface IProps {
  isLoading: boolean;
  isValid: boolean;
  triggerType: TriggerTypes;
}

const Header = ({isLoading, isValid, triggerType}: IProps) => {
  const form = Form.useFormInstance();

  return (
    <S.Header>
      <S.HeaderLeft>
        <EntryPointFactory type={triggerType} />
      </S.HeaderLeft>

      <S.HeaderRight>
        <AllowButton
          block
          ButtonComponent={CreateButton}
          data-cy="run-test-submit"
          disabled={!isValid}
          loading={isLoading}
          onClick={() => form.submit()}
          operation={Operation.Edit}
          type="primary"
        >
          Run
        </AllowButton>
      </S.HeaderRight>
    </S.Header>
  );
};

export default Header;
