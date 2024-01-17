import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import Wizard, {TRawWizard} from 'models/Wizard.model';
import {TTestApiEndpointBuilder} from '../Tracetest.api';

export const wizardEndpoints = (builder: TTestApiEndpointBuilder) => ({
  getWizard: builder.query<Wizard, unknown>({
    query: () => ({
      url: '/wizard',
      method: HTTP_METHOD.GET,
      headers: {'content-type': 'application/json'},
    }),
    providesTags: () => [{type: TracetestApiTags.WIZARD, id: 'LIST'}],
    transformResponse: (raw: TRawWizard) => Wizard(raw),
  }),
  updateWizard: builder.mutation<undefined, TRawWizard>({
    query: wizard => ({
      url: `/wizard`,
      method: HTTP_METHOD.PUT,
      body: wizard,
    }),
    invalidatesTags: [{type: TracetestApiTags.WIZARD, id: 'LIST'}],
  }),
});
