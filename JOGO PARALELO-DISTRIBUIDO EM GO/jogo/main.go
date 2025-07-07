package main

import (
	"fmt"
	"net/rpc"
	"os"
	"time"
)


type Jogador struct {
	Nome string
	X, Y int
}

type Posicao struct {
	Nome string
	X, Y int
	SequenceNumber int 
}

var nomeJogador string
var clienteRPC *rpc.Client
var C = make(chan string)

func main() {

	var err error

	var enderecoServidor string
   
    fmt.Print("Digite o endereço do servidor (ex: 192.168.234.86): ")
	fmt.Scanln(&enderecoServidor)

	
    fmt.Print("Digite seu nome de jogador: ")
    fmt.Scanln(&nomeJogador)

	clienteRPC, err = rpc.Dial("tcp", enderecoServidor + ":1234")
	if err != nil {
		panic(err)
	}

	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}


	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	interfaceDesenharJogo(&jogo)

	go func() {
		for msg := range C {

			jogo.StatusMsg = msg

		}
	}()

	go atualizaEstado(&jogo)

	for {
		
		ev := interfaceLerEventoTeclado()
		if !personagemExecutarAcao(ev, &jogo) {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}

var sequenceNumber int = 0

func enviaPosicao(jogo *Jogo) {
	const tentativas = 10
    seqNum := sequenceNumber + 1 // não incrementa ainda o global
    sucesso := false

    for i := 1; i <= tentativas; i++ {
        fmt.Printf("Tentativa %d com seq %d...\n", i, seqNum)
        
        var resposta string
        pos := Posicao{
            Nome: nomeJogador,
            X:    jogo.PosX,
            Y:    jogo.PosY,
            SequenceNumber: seqNum,
        }

        err := clienteRPC.Call("ServidorJogo.EnviarPosicao", pos, &resposta)
        if err == nil {
            jogo.StatusMsg = fmt.Sprintf("Sucesso na tentativa %d: %s", i, resposta)
            sucesso = true
            break // se quiser parar no primeiro sucesso
        } else {
            jogo.StatusMsg = fmt.Sprintf("Falha na tentativa %d: %v", i, err)
        }

        time.Sleep(500 * time.Millisecond) // pequena pausa entre as tentativas
    }


    if sucesso {
        sequenceNumber++ // só avança se algum envio foi bem-sucedido
    } else {
        fmt.Printf("Todas as %d tentativas falharam para seq %d.\n", tentativas, seqNum)
    }
}



func atualizaEstado(jogo *Jogo) {
	for {
		var estado []Jogador
		err := clienteRPC.Call("ServidorJogo.ObterEstado", "", &estado)
		if err == nil {
			jogo.OutrosJogadores = estado
			
            interfaceDesenharJogo(jogo)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
