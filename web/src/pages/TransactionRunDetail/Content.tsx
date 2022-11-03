import {/* Tag, */ Typography} from 'antd';
import {useNavigate, useParams} from 'react-router-dom';
import TransactionHeader from '../../components/TransactionHeader';
import {useTransaction} from '../../providers/TransactionRunDetail/TransactionRunDetailProvider';
import * as S from './Content.styled';

function iconBasedOnResult(result: 'success' | 'fail' | 'running', index: number) {
  switch (result) {
    case 'success':
      return <S.IconSuccess />;
    case 'fail':
      return <S.IconSuccess />;
    default:
      return index + 1;
  }
}

const Content = () => {
  const navigate = useNavigate();
  const {transaction} = useTransaction();
  const {transactionId = ''} = useParams();
  return (
    <>
      <TransactionHeader onBack={() => navigate(`/transaction/${transactionId}`)} />
      <S.Wrapper>
        <S.Container>
          <S.SectionLeft>
            <Typography style={{fontSize: 16, marginBottom: 24}}>Execution steps</Typography>

            {/* {transaction?.steps.map(({name, version, ...test}, index) => {
              return (
                <S.Containerr data-cy={`run-card-${name}`} key={test.id}>
                  <div>{iconBasedOnResult(test.result, index)}</div>
                  <S.Info>
                    <S.Title>{`${name} v${version}`}</S.Title>
                    <S.TagContainer>
                      {[test.trigger.method, test.trigger.type].map(d => (
                        <Tag key={d}>{d}</Tag>
                      ))}
                    </S.TagContainer>
                  </S.Info>
                </S.Containerr>
              );
            })} */}
          </S.SectionLeft>
          <S.SectionRight>
            <Typography style={{fontSize: 16, marginBottom: 24}}>Variables</Typography>
            {Object.keys(transaction?.env || {}).map(key => {
              const result = transaction?.env?.[key];
              return (
                <S.Containerr data-cy={`variable-card-${key}`} key={key}>
                  <S.Infoo>
                    <S.Stack>
                      <S.Text opacity={0.6}>Key</S.Text>
                      <S.Text>{`${key} `}</S.Text>
                    </S.Stack>
                    <S.Stack>
                      <S.Text opacity={0.6}>Value</S.Text>
                      <S.Text>{`${result}`}</S.Text>
                    </S.Stack>
                    <S.Stack>
                      <S.Text opacity={0.6}>Coming from</S.Text>
                      <S.Text>{`${'From'}`}</S.Text>
                    </S.Stack>
                  </S.Infoo>
                </S.Containerr>
              );
            })}
          </S.SectionRight>
        </S.Container>
      </S.Wrapper>
    </>
  );
};

export default Content;
