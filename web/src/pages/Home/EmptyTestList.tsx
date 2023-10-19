import {Button} from 'antd';
import Empty from 'components/Empty';
import {ADD_TEST_URL, OPENING_TRACETEST_URL} from 'constants/Common.constants';
import {withCustomization} from 'providers/Customization';
import * as S from './Home.styled';

interface IProps {
  onClick(): void;
}

const EmptyTestList = ({onClick}: IProps) => (
  <Empty
    title="Haven't Created a Test Yet"
    message={
      <>
        Hit the &apos;Create&apos; button below to kickstart your testing adventure. Want to learn more about tests?
        Just click{' '}
        <S.Link href={ADD_TEST_URL} target="_blank">
          here
        </S.Link>
        . If you don’t have an app that’s generating OpenTelemetry traces we have a demo for you. Follow these{' '}
        <S.Link href={OPENING_TRACETEST_URL} target="_blank">
          instructions
        </S.Link>
        !
      </>
    }
    action={
      <Button onClick={onClick} type="primary">
        Create Your First Test
      </Button>
    }
  />
);

export default withCustomization(EmptyTestList, 'emptyTestList');
