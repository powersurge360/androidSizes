
/* Copyright (c) 2013 Forest Giant, Inc.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ========================================================== */

package imageconverter

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type ImageConverter struct {
	paths     []string
	Directory string
	Type      string
}

var whitelistedTypes = []string{"jpg", "png", "jpeg"}
var typeWidthModifiers = map[string]map[string]float32{
	"ldpi": {
		"mdpi":  1.333,
		"hdpi":  2.0,
		"xhdpi": 2.666,
	},
	"mdpi": {
		"ldpi":  .75,
		"hdpi":  1.5,
		"xhdpi": 2.0,
	},
	"hdpi": {
		"ldpi":  .5,
		"mdpi":  .75,
		"xhdpi": 1.333,
	},
	"xhdpi": {
		"ldpi": .375,
		"mdpi": .5,
		"hdpi": .666,
	},
}

func (converter *ImageConverter) widthModifier(newType string) float32 {
	return typeWidthModifiers[converter.Type][newType]
}

func (converter *ImageConverter) convertToType(path string, newType string, channel chan int) {
	file, _ := os.Open(path)

	pathParts := strings.Split(path, "/")
	filename := pathParts[len(pathParts)-1]

	picture, format, _ := image.Decode(file)
	rectangle := picture.Bounds()

	modifier := converter.widthModifier(newType)
	newWidth := uint(modifier * float32(rectangle.Dx()))

	if newWidth < 1 {
		newWidth = 1
	}

	newImage := resize.Resize(newWidth, 0, picture, resize.NearestNeighbor)
	typeDirectory := strings.Join([]string{converter.Directory, "..", newType}, "/")

	os.Mkdir(typeDirectory, 0700)

	newFile, _ := os.Create(strings.Join([]string{typeDirectory, filename}, "/"))

	switch format {
	case "jpeg":
		jpeg.Encode(newFile, newImage, nil)
	case "png":
		png.Encode(newFile, newImage)
	}

	file.Close()
	newFile.Close()
	channel <- 1
}

func (converter *ImageConverter) convertImage(path string, channel chan int) {
	subchannel := make(chan int)
	switch converter.Type {
	case "ldpi":
		go converter.convertToType(path, "mdpi", subchannel)
		go converter.convertToType(path, "hdpi", subchannel)
		go converter.convertToType(path, "xhdpi", subchannel)
	case "mdpi":
		go converter.convertToType(path, "ldpi", subchannel)
		go converter.convertToType(path, "hdpi", subchannel)
		go converter.convertToType(path, "xhdpi", subchannel)
	case "hdpi":
		go converter.convertToType(path, "ldpi", subchannel)
		go converter.convertToType(path, "mdpi", subchannel)
		go converter.convertToType(path, "xhdpi", subchannel)
	case "xhdpi":
		go converter.convertToType(path, "ldpi", subchannel)
		go converter.convertToType(path, "mdpi", subchannel)
		go converter.convertToType(path, "hdpi", subchannel)
	}

	for i := 0; i < 3; i++ {
		<-subchannel
	}

	channel <- 1

}

func (converter *ImageConverter) Convert() {
	channel := make(chan int)
	// Get the list of files in the directory
	for _, extension := range whitelistedTypes {
		pathWithExtension := strings.Join([]string{
			"*", extension,
		}, ".")
		completePath := strings.Join([]string{
			converter.Directory, pathWithExtension,
		}, "/")

		files, _ := filepath.Glob(completePath)

		converter.paths = append(converter.paths, files...)
	}

	// Open each file and run converts
	for _, path := range converter.paths {
		go converter.convertImage(path, channel)
	}

	for _ = range converter.paths {
		<-channel
	}
}
