package main

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-xmlfmt/xmlfmt"
)

func main() {
	// files, _ := os.ReadDir("./files")

	// today := time.Now()

	// formattedToday := today.Format("020106")

	// for _, file := range files {
	// 	fmt.Println("BSPD.DB.YF.T520B3D.D1R" + file.Name()[26:31] + ".DX" + formattedToday + ".L00300.CPENV")
	// }

	streaming()
}
func streaming() {

	files, _ := os.ReadDir("./files")

	for _, file := range files {

		fmt.Println("Arquivo encontrado: " + file.Name())

		today := time.Now()

		formattedToday := today.Format("020106")

		openedFile, _ := os.Open("./files/" + file.Name())

		xmlDecoder := xml.NewDecoder(openedFile)

		createdFile, _ := os.Create("BSPD.DB.YF.T520B3D.D1R" + file.Name() + ".DX" + formattedToday + ".L00300.CPENV")

		xmlEncoder := xml.NewEncoder(createdFile)

		xmlEncoder.Indent("", "  ")

		for {
			tokenXml, err := xmlDecoder.RawToken()

			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				log.Println(err)
				break
			}

			switch t := tokenXml.(type) {
			case xml.ProcInst:
				continue
			case xml.StartElement:
				if t.Name.Local == "ARRCDOC" {
					continue
				}
			}

			xmlEncoder.EncodeToken(tokenXml)

		}

		if err := xmlEncoder.Close(); err != nil {
			log.Fatal(err)
		}
	}

}

func batch() {

	file, _ := os.Open("small_file.xml")

	content, _ := io.ReadAll(file)

	tempFile, _ := os.Create("temp.xml")

	tempFile.Write([]byte(xmlfmt.FormatXML(string(content), "", "  ")))

	file.Close()
	tempFile.Close()

	removeTags()
}

func removeTags() {
	file, _ := os.Open("temp.xml")

	defer os.Remove("temp.xml")

	newFile, _ := os.Create("batch.xml")

	defer newFile.Close()

	scanner := bufio.NewScanner(file)

	defer file.Close()

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(strings.ToLower(line), "xml") && !strings.Contains(strings.ToLower(line), "<arrcdoc>") && !strings.Contains(strings.ToLower(line), "</arrcdoc>") {
			fmt.Fprintf(newFile, "%s\n", line)
		}
	}
}
