import {Model, TConfigSchemas} from 'types/Common.types';

export type TRawDemo = TConfigSchemas['Demo'];
type Demo = Model<Model<TRawDemo, {}>['spec'], {}>;

const Demo = ({
  spec: {
    id = '',
    name = '',
    type = 'pokeshop',
    enabled,
    pokeshop: {httpEndpoint = '', grpcEndpoint = ''} = {},
    opentelemetryStore: {
      frontendEndpoint = '',
      productCatalogEndpoint = '',
      cartEndpoint = '',
      checkoutEndpoint = '',
    } = {},
  } = {enabled: false},
}: TRawDemo = {}): Demo => {
  return {
    id,
    name,
    type,
    enabled,
    pokeshop: {
      httpEndpoint,
      grpcEndpoint,
    },
    opentelemetryStore: {
      frontendEndpoint,
      productCatalogEndpoint,
      cartEndpoint,
      checkoutEndpoint,
    },
  };
};

export default Demo;
