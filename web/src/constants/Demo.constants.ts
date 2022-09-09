import pokeshopProtoData from 'assets/pokeshop.proto.json';
import otelProtoData from 'assets/otel-demo.proto.json';
import pokeshopPostmanData from 'assets/pokeshop.postman_collection.json';
import {HTTP_METHOD, SupportedPlugins} from './Common.constants';

const pokeshopProtoFile = new File([pokeshopProtoData?.proto], 'pokeshop.proto');
const otelProtoFile = new File([otelProtoData?.proto], 'otel-demo.proto');
const pokeshopPostmanFile = new File([JSON.stringify(pokeshopPostmanData)], 'pokeshop.postman_collection.json');

const {
  pokeshopDemoEnabled = 'false',
  pokeshopDemoHostname = '',
  otelDemoEndpoints = '{}',
  otelDemoEnabled = 'false',
} = window.ENV || {};

const isPokeshopEnabled = pokeshopDemoEnabled === 'true';
const isOtelEnabled = otelDemoEnabled === 'true';

const {
  frontend = '',
  productCatalog = '',
  cart = '',
  checkout = '',
}: Record<string, string> = JSON.parse(otelDemoEndpoints);
const userId = '2491f868-88f1-4345-8836-d5d8511a9f83';

export const PokeshopDemo = {
  [SupportedPlugins.REST]: [
    {
      name: 'Pokeshop - List',
      url: `http://${pokeshopDemoHostname}/pokemon?take=20&skip=0`,
      method: HTTP_METHOD.GET,
      body: '',
      description: 'Get a Pokemon',
    },
    {
      name: 'Pokeshop - Add',
      url: `http:/${pokeshopDemoHostname}/pokemon`,
      method: HTTP_METHOD.POST,
      body: '/{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
      description: 'Add a Pokemon',
    },
    {
      name: 'Pokeshop - Import',
      url: `http://${pokeshopDemoHostname}/pokemon/import`,
      method: HTTP_METHOD.POST,
      body: '{"id":52}',
      description: 'Import a Pokemon',
    },
  ],
  [SupportedPlugins.GRPC]: [
    {
      name: 'Pokeshop - List',
      url: `${pokeshopDemoHostname}:8082`,
      message: '',
      method: 'pokeshop.Pokeshop.getPokemonList',
      description: 'Get a Pokemon',
      protoFile: pokeshopProtoFile,
    },
    {
      name: 'Pokeshop - Add',
      url: `${pokeshopDemoHostname}:8082`,
      message:
        '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
      method: 'pokeshop.Pokeshop.createPokemon',
      protoFile: pokeshopProtoFile,
      description: 'Add a Pokemon',
    },
    {
      name: 'Pokeshop - Import',
      url: `${pokeshopDemoHostname}:8082`,
      message: '{"id":52}',
      method: 'pokeshop.Pokeshop.importPokemon',
      protoFile: pokeshopProtoFile,
      description: 'Import a Pokemon',
    },
  ],
  [SupportedPlugins.Postman]: [
    {
      name: 'Pokeshop - List',
      url: `http://${pokeshopDemoHostname}/pokemon?take=20&skip=0`,
      method: HTTP_METHOD.GET,
      body: '',
      description: 'Get a Pokemon',
      collectionTest: 'List',
      collectionFile: pokeshopPostmanFile,
    },
    {
      name: 'Pokeshop - Add',
      url: `http://${pokeshopDemoHostname}/pokemon`,
      method: HTTP_METHOD.POST,
      body: '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
      description: 'Add a Pokemon',
      collectionTest: 'Create',
      collectionFile: pokeshopPostmanFile,
    },
    {
      name: 'Pokeshop - Import',
      url: `http://${pokeshopDemoHostname}/pokemon/import`,
      method: HTTP_METHOD.POST,
      body: '{"id":52}',
      description: 'Import a Pokemon',
      collectionTest: 'Import',
      collectionFile: pokeshopPostmanFile,
    },
  ],
};

