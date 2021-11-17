//go:generate fyne bundle --output bundle-regular.go --package main font/NotoSans-Regular.ttf
//go:generate fyne bundle --output bundle-bold.go --package main font/NotoSans-Bold.ttf
//go:generate fyne bundle --output bundle-bold-italic.go --package main font/NotoSans-BoldItalic.ttf
//go:generate fyne bundle --output bundle-italic.go --package main font/NotoSans-Italic.ttf
//go:generate fyne bundle --output bundle-cjk.go --package main font/NotoSansCJK.ttc
//go:generate fyne package -os android
package main
