describe('Exercise 3: Test Responsive Design', () => {
  
  it('should work on mobile viewport (375x667)', () => {
    cy.viewport(375, 667);
    cy.visit('/');
    
    // Check all key elements are visible
    cy.get('[data-testid="breed-selector"]').should('be.visible');
    cy.get('[data-testid="fetch-dog-button"]').should('be.visible');
    
    // Test functionality on mobile
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
  });

  it('should work on tablet viewport (768x1024)', () => {
    cy.viewport(768, 1024);
    cy.visit('/');
    
    // Check all key elements are visible
    cy.get('[data-testid="breed-selector"]').should('be.visible');
    cy.get('[data-testid="fetch-dog-button"]').should('be.visible');
    
    // Test functionality on tablet
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
  });

  it('should work on desktop viewport (1920x1080)', () => {
    cy.viewport(1920, 1080);
    cy.visit('/');
    
    // Check all key elements are visible
    cy.get('[data-testid="breed-selector"]').should('be.visible');
    cy.get('[data-testid="fetch-dog-button"]').should('be.visible');
    
    // Test functionality on desktop
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
  });
});