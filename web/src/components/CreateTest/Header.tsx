import CreateButton from 'components/CreateButton';
import {Form} from 'antd';
import {TriggerTypes} from 'constants/Test.constants';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import EntryPointFactory from 'components/TestPlugins/EntryPointFactory';
import AllowButton, {Operation} from '../AllowButton';
import * as S from './CreateTest.styled';

interface IProps {
  triggerType: TriggerTypes;
  isValid: boolean;
}

const Header = ({triggerType, isValid}: IProps) => {
  const {isLoading} = useCreateTest();
  const form = Form.useFormInstance();

  return (
    <S.Header>
      <S.HeaderLeft>
        <EntryPointFactory type={triggerType} />
      </S.HeaderLeft>

      <S.HeaderRight>
        <AllowButton
          operation={Operation.Edit}
          block
          ButtonComponent={CreateButton}
          data-cy="run-test-submit"
          disabled={!isValid}
          loading={isLoading}
          onClick={() => form.submit()}
          type="primary"
        >
          Run
        </AllowButton>
      </S.HeaderRight>
    </S.Header>
  );
};

export default Header;
