package components

import (
	"strings"
	"testing"
)

func TestNewTabs(t *testing.T) {
	tests := []struct {
		name      string
		titles    []string
		wantTabs  int
		wantFirst bool
	}{
		{
			name:      "two tabs",
			titles:    []string{"Tab 1", "Tab 2"},
			wantTabs:  2,
			wantFirst: true,
		},
		{
			name:      "multiple tabs",
			titles:    []string{"Tab 1", "Tab 2", "Tab 3", "Tab 4"},
			wantTabs:  4,
			wantFirst: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := NewTabs(tt.titles)

			// Check number of tabs
			if got := len(tabs.Tabs); got != tt.wantTabs {
				t.Errorf("NewTabs() number of tabs = %v, want %v", got, tt.wantTabs)
			}

			// Check first tab is selected
			if got := tabs.Tabs[0].Selected; got != tt.wantFirst {
				t.Errorf("NewTabs() first tab selected = %v, want %v", got, tt.wantFirst)
			}

			// Check other tabs are not selected
			for i := 1; i < len(tabs.Tabs); i++ {
				if tabs.Tabs[i].Selected {
					t.Errorf("NewTabs() tab %d selected = true, want false", i)
				}
			}

			// Check Current is set to 0
			if got := tabs.Current; got != 0 {
				t.Errorf("NewTabs() Current = %v, want 0", got)
			}

			// Check all titles are set correctly
			for i, want := range tt.titles {
				if got := tabs.Tabs[i].Title; got != want {
					t.Errorf("NewTabs() tab %d title = %v, want %v", i, got, want)
				}
			}
		})
	}
}

func TestTabsNext(t *testing.T) {
	tests := []struct {
		name          string
		titles        []string
		moves         int
		wantCurrent   int
		wantSelected  int
		wantPrevious  int
		wantWraparound bool
	}{
		{
			name:         "move to next tab",
			titles:       []string{"Tab 1", "Tab 2", "Tab 3"},
			moves:        1,
			wantCurrent:  1,
			wantSelected: 1,
			wantPrevious: 0,
		},
		{
			name:          "wrap around to first tab",
			titles:        []string{"Tab 1", "Tab 2", "Tab 3"},
			moves:         3,
			wantCurrent:   0,
			wantSelected:  0,
			wantPrevious:  2,
			wantWraparound: true,
		},
		{
			name:         "multiple moves",
			titles:       []string{"Tab 1", "Tab 2", "Tab 3"},
			moves:        2,
			wantCurrent:  2,
			wantSelected: 2,
			wantPrevious: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := NewTabs(tt.titles)

			for range tt.moves {
				tabs.Next()
			}

			// Check Current value
			if got := tabs.Current; got != tt.wantCurrent {
				t.Errorf("Next() Current = %v, want %v", got, tt.wantCurrent)
			}

			// Check selected tab
			if got := tabs.Tabs[tt.wantSelected].Selected; !got {
				t.Errorf("Next() tab %d selected = false, want true", tt.wantSelected)
			}

			// Check previous tab is deselected
			if got := tabs.Tabs[tt.wantPrevious].Selected; got {
				t.Errorf("Next() previous tab %d selected = true, want false", tt.wantPrevious)
			}
		})
	}
}

func TestTabsPrev(t *testing.T) {
	tests := []struct {
		name          string
		titles        []string
		moves         int
		wantCurrent   int
		wantSelected  int
		wantPrevious  int
		wantWraparound bool
	}{
		{
			name:         "move to previous tab",
			titles:       []string{"Tab 1", "Tab 2", "Tab 3"},
			moves:        1,
			wantCurrent:  2,
			wantSelected: 2,
			wantPrevious: 0,
		},
		{
			name:          "wrap around to last tab",
			titles:        []string{"Tab 1", "Tab 2", "Tab 3"},
			moves:         3,
			wantCurrent:   0,
			wantSelected:  0,
			wantPrevious:  1,
			wantWraparound: true,
		},
		{
			name:         "multiple moves",
			titles:       []string{"Tab 1", "Tab 2", "Tab 3"},
			moves:        2,
			wantCurrent:  1,
			wantSelected: 1,
			wantPrevious: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := NewTabs(tt.titles)

			for range tt.moves {
				tabs.Prev()
			}

			// Check Current value
			if got := tabs.Current; got != tt.wantCurrent {
				t.Errorf("Prev() Current = %v, want %v", got, tt.wantCurrent)
			}

			// Check selected tab
			if got := tabs.Tabs[tt.wantSelected].Selected; !got {
				t.Errorf("Prev() tab %d selected = false, want true", tt.wantSelected)
			}

			// Check previous tab is deselected
			if got := tabs.Tabs[tt.wantPrevious].Selected; got {
				t.Errorf("Prev() previous tab %d selected = true, want false", tt.wantPrevious)
			}
		})
	}
}

func TestTabsView(t *testing.T) {
	tests := []struct {
		name     string
		titles   []string
		current  int
		contains []string
	}{
		{
			name:    "active and inactive tabs",
			titles:  []string{"Tab 1", "Tab 2"},
			current: 0,
			contains: []string{
				"Tab 1",
				"Tab 2",
			},
		},
		{
			name:    "multiple tabs with gap",
			titles:  []string{"Tab 1", "Tab 2", "Tab 3"},
			current: 1,
			contains: []string{
				"Tab 1",
				"Tab 2",
				"Tab 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := NewTabs(tt.titles)
			
			// Move to specified current tab
			for range tt.current {
				tabs.Next()
			}

			result := tabs.View()

			// Check that all titles are present in the rendered output
			for _, title := range tt.contains {
				if !strings.Contains(result, title) {
					t.Errorf("View() result does not contain %q", title)
				}
			}

			// Check that there's a gap between tabs (but not after the last tab)
			gapCount := strings.Count(result, tabGap)
			expectedGaps := len(tt.titles) - 1
			if gapCount != expectedGaps {
				t.Errorf("View() gap count = %v, want %v", gapCount, expectedGaps)
			}

			// Check active tab style
			activeTab := tt.titles[tt.current]
			styledActiveTab := activeTabStyle.Render(activeTab)
			if !strings.Contains(result, styledActiveTab) {
				t.Errorf("View() active tab style not applied correctly for %q", activeTab)
			}

			// Check inactive tab style
			for i, title := range tt.titles {
				if i != tt.current {
					styledInactiveTab := inactiveTabStyle.Render(title)
					if !strings.Contains(result, styledInactiveTab) {
						t.Errorf("View() inactive tab style not applied correctly for %q", title)
					}
				}
			}
		})
	}
}
