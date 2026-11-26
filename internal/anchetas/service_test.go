package anchetas

import "testing"

func TestAvailable(t *testing.T) {
	stock := map[string]int{"vino": 5, "queso": 9, "pan": 0}

	tests := []struct {
		name  string
		items []Item
		want  int
	}{
		{"sin items", nil, 0},
		{"un producto", []Item{{ProductID: "vino", Quantity: 2}}, 2},
		{"el producto más escaso manda", []Item{{ProductID: "vino", Quantity: 1}, {ProductID: "queso", Quantity: 3}}, 3},
		{"producto agotado", []Item{{ProductID: "vino", Quantity: 1}, {ProductID: "pan", Quantity: 1}}, 0},
		{"producto que ya no existe", []Item{{ProductID: "borrado", Quantity: 1}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := available(tt.items, stock); got != tt.want {
				t.Errorf("available() = %d, want %d", got, tt.want)
			}
		})
	}
}
