package launcher

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	iconColumns = 4
	iconRows    = 4
)

var iconCache = map[string][]string{}

func itemDisplayLines(item item, cursor string) []string {
	icon := renderItemIcon(item)
	if len(icon) < 2 {
		return []string{cursor + item.title}
	}
	return []string{
		cursor + icon[0] + "  " + item.title,
		"  " + icon[1],
	}
}

func renderItemIcon(item item) []string {
	if item.typ == commandItem {
		return []string{"⌘", ""}
	}
	if item.typ == fileItem {
		return []string{"□", ""}
	}
	if item.iconPath == "" {
		return []string{"◌", ""}
	}
	if icon, ok := iconCache[item.iconPath]; ok {
		return icon
	}

	icon, err := renderIconPath(item.iconPath)
	if err != nil {
		icon = []string{"◌", ""}
	}
	iconCache[item.iconPath] = icon
	return icon
}

func renderIconPath(iconPath string) ([]string, error) {
	pngPath, err := cachedIconPNG(iconPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(pngPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return renderANSIIcon(img), nil
}

func cachedIconPNG(iconPath string) (string, error) {
	cacheDir := filepath.Join(os.TempDir(), "launcher_icon_cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	sum := sha1.Sum([]byte(iconPath))
	outPath := filepath.Join(cacheDir, hex.EncodeToString(sum[:])+".png")
	if _, err := os.Stat(outPath); err == nil {
		return outPath, nil
	}

	if err := exec.Command("sips", "-s", "format", "png", iconPath, "--out", outPath).Run(); err != nil {
		return "", err
	}
	return outPath, nil
}

func renderANSIIcon(img image.Image) []string {
	lines := make([]string, 0, iconRows/2)
	for y := 0; y < iconRows; y += 2 {
		line := ""
		for x := 0; x < iconColumns; x++ {
			top := averageColor(img, x, y)
			bottom := averageColor(img, x, y+1)
			line += fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m",
				top.r, top.g, top.b,
				bottom.r, bottom.g, bottom.b,
			)
		}
		lines = append(lines, line)
	}
	return lines
}

type rgb struct {
	r int
	g int
	b int
}

func averageColor(img image.Image, cellX, cellY int) rgb {
	b := img.Bounds()
	x0 := b.Min.X + cellX*b.Dx()/iconColumns
	x1 := b.Min.X + (cellX+1)*b.Dx()/iconColumns
	y0 := b.Min.Y + cellY*b.Dy()/iconRows
	y1 := b.Min.Y + (cellY+1)*b.Dy()/iconRows

	var r, g, bl, count uint64
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			cr, cg, cb, ca := img.At(x, y).RGBA()
			alpha := ca >> 8
			r += blendOnDark(cr>>8, alpha)
			g += blendOnDark(cg>>8, alpha)
			bl += blendOnDark(cb>>8, alpha)
			count++
		}
	}
	if count == 0 {
		return rgb{30, 30, 30}
	}
	return rgb{
		r: int(r / count),
		g: int(g / count),
		b: int(bl / count),
	}
}

func blendOnDark(value, alpha uint32) uint64 {
	const bg = 24
	return uint64((value*alpha + bg*(255-alpha)) / 255)
}
