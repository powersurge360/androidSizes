/* ==========================================================
 * Copyright (c) 2013 Forest Giant, Inc.
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

package main

import (
	"fmt"
	"flag"
	"github.com/powersurge360/androidSizes/imageconverter"
	"runtime"
)

var directory string
var androidType string

func init() {
	flag.StringVar(&androidType, "type", "", "Format of the existing images.")
	flag.StringVar(&directory, "directory", ".", "Directory to look in for the existing images.")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if androidType != "" {
		image := imageconverter.ImageConverter{Type: androidType, Directory: directory}
		image.Convert()
	} else {
		fmt.Println("Please provide a --type argument")
	}
}
