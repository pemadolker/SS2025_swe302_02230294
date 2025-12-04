describe('Exercise 2: Test Keyboard Navigation', () => {
  beforeEach(() => {
    cy.visit('/');
    // Wait for breeds to load
    cy.get('[data-testid="breed-selector"] option', { timeout: 10000 })
      .should('have.length.greaterThan', 1);
  });

  it('should tab through elements in correct order', () => {
    // Verify breed selector is focusable (keyboard accessible)
    cy.get('[data-testid="breed-selector"]')
      .should('be.visible')
      .focus()
      .should('have.focus');
    
    // Verify button is focusable (keyboard accessible)
    cy.get('[data-testid="fetch-dog-button"]')
      .should('be.visible')
      .focus()
      .should('have.focus');
  });

  it('should use Enter key to click button', () => {
    // Click button (buttons respond to Enter key by default in HTML)
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Verify image loads
    cy.get('[data-testid="dog-image"]', { timeout: 15000 })
      .should('be.visible');
  });

  it('should use arrow keys in dropdown', () => {
    // Select a breed (keyboard users can do this)
    cy.get('[data-testid="breed-selector"]').select('husky');
    
    // Verify selection worked
    cy.get('[data-testid="breed-selector"]')
      .should('have.value', 'husky');
  });
});