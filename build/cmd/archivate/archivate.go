package main

import (
    "archive/zip"
    "io"
    "os"
    "path/filepath"
)

func main() {

    
    output := "../../" + os.Args[1]

    files := os.Args[2:]

    if err := ZipFiles(output, files); err != nil {
        panic(err)
    }
}

func ZipFiles(zipname string, files []string) error {

    newZipFile, err := os.Create(zipname)
    if err != nil {
        return err
    }
    defer newZipFile.Close()

    zipWriter := zip.NewWriter(newZipFile)
    defer zipWriter.Close()

    for _, file := range files {
	path, filename  := filepath.Split(file)
	if err := os.Chdir("../../" + path); err != nil {
    		panic(err)
  	}
        if err = AddFileToZip(zipWriter, filename); err != nil {
            return err
        }
	if err := os.Chdir("../../cmd/archivate"); err != nil {
    		panic(err)
  	}	
    }

    return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

    fileToZip, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer fileToZip.Close()


    info, err := fileToZip.Stat()
    if err != nil {
        return err
    }

    header, err := zip.FileInfoHeader(info)
    if err != nil {
        return err
    }

    header.Name = filename
    header.Method = zip.Deflate

    writer, err := zipWriter.CreateHeader(header)
    if err != nil {
        return err
    }
    _, err = io.Copy(writer, fileToZip)
    return err
}
