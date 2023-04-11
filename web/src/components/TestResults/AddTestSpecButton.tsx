import React, {useCallback, useEffect, useMemo, useRef} from 'react';
import {CaretDownOutlined} from '@ant-design/icons';
import {Dropdown, Menu} from 'antd';
import SpanService from 'services/Span.service';
import Span from 'models/Span.model';
import {TEST_SPEC_SNIPPETS, TSnippet} from 'constants/TestSpecs.constants';
import * as S from './TestResults.styled';
import {useTestSpecForm} from '../TestSpecForm/TestSpecForm.provider';

interface IProps {
  selectedSpan: Span;
  visibleByDefault?: boolean;
}

const AddTestSpecButton = ({selectedSpan, visibleByDefault = false}: IProps) => {
  const {open} = useTestSpecForm();
  const caretRef = useRef<HTMLElement>(null);
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

  useEffect(() => {
    if (visibleByDefault && caretRef.current) {
      caretRef.current?.click();
    }
  }, [visibleByDefault]);

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
            label: 'Empty Test Spec',
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
      overlay={menu}
      trigger={['click']}
      placement="bottomRight"
      onClick={handleEmptyTestSpec}
      type="primary"
      buttonsRender={([leftButton]) => [
        React.cloneElement(leftButton as React.ReactElement<any, string>, {'data-cy': 'add-test-spec-button'}),
        <S.CaretDropdownButton ref={caretRef} type="primary" data-cy="create-button">
          <CaretDownOutlined />
        </S.CaretDropdownButton>,
      ]}
    >
      Add Test Spec
    </Dropdown.Button>
  );
};

export default AddTestSpecButton;
