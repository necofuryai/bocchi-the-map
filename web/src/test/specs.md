# BDD Test Specifications for Bocchi The Map

## Feature: User Authentication
### Scenario: User can sign in with Google OAuth
**Given** the user is on the homepage  
**When** the user clicks the sign-in button  
**Then** the user should be redirected to Google OAuth  
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