/// <reference types="cypress" />

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Cypress {
    interface Chainable {
      /**
       * Custom command to fetch a dog image
       * @example cy.fetchDog()
       */
      fetchDog(): Chainable<void>;

      /**
       * Custom command to select a breed and fetch dog
       * @param breed - The breed name to select
       * @example cy.selectBreedAndFetch('husky')
       */
      selectBreedAndFetch(breed: string): Chainable<void>;

      /**
       * Custom command to wait for dog image to load
       * @example cy.waitForDogImage()
       */
      waitForDogImage(): Chainable<JQuery<HTMLElement>>;

      /**
       * Custom command to check if error is displayed
       * @example cy.checkError('Failed to load')
       */
      checkError(message: string): Chainable<void>;
    }
  }
}

// Fetch dog image command
Cypress.Commands.add('fetchDog', () => {
  cy.get('[data-testid="fetch-dog-button"]').click();
});

// Select breed and fetch command
Cypress.Commands.add('selectBreedAndFetch', (breed: string) => {
  cy.get('[data-testid="breed-selector"]').select(breed);
  cy.get('[data-testid="fetch-dog-button"]').click();
});

// Wait for dog image to load
Cypress.Commands.add('waitForDogImage', () => {
  return cy.get('[data-testid="dog-image"]', { timeout: 10000 })
    .should('be.visible');
});

// Check error message
Cypress.Commands.add('checkError', (message: string) => {
  cy.get('[data-testid="error-message"]')
    .should('be.visible')
    .and('contain.text', message);
});

export {};
// ***********************************************
// This example commands.ts shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
//
// declare global {
//   namespace Cypress {
//     interface Chainable {
//       login(email: string, password: string): Chainable<void>
//       drag(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
//       dismiss(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
//       visit(originalFn: CommandOriginalFn, url: string, options: Partial<VisitOptions>): Chainable<Element>
//     }
//   }
// }