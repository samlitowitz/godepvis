package color

import (
	"fmt"
	"github.com/go-playground/colors"
	"github.com/mitchellh/mapstructure"
	"image/color"
	"reflect"
)

type Color struct {
	color.Color `mapstructure:"color"`
}

func (c Color) Hex() string {
	r, g, b, _ := c.RGBA()
	r = r >> 8
	g = g >> 8
	b = b >> 8
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func colorHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(Color{}) {
			return data, nil
		}
		in := data.(string)
		if in == "" {
			return nil, colors.ErrBadColor
		}
		hex, err := colors.ParseHEX(in)
		if err == nil {
			return Color{hex}, nil
		}
		rgb, err := colors.ParseRGB(in)
		if err == nil {
			return Color{rgb}, nil
		}
		rgba, err := colors.ParseRGBA(in)
		if err == nil {
			return Color{rgba}, nil
		}

		if err != nil {
			return nil, err
		}
		return nil, colors.ErrBadColor
	}
}
