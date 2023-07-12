import {snakeCase} from 'lodash';
import Test from 'models/Test.model';
import {useState} from 'react';
import TestRun from 'models/TestRun.model';
import * as S from './RunDetailAutomate.styled';
import RunDetailAutomateDefinition from '../RunDetailAutomateDefinition';
import RunDetailAutomateMethods from '../RunDetailAutomateMethods/RunDetailAutomateMethods';

interface IProps {
  test: Test;
  run: TestRun;
}

const RunDetailAutomate = ({test, run}: IProps) => {
  const [fileName, setFileName] = useState<string>(`${snakeCase(test.name)}.yaml`);

  return (
    <S.Container>
      <S.SectionLeft>
        <RunDetailAutomateDefinition onFileNameChange={setFileName} fileName={fileName} test={test} />
      </S.SectionLeft>
      <S.SectionRight>
        <RunDetailAutomateMethods fileName={fileName} test={test} run={run} />
      </S.SectionRight>
    </S.Container>
  );
};

export default RunDetailAutomate;
