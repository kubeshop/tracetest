import {Col, Row, Typography} from 'antd';
import {TriggerTypeList} from 'components/TriggerTypeModal';
import {useCallback, useState} from 'react';
import {TriggerTypes} from 'constants/Test.constants';
import {ADD_TEST_URL} from 'constants/Common.constants';
import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import WizardAnalytics from 'services/Analytics/Wizard.service';
import * as S from './CreateTest.styled';

const CreateTest = () => {
  const [selectedType, setSelectedType] = useState<TriggerTypes | undefined>(undefined);
  const {navigate} = useDashboard();

  const handleNext = useCallback(() => {
    navigate(`test/create/${selectedType}`);
  }, [navigate, selectedType]);

  return (
    <>
      <Typography.Title>Create your first test</Typography.Title>
      <Typography.Paragraph type="secondary">
        At this point you are ready to run your first Tracetest test. Select the trigger type and create a test!{' '}
        <a href={ADD_TEST_URL} target="__blank">
          Learn more about it in our docs.
        </a>
      </Typography.Paragraph>

      <Typography.Title level={3}>Choose the kind of trigger you want to use to start your test</Typography.Title>
      <Row>
        <Col span={18}>
          <TriggerTypeList
            onClick={(plugin: TriggerTypes) => {
              setSelectedType(plugin);
              WizardAnalytics.onTriggerTypeSelect(plugin);
            }}
          />
        </Col>
      </Row>
      <S.ButtonContainer>
        <AllowButton
          ButtonComponent={CreateButton}
          data-cy="wizard-create-test-button"
          disabled={!selectedType}
          onClick={handleNext}
          operation={Operation.Edit}
          type="primary"
        >
          Continue
        </AllowButton>
      </S.ButtonContainer>
    </>
  );
};

export default CreateTest;
