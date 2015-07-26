/*
 *
 * Please see README.md for more information
 *
 * Usage: blur
 * Description: Blurr's the files in the "orig" directory and outputs to an "auto-dest" directory
 *
 */

package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	//
	// Define the factors we want for the blurrings
	//
	// factors := [...]string{"0.5", "0.25", "1"}
	//
	//
	// map from FactorMultiple to FactorName
	factors := map[string]string{
		"0":  "",
		"3":  "-blur3",
// 		"6":  "-blur6",
// 		"9":  "-blur9",
// 		"12": "-blur12",
// 		"15": "-blur15",
// 		"18": "-blur18",
// 		"21": "-blur21",
// 		"24": "-blur24",
// 		"27": "-blur27",
//		"30": "-blur30",
	}

	//
	// clean auto-dest directory
	//
	err := os.RemoveAll("auto-dest")
	errorCheck(err)

	//
	// foreach factor
	//
	for factor, factorName := range factors {
		//
		// convert the factor string to a float64
		//
		factorFloat, err := strconv.ParseFloat(factor, 64)
		errorCheck(err)

		//
		// Create the factor destination directory
		//
		fmt.Println("Factoring for", factorName, "...blur", factor)
		destPath := fmt.Sprintf("auto-dest")
		err = os.MkdirAll(destPath, 0777)
		errorCheck(err)

		//
		// foreach file in the original directory, resize it!
		//
		files, err := ioutil.ReadDir("orig")
		errorCheck(err)

		for _, file := range files {
			origFilename := file.Name()                                 // eg. surfboard.png
			origFilenameWithDir := fmt.Sprintf("orig/%s", origFilename) // eg. "orig/surfboard.png"

			//
			// Work out what our output filename should be
			//
			origFilenameFilepathBase := filepath.Base(origFilename)
			origFilenameExt := filepath.Ext(origFilenameFilepathBase)                                             // eg. ".png"
			origFilenameBasename := origFilenameFilepathBase[:len(origFilenameFilepathBase)-len(origFilenameExt)] // eg. "monk"

			destFilenameWithDir := fmt.Sprintf("%s/%s%s%s", destPath, origFilenameBasename, factorName, origFilenameExt)

			//
			// Detect filetype and process... 
			//
			if (strings.HasSuffix(origFilename, ".png")) {
				blueImage(origFilenameWithDir, destFilenameWithDir, factorFloat, imaging.PNG)
			} else if (strings.HasSuffix(origFilename, ".jpg")) { 
				blueImage(origFilenameWithDir, destFilenameWithDir, factorFloat, imaging.JPEG)

			}else {
				fmt.Println("Not processing", origFilenameWithDir, "as it's type could not be detected")
			}
		}

	}

	// blueImage("orig/monk.png", "result.png", float64(3.5))

}

/*
	resizeImage: resize the image and save it
*/
func blueImage(inFilename string, outFilename string, factor float64, format imaging.Format) {
	//
	// print a message about what we are converting
	//
	fmt.Println("\t\t", inFilename, "=>", outFilename, "using factor", factor)

	//
	// Open the file and read in the image
	//
	infile, err := os.Open(inFilename)
	errorCheck(err)
	defer infile.Close()

	//
	// load in the actual image data (decode image) and work out the new dimensions
	//
	srcImage, _, err := image.Decode(infile)
	errorCheck(err)

	//
	// blue the image...
	//

	blurredImage := imaging.Blur(srcImage, factor)

	//
	// save the new image
	//
	outfile, err := os.Create(outFilename)
	errorCheck(err)
	defer outfile.Close()

    imaging.Encode(outfile, blurredImage, format)

	// png.Encode(outfile, blurredImage)

}

/*
	errorCheck: Helper function to check for errors and log/panic on fail
*/
func errorCheck(e error) {
	if e != nil {
		log.Fatalf("%s", e)
		panic(e)
	}
}
