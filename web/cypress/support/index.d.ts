declare namespace Cypress {
  interface Chainable {
    createMultipleTestRuns(id: string, count: number): Chainable<Element>;
    createAssertion(): Chainable<Element>;
    openTestCreationModal(): Chainable<Element>;
    interceptTracePageApiCalls(): Chainable<Element>;
    interceptHomeApiCall(): Chainable<Element>;
    waitForTracePageApiCalls(): Chainable<Element>;
    createTest(): Chainable<Element>;
    clickNextOnCreateTestWizard(): Chainable<Element>;
    selectTestFromDemoList(): Chainable<Element>;
    selectOperator(index: number, text?: string): Chainable<Element>;
    editTestByTestId(testId: string): Chainable<Element>;
    submitCreateForm(mode?: string): Chainable<Element>;
    deleteTest(shouldIntercept?: boolean): Chainable<Element>;
    makeSureUserIsOnTracePage(): Chainable<Element>;
    cancelOnBoarding(): Chainable<Element>;
    makeSureUserIsOnTestDetailPage(): Chainable<Element>;
    goToTestDetailPageAndRunTest(pathname: string): Chainable<Element>;
    matchTestRunPageUrl(): Chainable<Element>;
    createTestByName(name: string): Chainable<Element>;
    submitAndMakeSureTestIsCreated(name: string): Chainable<Element>;
    createTestWithAuth(authMethod: string, keys: string[]): Chainable<string>;
    fillCreateFormBasicStep(name: string, mode?: string): Chainable<Element>;
    setCreateFormUrl(method: string, url: string): Chainable<Element>;
    selectRunDetailMode(index: number): Chainable<Element>;
    interceptEditTestCall(): Chainable<Element>;
    deleteTestSuiteTests(): Chainable<Element>;
    openTestSuiteCreationModal(): Chainable<Element>;
    deleteTestSuite(): Chainable<Element>;
    enableDemo(): Chainable<Element>;
  }
}
