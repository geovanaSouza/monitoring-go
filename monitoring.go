package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 2
const delay = 5

func main() {

	exibeNomes()

	nome, idade := devolveNomeEIdade()
	fmt.Println(nome, "E tenho", idade, "anos")

	exibeIntroducao(nome, idade)

	for {
		exibeMenu()

		comando := leComando()

		if comando == 1 {
			fmt.Println("Monitorando...")
		} else if comando == 2 {
			fmt.Println("Exibindo Logs...")
		} else if comando == 0 {
			fmt.Println("Saindo do programa")
		} else {
			fmt.Println("Não conheço este comando")
		}

		switch comando {
		case 1:
			iniciarMonitoramento(false)
			break
		case 2:
			imprimeLogsComIOUtil()
		case 3:
			imprimeLogsComOSOpen()
		case 4:
			iniciarMonitoramento(true)
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}

}

func devolveNomeEIdade() (string, int) {
	nome := "Geovana"
	var idade int

	fmt.Println("Qual a sua idade?")
	_, err := fmt.Scanf("%d", &idade)

	trataSaida(err, false)

	return nome, idade
}

func exibeIntroducao(nome string, idade int) {
	var versao float32 = 1.1
	var currentAno = 2021

	fmt.Println("Este programa está na versão", versao)

	fmt.Println("Olá, sr(a).", nome)

	anoNascimento := currentAno - idade
	fmt.Println("Você nasceu em", anoNascimento)

	fmt.Println("O tipo da variável nome é", reflect.TypeOf((nome)))
	fmt.Println("O tipo da variável versão é", reflect.TypeOf(versao))

	var preco float32
	fmt.Println("Digite o preço:")
	fmt.Scanf("%f", &preco)
	fmt.Println("O preço é", preco)
}

func exibeMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs com ioutil.ReadFile")
	fmt.Println("3- Exibir Logs com os.Open")
	fmt.Println("4- Iniciar Monitoramento com debug")
	fmt.Println("0- Sair do Programa")
}

func leComando() int {
	var comandoLido int

	_, err := fmt.Scan(&comandoLido)

	trataSaida(err, true)

	fmt.Println("O endereço da minha variável comando é", &comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	fmt.Println("")

	return comandoLido
}

func iniciarMonitoramento(debug bool) {
	fmt.Println("Monitorando...")

	if debug {
		sites2 := []string{"http://random-status-code.herokuapp.com", "https://www.alura.com.br"}

		fmt.Println("Tipo da variável sites:", reflect.TypeOf(sites2))
		fmt.Println("Valores do slice hard-coded:")
		retornaInfoAboutSlice(sites2)
	}

	sites := leSitesDoArquivo(debug)

	if debug {
		fmt.Println("Sites à serem monitorados:")
		fmt.Println(sites)
		for i := 0; i < len(sites); i++ {
			fmt.Println(sites[i])
		}
	}
	for i := 0; i < monitoramentos; i++ {
		for i, site := range sites {
			fmt.Println("Testando site", i, ":", site)
			testaSite(site, debug)
			if debug {
				fmt.Println("")
			}
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}

	fmt.Println("")
}

func exibeNomes() {
	nomes := [4]string{"Douglas", "Daniel", "Bernardo"}
	fmt.Println("Tipo da variável nome:", reflect.TypeOf(nomes))
	retornaInfoAboutArray(nomes)
	nomes[3] = "Aparecida"
	retornaInfoAboutArray(nomes)
}

func retornaInfoAboutArray(nomes [4]string) {
	fmt.Println(nomes)
	fmt.Println("O meu array tem", len(nomes), "itens")
	fmt.Println("O meu array tem capacidade para", cap(nomes))
}

func retornaInfoAboutSlice(values []string) {
	fmt.Println(values)
	fmt.Println("O meu slice tem", len(values), "itens")
	fmt.Println("O meu slice tem capacidade para", cap(values))
}

func trataSaida(erro error, aborta bool) (falhou bool) {
	if erro != nil {
		fmt.Println("Ocorreu um erro:", erro)
		if aborta {
			os.Exit(-1)
		}
		fmt.Println()
		return true
	}
	return false
}

func testaSite(site string, debug bool) {
	resp, err := http.Get(site)
	trataSaida(err, false)

	if debug {
		fmt.Println("Resposta", site, ":", resp)
	}

	if resp != nil {
		if resp.StatusCode == 200 {
			fmt.Println("Site:", site, "foi carregado com sucesso!")
			registraLog(site, true)
		} else {
			fmt.Println("Site", site, "está com problemas. Status Code", resp.StatusCode)
			registraLog(site, false)
		}
	} else {
		fmt.Println("Erro ao checar site")
	}
}

func leSitesDoArquivo(debug bool) []string {
	var sites []string

	arquivo, err := os.Open("sites.txt")

	if debug {
		fmt.Println("Conteúdo do arquivo com os.Open:", arquivo)
	}

	trataSaida(err, false)

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)
		if debug {
			retornaInfoAboutSlice(sites)
		}

		if trataSaida(err, false) {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Erro desconhecido")
				os.Exit(-1)
			}
		}

	}
	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	trataSaida(err, true)

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func imprimeLogsComOSOpen() {
	fmt.Println("Exibindo Logs...")
	var conteudo []string
	arquivo, err := os.Open("log.txt")
	trataSaida(err, true)
	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		if trataSaida(err, false) {
			if err == io.EOF {
				fmt.Println("Fim do arquivo")
				break
			} else {
				fmt.Println("Erro desconhecido!")
				os.Exit(-1)
			}
		}
		conteudo = append(conteudo, linha)
	}
	fmt.Println(conteudo)
	arquivo.Close()
}

func imprimeLogsComIOUtil() {
	arquivo, err := ioutil.ReadFile("log.txt")
	trataSaida(err, false)
	fmt.Println("Conteúdo do arquivo com iotuil.ReadFile:", arquivo)
	fmt.Println(string(arquivo))
}