export const OtelDemo = {
  [SupportedPlugins.REST]: [
    {
      name: 'Otel - List Products',
      url: `http://${frontend}/api/products`,
      method: HTTP_METHOD.GET,
      body: '',
      description: 'Otel - List Products',
    },
    {
      name: 'Otel - Get Product',
      url: `http://${frontend}/api/products/OLJCESPC7Z`,
      method: HTTP_METHOD.GET,
      body: '',
      description: 'Otel - Get Product',
    },
    {
      name: 'Otel - Add To Cart',
      url: `http://${frontend}/api/cart`,
      method: HTTP_METHOD.POST,
      body: JSON.stringify({
        item: {productId: 'OLJCESPC7Z', quantity: 1},
        userId,
      }),
      description: 'Otel - Add To Cart',
    },
    {
      name: 'Otel - Get Cart',
      url: `http://${frontend}/api/cart?sessionId=${userId}`,
      method: HTTP_METHOD.GET,
      body: '',
      description: 'Otel - Get Cart',
    },
    {
      name: 'Otel - Checkout',
      url: `http://${frontend}/api/checkout`,
      method: HTTP_METHOD.POST,
      body: JSON.stringify({
        userId,
        email: 'someone@example.com',
        address: {
          streetAddress: '1600 Amphitheatre Parkway',
          state: 'CA',
          country: 'United States',
          city: 'Mountain View',
          zipCode: 94043,
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
      url: productCatalog,
      message: '',
      method: 'hipstershop.ProductCatalogService.ListProducts',
      description: 'Otel - List Products',
      protoFile: otelProtoFile,
    },
    {
      name: 'Otel - Get Product',
      url: productCatalog,
      message: '{"id": "OLJCESPC7Z"}',
      method: 'hipstershop.ProductCatalogService.GetProduct',
      description: 'Otel - Get Product',
      protoFile: otelProtoFile,
    },
    {
      name: 'Otel - Add To Cart',
      url: cart,
      message: JSON.stringify({item: {product_id: 'OLJCESPC7Z', quantity: 1}, user_id: userId}),
      method: 'hipstershop.CartService.AddItem',
      description: 'Otel - Add To Cart',
      protoFile: otelProtoFile,
    },
    {
      name: 'Otel - Get Cart',
      url: cart,
      message: `{"user_id": "${userId}"}`,
      method: 'hipstershop.CartService.GetCart',
      description: 'Otel - Get Cart',
      protoFile: otelProtoFile,
    },
    {
      name: 'Otel - Checkout',
      url: checkout,
      message: JSON.stringify({
        user_id: userId,
        user_currency: 'USD',
        address: {
          street_address: '1600 Amphitheatre Parkway',
          state: 'CA',
          country: 'United States',
          city: 'Mountain View',
          zip_code: 94043,
        },
        email: 'someone@example.com',
        credit_card: {
          credit_card_number: '4432-8015-6152-0454',
          credit_card_cvv: 672,
          credit_card_expiration_year: 2030,
          credit_card_expiration_month: 1,
        },
      }),
      method: 'hipstershop.CheckoutService.PlaceOrder',
      description: 'Otel - Checkout',
      protoFile: otelProtoFile,
    },
  ],
};

export const DemoByPluginMap = {
  [SupportedPlugins.REST]: [
    ...((isPokeshopEnabled && PokeshopDemo[SupportedPlugins.REST]) || []),
    ...((isOtelEnabled && OtelDemo[SupportedPlugins.REST]) || []),
  ],
  [SupportedPlugins.GRPC]: [
    ...((isPokeshopEnabled && PokeshopDemo[SupportedPlugins.GRPC]) || []),
    ...((isOtelEnabled && OtelDemo[SupportedPlugins.GRPC]) || []),
  ],
  [SupportedPlugins.Postman]: (isPokeshopEnabled && PokeshopDemo[SupportedPlugins.Postman]) || [],
};
