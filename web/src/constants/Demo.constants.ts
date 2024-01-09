import Demo from 'models/Demo.model';
import SettingService from 'services/Setting.service';
import {SupportedDemos} from 'types/Settings.types';
import {HTTP_METHOD, SupportedPlugins} from './Common.constants';
import pokeshopProtoData from './demos/pokeshop.proto';
import otelDemoProtoData from './demos/otel-demo.proto';

const pokeshopProtoFile = new File([pokeshopProtoData?.proto], 'pokeshop.proto');
const otelProtoFile = new File([otelDemoProtoData?.proto], 'otel-demo.proto');

const userId = '2491f868-88f1-4345-8836-d5d8511a9f83';

export function getPokeshopDemo(demoSettings: Demo) {
  const {
    pokeshop: {httpEndpoint: pokeshopHttp = '', grpcEndpoint: pokeshopGrpc = '', kafkaBroker: pokeshopKafka = ''},
  } = demoSettings;

  return {
    [SupportedPlugins.REST]: [
      {
        name: 'Pokeshop - List',
        url: `${pokeshopHttp}/pokemon?take=20&skip=0`,
        method: HTTP_METHOD.GET,
        body: '',
        description: 'Get a Pokemon',
      },
      {
        name: 'Pokeshop - Add',
        url: `${pokeshopHttp}/pokemon`,
        method: HTTP_METHOD.POST,
        body: '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
        description: 'Add a Pokemon',
      },
      {
        name: 'Pokeshop - Import',
        url: `${pokeshopHttp}/pokemon/import`,
        method: HTTP_METHOD.POST,
        body: '{"id":52}',
        description: 'Import a Pokemon',
      },
    ],
    [SupportedPlugins.GRPC]: [
      {
        name: 'Pokeshop - List',
        url: pokeshopGrpc,
        message: '',
        method: 'pokeshop.Pokeshop.getPokemonList',
        description: 'Get a Pokemon',
        protoFile: pokeshopProtoFile,
      },
      {
        name: 'Pokeshop - Add',
        url: pokeshopGrpc,
        message:
          '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
        method: 'pokeshop.Pokeshop.createPokemon',
        protoFile: pokeshopProtoFile,
        description: 'Add a Pokemon',
      },
      {
        name: 'Pokeshop - Import',
        url: pokeshopGrpc,
        message: '{"id":52}',
        method: 'pokeshop.Pokeshop.importPokemon',
        protoFile: pokeshopProtoFile,
        description: 'Import a Pokemon',
      },
    ],
    [SupportedPlugins.Kafka]: [
      {
        name: 'Pokeshop - Import from Stream',
        brokerUrls: [`${pokeshopKafka}`],
        topic: 'pokemon',
        headers: [],
        messageKey: 'snorlax-key',
        messageValue: '{"id":143}',
        description: 'Import a Pokemon via Stream',
      },
    ],
  };
}

