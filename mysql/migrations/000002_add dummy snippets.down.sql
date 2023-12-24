-- Remove the records added in the up migration
DELETE FROM snippets WHERE title = 'An old silent pond';
DELETE FROM snippets WHERE title = 'Over the wintry forest';
DELETE FROM snippets WHERE title = 'First autumn morning';
