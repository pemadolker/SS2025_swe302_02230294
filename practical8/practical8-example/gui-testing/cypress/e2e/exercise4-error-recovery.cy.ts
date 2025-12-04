describe('Exercise 4: Error Recovery Tests', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should recover when API fails then succeeds on retry', () => {
    let callCount = 0;
    
    // First call fails, second succeeds
    cy.intercept('GET', '/api/dogs', (req) => {
      callCount++;
      if (callCount === 1) {
        req.reply({
          statusCode: 500,
          body: { error: 'Server Error' }
        });
      } else {
        req.reply({
          statusCode: 200,
          body: {
            message: 'https://images.dog.ceo/breeds/husky/n02110185_1469.jpg',
            status: 'success'
          }
        });
      }
    }).as('getDog');
    
    // First attempt - should fail
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDog');
    
    // Error should be displayed
    cy.get('[data-testid="error-message"]')
      .should('be.visible');
    
    // Second attempt - should succeed
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDog');
    
    // Error should be gone
    cy.get('[data-testid="error-message"]')
      .should('not.exist');
    
    // Image should be displayed
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
  });

  it('should clear error when selecting different breed', () => {
    // Wait for breeds to load
    cy.get('[data-testid="breed-selector"] option', { timeout: 10000 })
      .should('have.length.greaterThan', 1);
    
    // Mock API failure
    cy.intercept('GET', '/api/dogs?breed=husky', {
      statusCode: 500,
      body: { error: 'Server Error' }
    }).as('getHuskyError');
    
    // Mock success for corgi
    cy.intercept('GET', '/api/dogs?breed=corgi', {
      statusCode: 200,
      body: {
        message: ['https://images.dog.ceo/breeds/corgi/n02113186_1234.jpg'],
        status: 'success'
      }
    }).as('getCorgiSuccess');
    
    // Select husky and fetch (fails)
    cy.get('[data-testid="breed-selector"]').select('husky');
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getHuskyError');
    
    // Error should appear
    cy.get('[data-testid="error-message"]')
      .should('be.visible');
    
    // Select corgi and fetch (succeeds)
    cy.get('[data-testid="breed-selector"]').select('corgi');
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getCorgiSuccess');
    
    // Error should clear
    cy.get('[data-testid="error-message"]')
      .should('not.exist');
    
    // Image should display
    cy.get('[data-testid="dog-image"]')
      .should('be.visible');
  });

  it('should handle multiple consecutive errors gracefully', () => {
    // Mock multiple failures
    cy.intercept('GET', '/api/dogs', {
      statusCode: 500,
      body: { error: 'Server Error' }
    }).as('getDogError');
    
    // Try 3 times
    for (let i = 0; i < 3; i++) {
      cy.get('[data-testid="fetch-dog-button"]').click();
      cy.wait('@getDogError');
      
      // Error should be visible each time
      cy.get('[data-testid="error-message"]')
        .should('be.visible');
      
      // Button should still be enabled for retry
      cy.get('[data-testid="fetch-dog-button"]')
        .should('not.be.disabled');
    }
  });

  it('should clear error message when starting new successful fetch', () => {
    let callCount = 0;
    
    cy.intercept('GET', '/api/dogs', (req) => {
      callCount++;
      if (callCount === 1) {
        req.reply({
          statusCode: 500,
          body: { error: 'Server Error' }
        });
      } else {
        req.reply({
          statusCode: 200,
          body: {
            message: 'https://images.dog.ceo/breeds/husky/n02110185_1469.jpg',
            status: 'success'
          }
        });
      }
    }).as('getDog');
    
    // First attempt - fail
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDog');
    cy.get('[data-testid="error-message"]').should('be.visible');
    
    // Second attempt - success
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Error should clear immediately or shortly after
    cy.get('[data-testid="error-message"]', { timeout: 2000 })
      .should('not.exist');
  });

  it('should recover from network timeout', () => {
    let callCount = 0;
    
    cy.intercept('GET', '/api/dogs', (req) => {
      callCount++;
      if (callCount === 1) {
        // First call - quick fail
        req.reply({
          statusCode: 408,
          body: { error: 'Request Timeout' }
        });
      } else {
        // Second call - quick success
        req.reply({
          statusCode: 200,
          body: {
            message: 'https://images.dog.ceo/breeds/husky/n02110185_1469.jpg',
            status: 'success'
          }
        });
      }
    }).as('getDog');
    
    // First attempt
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDog');
    
    // Wait a bit for error to show
    cy.wait(1000);
    
    // Button should be enabled for retry
    cy.get('[data-testid="fetch-dog-button"]')
      .should('not.be.disabled')
      .click();
    
    // Should eventually succeed
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
  });

  it('should maintain app state after error recovery', () => {
    // Wait for breeds
    cy.get('[data-testid="breed-selector"] option', { timeout: 10000 })
      .should('have.length.greaterThan', 1);
    
    let callCount = 0;
    
    cy.intercept('GET', '/api/dogs?breed=husky', (req) => {
      callCount++;
      if (callCount === 1) {
        req.reply({
          statusCode: 500,
          body: { error: 'Server Error' }
        });
      } else {
        req.reply({
          statusCode: 200,
          body: {
            message: ['https://images.dog.ceo/breeds/husky/n02110185_1469.jpg'],
            status: 'success'
          }
        });
      }
    }).as('getHusky');
    
    // Select husky
    cy.get('[data-testid="breed-selector"]').select('husky');
    
    // First attempt - fail
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getHusky');
    cy.get('[data-testid="error-message"]').should('be.visible');
    
    // Breed selection should still be husky
    cy.get('[data-testid="breed-selector"]')
      .should('have.value', 'husky');
    
    // Retry - success
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getHusky');
    
    // Should show husky image
    cy.get('[data-testid="dog-image"]')
      .should('be.visible')
      .invoke('attr', 'src')
      .should('include', 'husky');
    
    // Breed selection should still be husky
    cy.get('[data-testid="breed-selector"]')
      .should('have.value', 'husky');
  });

  it('should handle partial API failures (breeds fail, dogs succeed)', () => {
    // Mock breeds API failure
    cy.intercept('GET', '/api/dogs/breeds', {
      statusCode: 500,
      body: { error: 'Failed to fetch breeds' }
    }).as('getBreedsError');
    
    // Mock dogs API success
    cy.intercept('GET', '/api/dogs', {
      statusCode: 200,
      body: {
        message: 'https://images.dog.ceo/breeds/husky/n02110185_1469.jpg',
        status: 'success'
      }
    }).as('getDogsSuccess');
    
    // Reload to trigger breeds fetch
    cy.reload();
    cy.wait('@getBreedsError');
    
    // Even if breeds fail, should still be able to fetch random dogs
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDogsSuccess');
    
    cy.get('[data-testid="dog-image"]')
      .should('be.visible');
  });

  it('should provide feedback during error recovery', () => {
    let callCount = 0;
    
    cy.intercept('GET', '/api/dogs', (req) => {
      callCount++;
      if (callCount === 1) {
        req.reply({
          statusCode: 500,
          body: { error: 'Server Error' }
        });
      } else {
        req.reply({
          statusCode: 200,
          body: {
            message: 'https://images.dog.ceo/breeds/husky/n02110185_1469.jpg',
            status: 'success'
          }
        });
      }
    }).as('getDog');
    
    // First attempt - fail
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getDog');
    
    // Should show error
    cy.get('[data-testid="error-message"]')
      .should('be.visible')
      .and('contain.text', 'Failed');
    
    // Retry
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // FIXED: Check loading state properly
    cy.get('[data-testid="fetch-dog-button"]').then(($btn) => {
      const text = $btn.text();
      const isDisabled = $btn.prop('disabled');
      expect(text.includes('Loading') || isDisabled).to.be.true;
    });
    
    // Then success
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
    
    // Button should return to normal
    cy.get('[data-testid="fetch-dog-button"]')
      .should('not.be.disabled')
      .and('not.contain.text', 'Loading');
  });
});