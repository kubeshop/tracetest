import {Form, FormInstance} from 'antd';
import {useEffect, useMemo} from 'react';
import {debounce} from 'lodash';
import {useSpan} from 'providers/Span/Span.provider';
import {useLazyGetSelectedSpansQuery} from 'redux/apis/TraceTest.api';
import SelectorService from 'services/Selector.service';
import {useAppSelector} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TPseudoSelector, TSpanSelector} from 'types/Assertion.types';
import {IValues} from './AssertionForm';
import * as S from './AssertionForm.styled';
import AssertionFormSelectorInput from './AssertionFormSelectorInput';
import AssertionFormPseudoSelectorInput from './AssertionFormPseudoSelectorInput';
import useAssertionFormValues from './hooks/useAssertionFormValues';
import AdvancedEditor from '../AdvancedEditor';

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
  const {onSetAffectedSpans, onClearAffectedSpans} = useSpan();
  const [onTriggerSelectedSpans, {data: spanIdList = [], isError: isInvalidSelector}] = useLazyGetSelectedSpansQuery();

  const {currentIsAdvancedSelector, currentPseudoSelector, currentSelector, currentSelectorList} =
    useAssertionFormValues(form);

  const query = useMemo(
    () =>
      currentIsAdvancedSelector
        ? currentSelector
        : SelectorService.getSelectorString(currentSelectorList || [], currentPseudoSelector),
    [currentIsAdvancedSelector, currentPseudoSelector, currentSelector, currentSelectorList]
  );

  const handleSelector = useMemo(
    () =>
      debounce(async ({q, tId, rId}: {q: string; rId: string; tId: string}) => {
        const idList = await onTriggerSelectedSpans({
          query: q,
          testId: tId,
          runId: rId,
        }).unwrap();

        onSetAffectedSpans(idList);
      }, 500),
    [onSetAffectedSpans, onTriggerSelectedSpans]
  );

  useEffect(() => {
    handleSelector({q: query, tId: testId, rId: runId});
  }, [handleSelector, query, runId, testId]);

  useEffect(() => {
    form.setFields([
      {
        name: 'selector',
        errors: isInvalidSelector ? ['Invalid selector'] : [],
      },
    ]);
    onValidSelector(!isInvalidSelector);
  }, [form, isInvalidSelector, onValidSelector]);

  useEffect(() => {
    return () => {
      onClearAffectedSpans();
    };
  }, []);

  const selectorAttributeList = useAppSelector(state =>
    AssertionSelectors.selectSelectorAttributeList(state, testId, runId, spanIdList, currentSelectorList)
  );
  const definitionSelectorList = useAppSelector(state => TestDefinitionSelectors.selectDefinitionSelectorList(state));

  return !currentIsAdvancedSelector ? (
    <S.SelectorInputContainer>
      <Form.Item
        name="selectorList"
        rules={[
          {required: true, message: 'At least one selector is required'},
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
        rules={[{required: true, message: 'The selector cannot be empty'}]}
        validateTrigger={[]}
        hasFeedback
        help={isInvalidSelector ? 'Invalid selector' : ''}
        validateStatus={isInvalidSelector ? 'error' : ''}
      >
        <AdvancedEditor runId={runId} testId={testId} />
      </Form.Item>
    </S.AdvancedSelectorInputContainer>
  );
};

export default AssertionFormSelector;
