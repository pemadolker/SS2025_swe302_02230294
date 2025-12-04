# Practical 8: GUI Testing with Cypress - Final Report

**Student Name:** Pema Dolker   
**Student ID:** 02230294    
**Module:** Software Testing and Qaulity Assurance 

---

## Introduction

This practical focused on implementing automated GUI testing using Cypress for a Next.js Dog Image Browser application. The goal was to ensure the application works correctly from a user's perspective by testing various aspects including image loading, keyboard navigation, responsive design, error handling, and accessibility.

### What is GUI Testing?

GUI (Graphical User Interface) testing verifies that an application's user interface works as expected. Instead of testing individual functions in code, GUI testing simulates real user interactions like clicking buttons, typing in forms, and navigating with the keyboard.

### Why Cypress?

We used Cypress for this practical because:
- It's designed specifically for modern web applications
- Tests run in real browsers, showing exactly what users see
- It automatically waits for elements to appear (no manual delays needed)
- It provides excellent debugging tools with time-travel and screenshots
- It's written in JavaScript/TypeScript, which I'm already familiar with

---

## Setup and Configuration

### Installation Process

First, I installed all necessary dependencies:

```bash
# Install Cypress
pnpm add -D cypress

# Install helper tools
pnpm add -D start-server-and-test

# Install accessibility testing tools
pnpm add -D cypress-axe axe-core
```

### Configuration Files Created

#### 1. Cypress Configuration (`cypress.config.ts`)

Created the main configuration file that tells Cypress how to run tests:

```typescript
import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    setupNodeEvents() {},
    specPattern: 'cypress/e2e/**/*.cy.{js,jsx,ts,tsx}',
    supportFile: 'cypress/support/e2e.ts',
    viewportWidth: 1280,
    viewportHeight: 720,
    video: true,
    screenshotOnRunFailure: true,
    defaultCommandTimeout: 10000,
  },
});
```

**What this does:**
- `baseUrl`: Points to my local development server
- `video: true`: Records videos of test runs
- `screenshotOnRunFailure: true`: Takes screenshots when tests fail
- `defaultCommandTimeout: 10000`: Waits up to 10 seconds for elements

#### 2. Support File (`cypress/support/e2e.ts`)

This file loads before every test:

```typescript
import './commands';
import 'cypress-axe';
```

#### 3. Package.json Scripts

Added commands to easily run tests:

```json
{
  "scripts": {
    "cypress:open": "cypress open",
    "cypress:run": "cypress run",
    "test:e2e": "start-server-and-test dev http://localhost:3000 cypress:run",
    "test:e2e:open": "start-server-and-test dev http://localhost:3000 cypress:open"
  }
}
```

### Adding Test IDs to Components

To make testing reliable, I added `data-testid` attributes to key elements in my Dog Browser app:

```tsx
<h1 data-testid="page-title">Dog Image Browser</h1>

<select data-testid="breed-selector">
  {/* breed options */}
</select>

<button data-testid="fetch-dog-button">
  Get Random Dog
</button>

<img data-testid="dog-image" src={dogImage} alt="Random dog" />

{error && (
  <div data-testid="error-message" role="alert">
    {error}
  </div>
)}
```

**Why use data-testid?**
- CSS classes might change when styling updates
- Element positions might change when layout changes  
- Text content might change or be translated
- `data-testid` is specifically for testing and won't change

---

## Exercise Implementations

### Exercise 1: Test Image Loading

**Goal:** Verify that dog images load correctly and display properly.

**Tests Implemented:**
1. **Verify image has correct src attribute** - Checks that the image URL points to the Dog API
2. **Verify image loads successfully** - Uses `naturalWidth` to ensure image isn't broken
3. **Verify image is visible in viewport** - Confirms image appears on screen

![alt text](image-1.png)

**Code Example:**
```typescript
it('should verify image loads successfully (not broken)', () => {
  cy.get('[data-testid="fetch-dog-button"]').click();
  
  cy.get('[data-testid="dog-image"]', { timeout: 10000 })
    .should('be.visible')
    .and(($img) => {
      expect($img[0].naturalWidth).to.be.greaterThan(0);
    });
});
```

**What I learned:**
- How to check if an image actually loaded (not just appeared)
- The difference between an element existing and being visible
- How to handle async operations with timeouts

---

### Exercise 2: Test Keyboard Navigation

**Goal:** Ensure users who rely on keyboards (including people with disabilities) can use the app.

**Tests Implemented:**
1. **Tab through elements in correct order** - Verifies logical tab sequence
2. **Use Enter key to click button** - Tests keyboard activation
3. **Use arrow keys in dropdown** - Validates keyboard navigation in select element

![alt text](image-2.png)

**Code Example:**
```typescript
it('should use Enter key to click button', () => {
  cy.get('[data-testid="fetch-dog-button"]').focus();
  cy.focused().type('{enter}');
  
  cy.get('[data-testid="dog-image"]', { timeout: 10000 })
    .should('be.visible');
});
```

**What I learned:**
- Not everyone uses a mouse - keyboard accessibility is crucial
- Tab order matters for usability
- Cypress can simulate keyboard events realistically

---

### Exercise 3: Test Responsive Design

**Goal:** Verify the app works on mobile phones, tablets, and desktop computers.

**Tests Implemented:**
1. **Mobile viewport (375x667)** - Tests on phone-sized screens
2. **Tablet viewport (768x1024)** - Tests on tablet-sized screens
3. **Desktop viewport (1920x1080)** - Tests on computer monitors

![alt text](image-3.png)


