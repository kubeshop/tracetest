import {CaretDownOutlined} from '@ant-design/icons';
import {Dropdown, Menu} from 'antd';
import React, {useCallback, useMemo} from 'react';

import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {TEST_SPEC_SNIPPETS, TSnippet} from 'constants/TestSpecs.constants';
import Span from 'models/Span.model';
import {isRunStateSucceeded} from 'models/TestRun.model';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import SpanService from 'services/Span.service';
import * as S from './TestResults.styled';

interface IProps {
  selectedSpan: Span;
}

const AddTestSpecButton = ({selectedSpan}: IProps) => {
  const {
    run: {state},
  } = useTestRun();
  const {open} = useTestSpecForm();

  const handleEmptyTestSpec = useCallback(() => {
    const selector = SpanService.getSelectorInformation(selectedSpan);

    open({
      isEditing: false,
      selector,
      defaultValues: {
        selector,
      },
    });
  }, [open, selectedSpan]);

  const onSnippetClick = useCallback(
    (snippet: TSnippet) => {
      open({
        isEditing: false,
        selector: snippet.selector,
        defaultValues: snippet,
      });
    },
    [open]
  );

  const menu = useMemo(
    () => (
      <Menu
        items={[
          {
            label: 'Try these snippets for quick testing:',
            key: 'test',
            type: 'group',
            children: TEST_SPEC_SNIPPETS.map(snippet => ({
              label: snippet.name,
              key: snippet.name,
              onClick: () => onSnippetClick(snippet),
            })),
          },
          {type: 'divider'},
          {
            label: 'Add empty Test Spec',
            key: 'empty-test-spec',
            onClick: handleEmptyTestSpec,
          },
        ]}
      />
    ),
    [handleEmptyTestSpec, onSnippetClick]
  );

  return (
    <Dropdown.Button
      disabled={!isRunStateSucceeded(state)}
      overlay={menu}
      trigger={['click']}
      placement="bottomRight"
      onClick={handleEmptyTestSpec}
      type="primary"
      buttonsRender={([leftButton]) => [
        React.cloneElement(leftButton as React.ReactElement<any, string>, {'data-cy': 'add-test-spec-button'}),
        <S.CaretDropdownButton type="primary" data-cy="create-button">
          <CaretDownOutlined />
        </S.CaretDropdownButton>,
      ]}
    >
      Add Test Spec
    </Dropdown.Button>
  );
};

export default AddTestSpecButton;
