declare namespace Cypress {
  interface Chainable {
    createMultipleTestRuns(id: string, count: number): Chainable<Element>;
    createAssertion(index?: number): Chainable<Element>;
    openTestCreationModal(): Chainable<Element>;
    interceptTracePageApiCalls(): Chainable<Element>;
    inteceptHomeApiCall(): Chainable<Element>;
    waitForTracePageApiCalls(): Chainable<Element>;
    createTest(): Chainable<Element>;
    clickNextOnCreateTestWizard(): Chainable<Element>;
    selectTestFromDemoList(): Chainable<Element>;
    selectOperator(index: number, text?: string): Chainable<Element>;
    editTestByTestId(testId: string): Chainable<Element>;
    submitCreateTestForm(): Chainable<Element>;
    deleteTest(shouldIntercept?: boolean): Chainable<Element>;
    makeSureUserIsOnTracePage(shouldCancelOnboarding?: boolean): Chainable<Element>;
    cancelOnBoarding(): Chainable<Element>;
    makeSureUserIsOnTestDetailPage(): Chainable<Element>;
    goToTestDetailPageAndRunTest(pathname: string): Chainable<Element>;
    matchTestRunPageUrl(): Chainable<Element>;
    createTestByName(name: string): Chainable<Element>;
    submitAndMakeSureTestIsCreated(name: string): Chainable<Element>;
    createTestWithAuth(authMethod: string, keys: string[]): Chainable<string>;
    fillCreateFormBasicStep(name: string, description?: string): Chainable<Element>;
    setCreateFormUrl(method: string, url: string): Chainable<Element>;
    selectRunDetailMode(index: number): Chainable<Element>;
    interceptEditTestCall(): Chainable<Element>;
  }
}