**Code Example:**
```typescript
it('should work on mobile viewport (375x667)', () => {
  cy.viewport(375, 667);
  cy.visit('/');
  
  cy.get('[data-testid="breed-selector"]').should('be.visible');
  cy.get('[data-testid="fetch-dog-button"]').should('be.visible');
  
  cy.get('[data-testid="fetch-dog-button"]').click();
  cy.get('[data-testid="dog-image"]', { timeout: 10000 })
    .should('be.visible');
});
```

**What I learned:**
- Same app needs to work on screens of vastly different sizes
- Mobile users are just as important as desktop users
- Responsive design isn't just about layout - functionality must work too

---

### Exercise 4: Test Error Recovery

**Goal:** Verify users can recover from errors gracefully.

**Tests Implemented:**
1. **API fails then succeeds on retry** - Tests that users can try again after an error
2. **Error clears when selecting different breed** - Validates error state management

![alt text](image-4.png)

**Code Example:**
```typescript
it('should recover when API fails then succeeds on retry', () => {
  let callCount = 0;
  
  cy.intercept('GET', '/api/dogs', (req) => {
    callCount++;
    if (callCount === 1) {
      req.reply({ statusCode: 500, body: { error: 'Server Error' } });
    } else {
      req.reply({ 
        statusCode: 200, 
        body: { message: 'https://images.dog.ceo/...', status: 'success' }
      });
    }
  }).as('getDog');
  
  // First click - fails
  cy.get('[data-testid="fetch-dog-button"]').click();
  cy.wait('@getDog');
  cy.get('[data-testid="error-message"]').should('be.visible');
  
  // Second click - succeeds
  cy.get('[data-testid="fetch-dog-button"]').click();
  cy.wait('@getDog');
  cy.get('[data-testid="error-message"]').should('not.exist');
  cy.get('[data-testid="dog-image"]').should('be.visible');
});
```

**What I learned:**
- Apps should never trap users in error states
- API calls can fail - apps must handle this gracefully
- `cy.intercept()` lets me test error scenarios without breaking the API

---

### Exercise 5: Accessibility Testing

**Goal:** Ensure the app is usable by people with disabilities.

**Tests Implemented:**
1. **No accessibility violations** - Automated accessibility checks
2. **Proper focus indicators** - Visual feedback when tabbing
3. **Keyboard navigation** - Full keyboard support

![alt text](image-5.png)

**Code Example:**
```typescript
describe('Exercise 5: Accessibility Testing', () => {
  beforeEach(() => {
    cy.visit('/');
    cy.injectAxe();
  });

  it('should have no detectable accessibility violations', () => {
    cy.checkA11y();
  });

  it('should have proper focus indicators', () => {
    cy.get('[data-testid="fetch-dog-button"]')
      .focus()
      .should('have.focus');
  });
});
```

**What I learned:**
- Accessibility isn't optional - it's a legal requirement in many places
- Automated tools catch many (but not all) accessibility issues
- Small changes like adding `alt` text make huge differences

---

## Test Results

### Summary Statistics

![alt text](image.png)


## Key Learnings

### Technical Skills Gained

1. **Cypress Test Framework**
   - Writing E2E tests with Cypress
   - Using selectors and assertions
   - Handling async operations
   - Mocking API responses with `cy.intercept()`

2. **Testing Best Practices**
   - Using `data-testid` for reliable selectors
   - Testing user behavior, not implementation
   - Writing independent, isolated tests
   - Proper error handling in tests

3. **Accessibility Awareness**
   - Importance of keyboard navigation
   - Screen reader compatibility
   - WCAG standards basics
   - Using automated tools like axe-core

4. **Responsive Design Testing**
   - Testing multiple viewport sizes
   - Ensuring functionality across devices
   - Mobile-first considerations

### Conceptual Understanding

1. **GUI Testing Philosophy**
   - Test what users see and do, not internal code
   - Real browser testing catches real bugs
   - Automated tests provide regression safety

2. **Quality Assurance**
   - Testing is an investment, not a cost
   - Early bug detection saves time later
   - Tests document expected behavior

3. **User Experience**
   - Error recovery is critical
   - Accessibility benefits everyone
   - Responsive design is non-negotiable

---

## Screenshots

### 1. All Tests Passing (Interactive Mode)

![alt text](image-8.png)

**Description:** This shows all tests running successfully in Cypress's interactive mode. The green checkmarks indicate passing tests, and you can see the execution time for each test.

---

### 2. Headless Mode Results (Terminal)

![alt text](image-7.png)



### 3. Test Video Example

<video controls src="practical8-example/gui-testing/cypress/videos/exercise3-responsive.cy.ts.mp4" title="Title"></video>

**Description:** This video shows Exercise 3 running in real-time. 

---


## Conclusion

This practical successfully demonstrated the implementation of comprehensive GUI testing using Cypress for a Next.js application. I completed all five exercises, achieving 100% test pass rate across 15 tests covering:

- Image loading verification
- Keyboard navigation
- Responsive design across three viewport sizes
- Error recovery scenarios
- Accessibility compliance

### Project Outcomes

 **All exercises completed** - 5/5 exercises implemented  
 **All tests passing** -  54/54 tests successful  
 **Accessibility compliant** - No violations detected  
 **Cross-device tested** - Mobile, tablet, and desktop  
**Error handling verified** - Graceful failure and recovery




---

## Appendix

### Commands Reference

```bash
# Run tests interactively
pnpm cypress:open

# Run tests headlessly
pnpm cypress:run

# Run specific test file
pnpm cypress:run --spec "cypress/e2e/exercise1-image-loading.cy.ts"

# Start app and run tests
pnpm test:e2e
```

### Resources Used

- [Cypress Documentation](https://docs.cypress.io)
- [cypress-axe Documentation](https://github.com/component-driven/cypress-axe)
- [WCAG Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [Next.js Testing Documentation](https://nextjs.org/docs/testing)

---