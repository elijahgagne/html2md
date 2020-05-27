package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: cat index.html | html2md")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	html := string(output)
	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())
	converter.Use(plugin.ConfluenceCodeBlock())
	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
}

// <p>Here is a three line code block</p><ac:structured-macro ac:name="code" ac:schema-version="1" ac:macro-id="17182db8-062e-4400-a53f-04cd1adda8ac"><ac:parameter ac:name="language">bash</ac:parameter><ac:plain-text-body><![CDATA[echo 'hi'
// cat file.txt | grep begin
// mv abc xyz]]></ac:plain-text-body></ac:structured-macro><p />

// <p>Here is a three line code block</p>
// <ac:structured-macro ac:name="code" ac:schema-version="1" ac:macro-id="17182db8-062e-4400-a53f-04cd1adda8ac">
// 	<ac:parameter ac:name="language">bash</ac:parameter>
// 	<ac:plain-text-body>
// 	<![CDATA[echo 'hi'
// cat file.txt | grep begin
// mv abc xyz]]>
// </ac:plain-text-body>
// </ac:structured-macro>
// <p />

// Find:
// <ac:structured-macro ac:name="code" ac:schema-version="1" ac:macro-id="17182db8-062e-4400-a53f-04cd1adda8ac"><ac:parameter ac:name="language">bash</ac:parameter><ac:plain-text-body><![CDATA[
// Replace: ```bash

// Find:
// ]]></ac:plain-text-body></ac:structured-macro>
// Replace: ```

// html := strings.ReplaceAll(string(output), "]]></ac:plain-text-body></ac:structured-macro>", "\n```\n")
// 	html = strings.ReplaceAll(html, "<ac:structured-macro ac:name=\"code\" ac:schema-version=\"1\" ac:macro-id=\"17182db8-062e-4400-a53f-04cd1adda8ac\"><ac:parameter ac:name=\"language\">bash</ac:parameter><ac:plain-text-body><![CDATA[", "\n```\n")
