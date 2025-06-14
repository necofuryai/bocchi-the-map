# BDD Test Specifications for Bocchi The Map

## Feature: User Authentication
### Scenario: User can sign in with Google OAuth
**Given** the user is on the homepage  
**When** the user clicks the sign-in button  
**Then** the user should be redirected to `Google OAuth`  
**And** after successful authentication, the user should be signed in  
**And** the user's profile should be displayed in the header  

### Scenario: User can sign out
**Given** the user is signed in  
**When** the user clicks the sign-out button  
**Then** the user should be signed out  
**And** the sign-in button should be visible again  

## Feature: Map Display and Interaction
### Scenario: Map loads correctly on homepage
**Given** the user visits the homepage  
**When** the page loads  
**Then** the map should be displayed  
**And** the map should show the default location  
**And** the map controls should be visible  

### Scenario: User can interact with map controls
**Given** the map is displayed  
**When** the user uses navigation controls  
**Then** the map should zoom in/out accordingly  
**And** the map position should change as expected  

### Scenario: POI features are displayed on map
**Given** the map is loaded  
**When** POI data is available  
**Then** POI markers should be displayed on the map  
**And** clicking a POI marker should show details  

### Scenario: POI API returns 4xx client error
**Given** the map is loaded  
**When** the POI API returns a 4xx client error (e.g., 404, 400)  
**Then** an error message should be displayed to the user  
**And** the error UI should explain the issue clearly  
**And** a retry button should be available  

### Scenario: POI API returns 5xx server error
**Given** the map is loaded  
**When** the POI API returns a 5xx server error (e.g., 500, 503)  
**Then** a server error message should be displayed to the user  
**And** the error UI should indicate it's a temporary issue  
**And** a retry button should be available  
**And** an automatic retry should be attempted after a short delay  

### Scenario: Network failure during POI data fetch
**Given** the map is loaded  
**When** a network failure occurs during POI data fetching  
**Then** a network error message should be displayed to the user  
**And** the error UI should indicate connectivity issues  
**And** a retry button should be available  
**And** the app should retry automatically when connectivity is restored  

### Scenario: User can retry after POI API failure
**Given** the POI API has failed with any error type  
**When** the user clicks the retry button  
**Then** the error UI should show a loading state  
**And** a new API request should be made  
**And** on success, POI markers should be displayed normally  
**And** on failure, the appropriate error message should be shown again  

### Scenario: POI data takes too long to load (timeout)
**Given** the map is loaded  
**When** the POI API request times out  
**Then** a timeout error message should be displayed  
**And** the error UI should suggest checking connectivity  
**And** a retry button should be available  

## Feature: Theme Switching
### Scenario: User can switch to dark mode
**Given** the user is on the homepage  
**And** the current theme is light mode  
**When** the user clicks the theme toggle button  
**Then** the application should switch to dark mode  
**And** the theme preference should be saved  

### Scenario: User can switch to light mode
**Given** the user is on the homepage  
**And** the current theme is dark mode  
**When** the user clicks the theme toggle button  
**Then** the application should switch to light mode  
**And** the theme preference should be saved  

## Feature: Navigation and Routing
### Scenario: User can navigate through the application
**Given** the user is on the homepage  
**When** the user clicks on navigation menu items  
**Then** the user should be navigated to the correct page  
**And** the URL should update accordingly  
**And** the page content should load correctly  

## Feature: Responsive Design
### Scenario: Application works on mobile devices
**Given** the user accesses the site on a mobile device  
**When** the page loads  
**Then** the layout should be mobile-responsive  
**And** the map should fit the mobile screen  
**And** navigation should be accessible on mobile  

### Scenario: Application works on desktop
**Given** the user accesses the site on a desktop  
**When** the page loads  
**Then** the layout should use the full desktop width  
**And** all components should be properly positioned  
**And** desktop-specific features should be available  

### Scenario: Application maintains layout stability in landscape mode
**Given** the user accesses the site in landscape orientation on mobile or tablet  
**When** the page loads  
**Then** the layout should not break in landscape orientation  
**And** the map should adapt properly to the wider viewport  
**And** navigation elements should remain accessible and properly positioned  

## Feature: Component Behavior
### Scenario: Header component displays correctly
**Given** the user visits any page  
**When** the page loads  
**Then** the header should be visible  
**And** the logo/title should be displayed  
**And** navigation elements should be present  

### Scenario: Map status indicator works
**Given** the map is loading  
**When** the map data is being fetched  
**Then** a loading indicator should be shown  
**And** when loading completes, the indicator should disappear  
**And** if there's an error, an error message should be shown  

## Feature: Error Handling
### Scenario: Application handles map loading errors gracefully
**Given** the user visits the homepage  
**When** the map fails to load  
**Then** an error message should be displayed  
**And** the user should be able to retry loading the map  

### Scenario: Application handles authentication errors
**Given** the user attempts to sign in  
**When** the authentication process fails  
**Then** an error message should be displayed  
**And** the user should remain on the sign-in flow

### Scenario: Retry button functions correctly after map loading error
**Given** the map fails to load and displays an error message  
**When** the user clicks the retry button  
**Then** the map should attempt to reload  
**And** the error message should be cleared  
**And** a loading indicator should be shown during retry  

### Scenario: Retry button handles multiple rapid clicks gracefully
**Given** the map fails to load and displays a retry button  
**When** the user rapidly clicks the retry button multiple times  
**Then** only one retry attempt should be initiated  
**And** subsequent clicks should be ignored until the current retry completes  
**And** the button should be disabled during the retry process  

### Scenario: Retry button recovers from failed retry attempts
**Given** the user clicks retry and the retry attempt also fails  
**When** the retry failure occurs  
**Then** the error message should be updated to reflect the retry failure  
**And** the retry button should remain available for another attempt  
**And** the system should not enter an infinite retry loop  