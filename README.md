## Server-side Rendered Application

### Architecture

- Modular Monolith

### Tech

- go for backend
- htmx for ajax requests
- tailwind for css
- Mysql for persistance

### Frameworks

- no frameworks, I use go stdlib only.

## Use Case: Snippet Creation and Timeline Update with Flash Notification

**User Scenario:**

1. **User Intent:**

   - The user wants to create a new code snippet in a code snippet management application.

2. **User Action:**

   - The user navigates to the application and selects the option to create a new code snippet.

3. **Snippet Creation:**

   - The application opens a new snippet creation interface where the user can input code, add a title, and provide a description for the snippet.

4. **Save and Publish:**

   - After completing the snippet, the user clicks on the "Save" or "Publish" button to save the snippet and make it visible to others.

5. **Timeline Update:**

   - The application automatically updates the user's timeline with the newly created snippet, showcasing a preview or relevant details of the code.

6. **Flash Notification:**

   - Simultaneously, a subtle and visually pleasing flash notification appears on the user interface, indicating that the snippet has been successfully created and added to the timeline.

7. **Duration:**

   - The flash notification remains visible for 3 seconds, providing a brief confirmation to the user without causing distraction.

8. **Notification Content:**

   - The flash notification may include a short message like "Snippet Created Successfully!" or any other informative message.

9. **Disappearance:**
   - After 3 seconds, the flash notification smoothly disappears, leaving the timeline as the focal point for the user.

## TODO

- Prevent Access To Partial Page Components

- Feature: Edit Snippets Form
