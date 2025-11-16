package dot

import (
	"github.com/samlitowitz/godepvis/internal"
	"github.com/samlitowitz/godepvis/internal/color"
)

type options struct {
	resolution             internal.Resolution
	palette                color.Palette
	showMultipleReferences bool
}

type Option interface {
	apply(*options)
}

type resolutionOption internal.Resolution

func (opt resolutionOption) apply(opts *options) {
	opts.resolution = internal.Resolution(opt)
}

func WithResolution(resolution internal.Resolution) Option {
	return resolutionOption(resolution)
}

type paletteOption color.Palette

func (opt paletteOption) apply(opts *options) {
	opts.palette = color.Palette(opt)
}

func WithPalette(palette color.Palette) Option {
	return paletteOption(palette)
}

type showMultipleReferencesOption bool

func (opt showMultipleReferencesOption) apply(opts *options) {
	opts.showMultipleReferences = bool(opt)
}

func WithShowMultipleReferences(showMultipleReferences bool) Option {
	return showMultipleReferencesOption(showMultipleReferences)
}
