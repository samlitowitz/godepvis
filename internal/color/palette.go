package color

import (
	"fmt"
	"github.com/spf13/viper"
	"image/color"
)

type HalfPalette struct {
	PackageName       Color `mapstructure:"packagename"`
	PackageBackground Color `mapstructure:"packagebackground"`
	FileName          Color `mapstructure:"filename"`
	FileBackground    Color `mapstructure:"filebackground"`
	ImportArrow       Color `mapstructure:"importarrow"`
}

type Palette struct {
	Base  *HalfPalette `mapstructure:"base"`
	Cycle *HalfPalette `mapstructure:"cycle"`
}

var (
	DefaultPalette = &Palette{
		Base: &HalfPalette{
			PackageName: Color{
				Color: &color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				},
			},
			PackageBackground: Color{
				Color: &color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 0,
				},
			},
			FileName: Color{
				Color: &color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				},
			},
			FileBackground: Color{
				Color: &color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 0,
				},
			},
			ImportArrow: Color{
				Color: &color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				},
			},
		},
		Cycle: &HalfPalette{
			PackageName: Color{
				Color: &color.RGBA{
					R: 255,
					G: 0,
					B: 0,
					A: 0,
				},
			},
			PackageBackground: Color{
				Color: &color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 0,
				},
			},
			FileName: Color{
				Color: &color.RGBA{
					R: 255,
					G: 0,
					B: 0,
					A: 0,
				},
			},
			FileBackground: Color{
				Color: &color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 0,
				},
			},
			ImportArrow: Color{
				Color: &color.RGBA{
					R: 255,
					G: 0,
					B: 0,
					A: 0,
				},
			},
		},
	}
	InvertedDefaultPalette = &Palette{
		Base: &HalfPalette{
			PackageName: Color{
				Color: &color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 0,
				},
			},
			PackageBackground: Color{
				Color: &color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				},
			},
			FileName: Color{
				Color: &color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 0,
				},
			},
			FileBackground: Color{
				Color: &color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				},
			},
			ImportArrow: Color{
				Color: &color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 0,
				},
			},
		},
		Cycle: &HalfPalette{
			PackageName: Color{
				Color: &color.RGBA{
					R: 0,
					G: 255,
					B: 255,
					A: 0,
				},
			},
			PackageBackground: Color{
				Color: &color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				},
			},
			FileName: Color{
				Color: &color.RGBA{
					R: 0,
					G: 255,
					B: 255,
					A: 0,
				},
			},
			FileBackground: Color{
				Color: &color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				},
			},
			ImportArrow: Color{
				Color: &color.RGBA{
					R: 0,
					G: 255,
					B: 255,
					A: 0,
				},
			},
		},
	}
)

func GetPaletteFromFile(file string) (*Palette, error) {
	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load palette: %w", err)
	}

	p := DefaultPalette
	err = v.Unmarshal(&p, viper.DecodeHook(colorHookFunc()))
	if err != nil {
		return nil, fmt.Errorf("failed to load palette: %w", err)
	}
	return p, nil
}
