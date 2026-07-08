package TUI

import "fmt"

type RGB struct {
	R, G, B uint8
}

const Reset = "\033[0m"

var Catppuccin = map[string]RGB{
	"Rosewater": {245, 224, 220},
	"Flamingo":  {242, 205, 205},
	"Pink":      {245, 194, 231},
	"Mauve":     {203, 166, 247},
	"Red":       {243, 139, 168},
	"Maroon":    {235, 160, 172},
	"Peach":     {250, 179, 135},
	"Yellow":    {249, 226, 175},
	"Green":     {166, 227, 161},
	"Teal":      {148, 226, 213},
	"Sky":       {137, 220, 235},
	"Sapphire":  {116, 199, 236},
	"Blue":      {137, 180, 250},
	"Lavender":  {180, 190, 254},

	"Text":     {205, 214, 244},
	"Subtext1": {186, 194, 222},
	"Subtext0": {166, 173, 200},
	"Overlay2": {147, 153, 178},
	"Overlay1": {127, 132, 156},
	"Overlay0": {108, 112, 134},
	"Surface2": {88, 91, 112},
	"Surface1": {69, 71, 90},
	"Surface0": {49, 50, 68},
	"Base":     {30, 30, 46},
	"Mantle":   {24, 24, 37},
	"Crust":    {17, 17, 27},
}

func (c RGB) Foreground() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

func (c RGB) Background() string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", c.R, c.G, c.B)
}
