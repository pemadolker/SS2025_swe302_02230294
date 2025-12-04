describe('Exercise 1: Image Loading Tests', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should verify image has correct src attribute', () => {
    // Click fetch button to load an image
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Wait for image to appear and verify src attribute
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .and('have.attr', 'src')
      // FIXED: Next.js Image component transforms URLs
      .and('include', 'images.dog.ceo');
    
    // Verify src contains image file reference
    cy.get('[data-testid="dog-image"]')
      .invoke('attr', 'src')
      .should('match', /\.(jpg|jpeg|png|gif)/i);
  });

  it('should verify image loads successfully (not broken)', () => {
    // Fetch dog image
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Wait for image and check if it loaded successfully
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .and(($img) => {
        // naturalWidth > 0 means image loaded successfully
        expect($img[0].naturalWidth).to.be.greaterThan(0);
        expect($img[0].naturalHeight).to.be.greaterThan(0);
      });
  });

  it('should verify image is visible in viewport', () => {
    // Fetch dog image
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    // Verify image is visible
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible');
    
    // FIXED: Check if image is within viewport bounds manually
    cy.get('[data-testid="dog-image"]').then(($img) => {
      const rect = $img[0].getBoundingClientRect();
      const windowHeight = Cypress.config('viewportHeight');
      const windowWidth = Cypress.config('viewportWidth');
      
      expect(rect.top).to.be.lessThan(windowHeight);
      expect(rect.bottom).to.be.greaterThan(0);
      expect(rect.left).to.be.lessThan(windowWidth);
      expect(rect.right).to.be.greaterThan(0);
    });
  });

  it('should load multiple images successfully', () => {
    const imageUrls = [];
    
    // Load first image
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .and(($img) => {
        expect($img[0].naturalWidth).to.be.greaterThan(0);
      })
      .invoke('attr', 'src')
      .then((src) => {
        imageUrls.push(src);
      });
    
    // Load second image
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .and(($img) => {
        expect($img[0].naturalWidth).to.be.greaterThan(0);
      });
  });

  it('should handle image loading errors gracefully', () => {
    // Mock a broken image URL
    cy.intercept('GET', '/api/dogs', {
      statusCode: 200,
      body: {
        message: 'https://invalid-url.example.com/broken-image.jpg',
        status: 'success',
      },
    }).as('getBrokenImage');
    
    cy.get('[data-testid="fetch-dog-button"]').click();
    cy.wait('@getBrokenImage');
    
    // The app should handle this gracefully
    // Either show error or have a fallback
  });

  it('should display image with proper aspect ratio', () => {
    cy.get('[data-testid="fetch-dog-button"]').click();
    
    cy.get('[data-testid="dog-image"]', { timeout: 10000 })
      .should('be.visible')
      .and(($img) => {
        const width = $img[0].naturalWidth;
        const height = $img[0].naturalHeight;
        
        // Images should have reasonable dimensions
        expect(width).to.be.greaterThan(100);
        expect(height).to.be.greaterThan(100);
        
        // Aspect ratio should be reasonable (not too extreme)
        const aspectRatio = width / height;
        expect(aspectRatio).to.be.within(0.5, 2);
      });
  });
});