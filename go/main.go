package main

func main() {
	// beforeFile, err := os.Open("testdata/before.txt")
	// if err != nil {
	// 	log.Fatalf("opening file failed %s", err)
	// }
	// defer beforeFile.Close()

	// afterFile, err := os.Open("testdata/after.txt")
	// if err != nil {
	// 	log.Fatalf("opening file failed %s", err)
	// }
	// defer afterFile.Close()

	// // copy before.txt to temp file
	// newFileName := "testdata/temp.txt"
	// newFile, err := os.Create(newFileName)
	// if err != nil {
	// 	log.Fatalf("opening file failed %s", err)
	// }
	// defer newFile.Close()
	// n, err := io.Copy(newFile, beforeFile)
	// if err != nil {
	// 	log.Fatalf("copying file failed %s", err)
	// }
	// fmt.Printf("%d chars copied \n", n)
	// if newFile.Close() != nil {
	// 	log.Fatalf("closing file failed %s", err)
	// }

	// // modify new file
	// writeFile, err := vscode.NewFileHandler(newFileName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer writeFile.Close()
	// // writeFile.Insert(vscode.Position{Line: 33, Character: 17}, "props.todos.reverse()")

	// // os.Remove(newFileName)
}
