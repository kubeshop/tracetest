import {Col, Row, Typography} from 'antd';
import {TriggerTypeList} from 'components/TriggerTypeModal';
import {useCallback, useState} from 'react';
import {TriggerTypes} from 'constants/Test.constants';
import {ADD_TEST_URL} from 'constants/Common.constants';
import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
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
        This step involves creating tests by defining test assertions based on your application&apos;s trace data, which
        helps validate the performance and behavior of your services.{' '}
        <a href={ADD_TEST_URL} target="__blank">
          Learn more in docs
        </a>
      </Typography.Paragraph>

      <Typography.Title level={3}>
        Choose the kind of trigger to initiate the trace is presented you will be redirected to the test page
      </Typography.Title>
      <Row>
        <Col span={18}>
          <TriggerTypeList onClick={setSelectedType} />
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
          Next
        </AllowButton>
      </S.ButtonContainer>
    </>
  );
};

export default CreateTest;
