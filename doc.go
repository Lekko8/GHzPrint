package main

import (
	"fmt"
	"os"
	"strings"
)

// Функция генерации .doc файла (Word отлично понимает HTML-разметку с расширением .doc)
func generateDocFile(files []FileWData, filename, date string) error {
	content := fmt.Sprintf(`
		<html>
		<head><meta charset="utf-8">
		<style>
        	body { font-family: Arial, sans-serif; font-size: 10pt; line-height: 1.0;}
        	h1 { text-align: center; font-size: 16pt; }
        	table { border-collapse: collapse; width: 100%%; page-break-before: always; page-break-inside: avoid; text-align: center;
            break-inside: avoid;}
        	th, td { border: 1px solid black; padding: 4px; }
    	</style>
		</head>
		<body>
			<h1 style="text-align:center;">Отчёт по результатам измерения на ... </h1>
			<p><b>Дата составления:</b> %s</p>
			<p><b>Исполнитель: _________________</b> </p>
			<p><b>Данные измерений:</b> %s</p>
			<br><br>
			<p></p><br><br><br><br>
			<p style="text-align:right;">Подпись: _______________</p>
		</body>
		</html>
	`, date, dataTostring(dataCreate(files)))

	return os.WriteFile(filename, []byte(content), 0644)
}

// собирает массив данных из файлов в 1 строку для вывода в .doc
func dataTostring(datamass []FileWData) string {
	var str strings.Builder
	var number = 0
	str.WriteString("<table><tr>")
	for _, file := range datamass {
		if number%3 != 0 || number == 0 {
			str.WriteString("<td>")
			number++
		} else if number > 0 {
			str.WriteString("</tr></table>")
			str.WriteString("<div style=\"page-break-before: always; clear: both;\"><br></div><table><tr><td>")
			number++
		}

		str.WriteString(file.sample)
		str.WriteString(" ")
		str.WriteString(file.sampletime.Format("02.01.2006 15:04"))
		str.WriteString("<table><tr><th>avgvalue</th><th>maxvalue</th></tr>")
		for _, data := range file.data {
			str.WriteString("<tr><td>")
			str.WriteString(data.avgvalue)
			str.WriteString("</td><td>")
			str.WriteString(data.maxvalue)
			str.WriteString("</td></tr>")
		}
		str.WriteString("</td>")

		str.WriteString("</table></td>")
	}
	str.WriteString("</tr></table>")
	return str.String()
}
