package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func checkInternet() string {
	host := "8.8.8.8"
	port := "53"
	_, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return " ☒ Интернет-соединение отсутствует."
	} else {
		return " ☑ Интернет-соединение присутствует."
	}
}

func checkFirewall() string {
	cmd, _ := exec.Command("netsh", "advfirewall", "show", "currentprofile").Output()
	if strings.Contains(string(cmd), "Ok.") {
		return " ☑ Брандмауэр обнаружен."
	}
	return " ☒ Брандмауэр НЕ обнаружен."
}

func checkFirewall2() string {
	hasInternet := strings.Contains(checkInternet(), "☑")
	if !hasInternet {
		return " ☒ Ошибка: подключите устройству к Интернету."
	}
	err := exec.Command("netsh", "advfirewall", "set", "allprofiles", "firewallpolicy", "blockinbound,blockoutbound").Run()
	if err != nil {
		log.Fatal(err)
	}
	hasInternet = strings.Contains(checkInternet(), "☑")
	err = exec.Command("netsh", "advfirewall", "set", "allprofiles", "firewallpolicy", "blockinbound,allowoutbound").Run()
	if err != nil {
		log.Fatal(err)
	}
	if !hasInternet {
		return " ☑ Брандмауэр работает нормально."
	}
	return " ☒ Брандмауэр НЕ работает нормально."
}

func checkAntivirus() string {
	cmd, err := exec.Command("powershell", "-NoProfile", "Get-WmiObject", "-Namespace", "\"root\\SecurityCenter2\"", "-Class", "AntiVirusProduct").Output()
	if err != nil {
		log.Fatal(err)
	}
	t1 := "displayName              : "
	t2 := "\ninstanceGuid"
	scmd := string(cmd)
	antiviruses := ""
	for i := strings.Index(scmd, t1); i != -1; i = strings.Index(scmd, t1) {
		a := i + len(t1)
		b := strings.Index(scmd, t2)
		antiviruses += scmd[a:b-1] + ", "
		scmd = scmd[b+len(t2):]
	}
	if antiviruses == "" {
		return " ☒ Антивирусное ПО на устройстве не обнаружено"
	}
	return " ☑ Обнаружено антивирусное ПО: " + antiviruses[:len(antiviruses)-2]
}

func writeNewFile() {
	file, err := os.Create("checkTestVirus.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	_, err = file.WriteString("X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*")
	if err != nil {
		log.Fatal(err)
	}
}

func checkEICAR() string {
	writeNewFile()
	file, err := os.Open("checkTestVirus.txt")
	if err != nil {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove("checkTestVirus.txt")
		if err != nil {
			log.Fatal(err)
		}
		return " ☒ Антивирус НЕ работает."
	}
	return " ☑ Антивирус работает."
}

func writeResult(filename string, args ...string) error {
	if len(args) == 0 {
		return errors.New("there is nothing to write to the file")
	}
	file, err := os.Create(filename)
	if err != nil {
		file, err = os.Open(filename)
		if err != nil {
			return err
		}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	for _, v := range args {
		_, err = file.WriteString(v + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func checkAll() string {
	return "\n Результаты комплесной проверки:\n" +
		checkInternet() + "\n" +
		checkFirewall() + "\n" +
		checkFirewall2() + "\n" +
		checkAntivirus() + "\n" +
		checkEICAR() + "\n"
}

func checkAndSave() {
	s := checkAll()
	fmt.Print(" Введите название/полный путь к файлу.\n" +
		" Пример: \n" +
		"	result.txt\n" +
		" 	C:/Users/User1/Desktop/result.txt\n" +
		" Ввод: ")
	filename := ""
	_, err := fmt.Scanf("%s\n", &filename)
	if err != nil && strings.ContainsRune(filename, '<') || strings.ContainsRune(filename, '<') ||
		strings.ContainsRune(filename, '>') || strings.ContainsRune(filename, ':') ||
		strings.ContainsRune(filename, '"') || strings.ContainsRune(filename, '/') ||
		strings.ContainsRune(filename, '\\') || strings.ContainsRune(filename, '|') ||
		strings.ContainsRune(filename, '?') || strings.ContainsRune(filename, '*') {
		fmt.Println(" Ошибка ввода, введите корректный номер")
	}
	err = writeResult(filename, "Проверка от ", time.Now().String(), s)
	if err != nil {
		fmt.Println(" Ошибка при записи результатов")
	} else {
		fmt.Println(" Результаты проверки сохранены в:", filename)
	}
}

func menu() {
	fmt.Print("\n Введите номер действия: " +
		"\n	0 - выход" +
		"\n	1 - проверка наличия соединения с Интернетом" +
		"\n	2 - проверка наличия установленного межсетевого экрана" +
		"\n	3 - проверка работоспособности межсетевого экрана" +
		"\n	4 - проверка наличия установленного антивируса" +
		"\n	5 - проверка работоспособности антивирусного ПО" +
		"\n	6 - комплексная проверка безопасности ПК" +
		"\n	7 - комплексная проверка с записью в файл" +
		"\n Ввод >> ")
}

func main() {
	err := exec.Command("chcp", "65001").Run()
	if err != nil {
		log.Fatal(err)
	}
	for i := -1; i != 0; {
		menu()
		_, err = fmt.Scanf("%d\n", &i)
		if err != nil {
			fmt.Println("Ошибка ввода, введите корректный номер")
			continue
		}
		switch i {
		case 0:
			return
		case 1:
			fmt.Print("\n", checkInternet())
		case 2:
			fmt.Print("\n", checkFirewall())
		case 3:
			fmt.Print("\n", checkFirewall2())
		case 4:
			fmt.Print("\n", checkAntivirus())
		case 5:
			fmt.Print("\n", checkEICAR())
		case 6:
			fmt.Print(checkAll())
		case 7:
			checkAndSave()
		default:
			fmt.Println("Ошибка ввода, введите корректный номер")
		}
	}
}
