# Test Plan for Tabs Component

## Overview
This document outlines the test plan for the Tabs component in the gh-portrait project.

## Test Cases

### TestNewTabs
Tests the initialization of tabs with various configurations.

Test cases:
- Multiple titles (2 or more)
  - Verify first tab is selected (Selected = true)
  - Verify other tabs are not selected (Selected = false)
  - Verify Current is set to 0
  - Verify all titles are correctly set

### TestTabsNext
Tests the next tab selection functionality.

Test cases:
- Normal movement
  - Verify Current value changes
  - Verify Selected flags are updated correctly
  - Verify previous tab is deselected
- Circular movement (from last to first tab)
  - Verify Current wraps to 0
  - Verify Selected flags are updated correctly
- Multiple movements
  - Verify correct tab selection after multiple Next calls

### TestTabsPrev
Tests the previous tab selection functionality.

Test cases:
- Normal movement
  - Verify Current value changes
  - Verify Selected flags are updated correctly
  - Verify previous tab is deselected
- Circular movement (from first to last tab)
  - Verify Current wraps to last index
  - Verify Selected flags are updated correctly
- Multiple movements
  - Verify correct tab selection after multiple Prev calls

### TestTabsView
Tests the rendering of tabs including styles.

Test cases:
- Active tab style verification
  - Verify bold style is applied
  - Verify correct color (86) is applied
- Inactive tab style verification
  - Verify correct color (241) is applied
- Tab gap verification
  - Verify gaps are inserted between tabs
  - Verify no gap after last tab
- Multiple tabs rendering
  - Verify correct horizontal joining of tabs
  - Verify overall layout matches expected format

## Implementation Priority
1. TestNewTabs (foundation)
2. TestTabsNext/TestTabsPrev (state management)
3. TestTabsView (rendering and styling)

## Notes
- All tests will use table-driven test pattern for consistency with existing tests
- Tabs component always has 2 or more tabs
- Style testing is included for complete coverage