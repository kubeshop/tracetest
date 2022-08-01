import {Form, FormInstance} from 'antd';
import SelectorService from 'services/Selector.service';
import {useAppSelector} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TPseudoSelector, TSpanSelector} from 'types/Assertion.types';
import {IValues} from './AssertionForm';
import * as S from './AssertionForm.styled';
import AssertionFormSelectorInput from './AssertionFormSelectorInput';
import AssertionFormPseudoSelectorInput from './AssertionFormPseudoSelectorInput';
import AdvancedEditor from '../AdvancedEditor';
import useQuerySelector from './hooks/useQuerySelector';
import useAssertionFormValues from './hooks/useAssertionFormValues';

interface IProps {
  form: FormInstance<IValues>;
  testId: string;
  runId: string;
  isEditing: boolean;
  selectorList: TSpanSelector[];
  pseudoSelector?: TPseudoSelector;
  onValidSelector(isValid: boolean): void;
}

const AssertionFormSelector = ({
  form,
  testId,
  runId,
  isEditing,
  selectorList,
  pseudoSelector,
  onValidSelector,
}: IProps) => {
  const {spanIdList, isValid} = useQuerySelector({
    form,
    runId,
    testId,
    onValidSelector,
  });
  const {currentIsAdvancedSelector, currentPseudoSelector, currentSelectorList} = useAssertionFormValues(form);

  const selectorAttributeList = useAppSelector(state =>
    AssertionSelectors.selectSelectorAttributeList(state, testId, runId, spanIdList, currentSelectorList)
  );
  const definitionSelectorList = useAppSelector(state => TestDefinitionSelectors.selectDefinitionSelectorList(state));

  return !currentIsAdvancedSelector ? (
    <S.SelectorInputContainer>
      <Form.Item
        name="selectorList"
        rules={[
          {
            validator: (_, value: TSpanSelector[]) =>
              SelectorService.validateSelector(
                definitionSelectorList,
                isEditing,
                selectorList,
                value,
                currentPseudoSelector,
                pseudoSelector
              ),
          },
        ]}
      >
        <AssertionFormSelectorInput attributeList={selectorAttributeList} />
      </Form.Item>
      <Form.Item name="pseudoSelector">
        <AssertionFormPseudoSelectorInput />
      </Form.Item>
    </S.SelectorInputContainer>
  ) : (
    <S.AdvancedSelectorInputContainer>
      <Form.Item
        name="selector"
        validateTrigger={[]}
        hasFeedback
        help={!isValid ? 'Invalid selector' : ''}
        validateStatus={!isValid ? 'error' : ''}
      >
        <AdvancedEditor runId={runId} testId={testId} />
      </Form.Item>
    </S.AdvancedSelectorInputContainer>
  );
};

export default AssertionFormSelector;
