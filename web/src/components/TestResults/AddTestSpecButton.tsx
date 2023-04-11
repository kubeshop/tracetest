import React, {useCallback, useMemo} from 'react';
import {CaretDownOutlined} from '@ant-design/icons';
import {Dropdown, Menu} from 'antd';
import SpanService from 'services/Span.service';
import Span from 'models/Span.model';
import * as S from './TestResults.styled';
import {useTestSpecForm} from '../TestSpecForm/TestSpecForm.provider';
import {TEST_SPEC_SNIPPETS, TSnippet} from '../../constants/TestSpecs.constants';

interface IProps {
  selectedSpan: Span;
}

const AddTestSpecButton = ({selectedSpan}: IProps) => {
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
        ]}
      />
    ),
    [onSnippetClick]
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