export function getOtelDemo(demoSettings: Demo) {
  const {
    opentelemetryStore: {
      cartEndpoint: otelCart = '',
      checkoutEndpoint: otelCheckout = '',
      frontendEndpoint: otelFrontend = '',
      productCatalogEndpoint: otelProductCatalog = '',
    },
  } = demoSettings;

  return {
    [SupportedPlugins.REST]: [
      {
        name: 'Otel - List Products',
        url: `${otelFrontend}/api/products`,
        method: HTTP_METHOD.GET,
        body: '',
        description: 'Otel - List Products',
      },
      {
        name: 'Otel - Get Product',
        url: `${otelFrontend}/api/products/OLJCESPC7Z`,
        method: HTTP_METHOD.GET,
        body: '',
        description: 'Otel - Get Product',
      },
      {
        name: 'Otel - Add To Cart',
        url: `${otelFrontend}/api/cart`,
        method: HTTP_METHOD.POST,
        body: JSON.stringify({
          item: {productId: 'OLJCESPC7Z', quantity: 1},
          userId,
        }),
        description: 'Otel - Add To Cart',
      },
      {
        name: 'Otel - Get Cart',
        url: `${otelFrontend}/api/cart?sessionId=${userId}`,
        method: HTTP_METHOD.GET,
        body: '',
        description: 'Otel - Get Cart',
      },
      {
        name: 'Otel - Checkout',
        url: `${otelFrontend}/api/checkout`,
        method: HTTP_METHOD.POST,
        body: JSON.stringify({
          userId,
          email: 'someone@example.com',
          address: {
            streetAddress: '1600 Amphitheatre Parkway',
            state: 'CA',
            country: 'United States',
            city: 'Mountain View',
            zipCode: '94043',
          },
          userCurrency: 'USD',
          creditCard: {
            creditCardCvv: 672,
            creditCardExpirationMonth: 1,
            creditCardExpirationYear: 2030,
            creditCardNumber: '4432-8015-6152-0454',
          },
        }),
        description: 'Otel - Checkout',
      },
    ],
    [SupportedPlugins.GRPC]: [
      {
        name: 'Otel - List Products',
        url: otelProductCatalog,
        message: '',
        method: 'oteldemo.ProductCatalogService.ListProducts',
        description: 'Otel - List Products',
        protoFile: otelProtoFile,
      },
      {
        name: 'Otel - Get Product',
        url: otelProductCatalog,
        message: '{"id": "OLJCESPC7Z"}',
        method: 'oteldemo.ProductCatalogService.GetProduct',
        description: 'Otel - Get Product',
        protoFile: otelProtoFile,
      },
      {
        name: 'Otel - Add To Cart',
        url: otelCart,
        message: JSON.stringify({item: {product_id: 'OLJCESPC7Z', quantity: 1}, user_id: userId}),
        method: 'oteldemo.CartService.AddItem',
        description: 'Otel - Add To Cart',
        protoFile: otelProtoFile,
      },
      {
        name: 'Otel - Get Cart',
        url: otelCart,
        message: `{"user_id": "${userId}"}`,
        method: 'oteldemo.CartService.GetCart',
        description: 'Otel - Get Cart',
        protoFile: otelProtoFile,
      },
      {
        name: 'Otel - Checkout',
        url: otelCheckout,
        message: JSON.stringify({
          user_id: userId,
          user_currency: 'USD',
          address: {
            street_address: '1600 Amphitheatre Parkway',
            state: 'CA',
            country: 'United States',
            city: 'Mountain View',
            zip_code: '94043',
          },
          email: 'someone@example.com',
          credit_card: {
            credit_card_number: '4432-8015-6152-0454',
            credit_card_cvv: 672,
            credit_card_expiration_year: 2030,
            credit_card_expiration_month: 1,
          },
        }),
        method: 'oteldemo.CheckoutService.PlaceOrder',
        description: 'Otel - Checkout',
        protoFile: otelProtoFile,
      },
    ],
  };
}

export function getDemoByPluginMap(demos: Demo[]) {
  const enabledDemos = SettingService.getEnabledDemos(demos);
  const pokeShopDemo = enabledDemos.find(demo => demo.type === SupportedDemos.Pokeshop);
  const otelDemo = enabledDemos.find(demo => demo.type === SupportedDemos.OpentelemetryStore);

  const pokeshopDemoMap = pokeShopDemo ? getPokeshopDemo(pokeShopDemo) : undefined;
  const otelDemoMap = otelDemo ? getOtelDemo(otelDemo) : undefined;

  return {
    [SupportedPlugins.REST]: [
      ...((pokeshopDemoMap && pokeshopDemoMap[SupportedPlugins.REST]) || []),
      ...((otelDemoMap && otelDemoMap[SupportedPlugins.REST]) || []),
    ],
    [SupportedPlugins.GRPC]: [
      ...((pokeshopDemoMap && pokeshopDemoMap[SupportedPlugins.GRPC]) || []),
      ...((otelDemoMap && otelDemoMap[SupportedPlugins.GRPC]) || []),
    ],
    [SupportedPlugins.TraceID]: [],
    [SupportedPlugins.Cypress]: [],
    [SupportedPlugins.Playwright]: [],
    [SupportedPlugins.Kafka]: (pokeshopDemoMap && pokeshopDemoMap[SupportedPlugins.Kafka]) || [],
  };
}
