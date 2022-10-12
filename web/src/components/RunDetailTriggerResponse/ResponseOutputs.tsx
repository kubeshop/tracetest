import {useCallback, useState} from 'react';
import {TOutput} from 'types/Output.types';
import SkeletonResponse from './SkeletonResponse';
import * as S from './RunDetailTriggerResponse.styled';
import OutputRow from '../OutputRow';
import OutputModal from '../OutputModal/OutputModal';
import {useConfirmationModal} from '../../providers/ConfirmationModal/ConfirmationModal.provider';

// replace this with the backend response
const testOutputs: TOutput[] = [
  {
    id: '1',
    source: 'trigger',
    attribute: 'body',
  },
  {
    id: '3',
    source: 'trace',
    attribute: 'http.status_code',
    selector: 'span[tracetest.span.type = "http"]:first',
  },
  {
    id: '2',
    source: 'trigger',
    attribute: 'headers.content-type',
  },
  {
    id: '4',
    source: 'trace',
    attribute: 'http.route',
    selector: 'span[tracetest.span.type = "http"]:first',
  },
];

interface IProps {
  outputs?: TOutput[];
}

const ResponseOutputs = ({outputs = testOutputs}: IProps) => {
  const [selectedOutput, setSelectedOutput] = useState<TOutput>();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const {onOpen} = useConfirmationModal();

  const onDelete = useCallback(
    (id: string) => {
      onOpen(`Are you sure you want to delete the output?`, () => console.log('@@confirm delete', id));
    },
    [onOpen]
  );

  const onEdit = useCallback((output: TOutput) => {
    setSelectedOutput(output);
    setIsModalOpen(true);
  }, []);

  const handleSubmit = useCallback((values: TOutput) => {
    setIsModalOpen(false);
    console.log('@@onSubmit', values);
  }, []);

  return !outputs ? (
    <SkeletonResponse />
  ) : (
    <>
      <S.Actions>
        <S.AddOutputButton
          ghost
          type="primary"
          onClick={() => {
            setIsModalOpen(true);
            setSelectedOutput(undefined);
          }}
        >
          Add Output
        </S.AddOutputButton>
      </S.Actions>
      <S.HeadersList>
        {outputs.map(output => (
          <OutputRow output={output} key={output.id} onDelete={onDelete} onEdit={onEdit} />
        ))}
      </S.HeadersList>
      <OutputModal
        draftOutput={selectedOutput}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
      />
    </>
  );
};

export default ResponseOutputs;
