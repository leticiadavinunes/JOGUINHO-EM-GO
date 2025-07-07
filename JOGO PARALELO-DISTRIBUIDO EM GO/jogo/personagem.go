package main

import "fmt"

func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1 // Move para cima
	case 'a':
		dx = -1 // Move para a esquerda
	case 's':
		dy = 1 // Move para baixo
	case 'd':
		dx = 1 // Move para a direita
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
		enviaPosicao(jogo)
	}

}

func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
		personagemInteragir(jogo)
	case "mover":
		personagemMover(ev.Tecla, jogo)
	}
	return true 
}

type InteracaoAlavanca struct {
    NomeJogador    string
    NomeAlavanca   string
    SequenceNumber int
}

var alavancaSeqNum int = 0

func personagemInteragir(jogo *Jogo) {
	// Verifica as 4 direções ao redor do personagem
    direcoes := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
    encontrou := false
    var ax, ay int

    for _, d := range direcoes {
        nx, ny := jogo.PosX+d[0], jogo.PosY+d[1]
        if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
            if jogo.Mapa[ny][nx].simbolo == Alavanca.simbolo {
                encontrou = true
                ax, ay = nx, ny
                break
            }
        }
    }

    if encontrou {
        nomeAlavanca := fmt.Sprintf("alavanca_%d_%d", ax, ay)
        alavancaSeqNum++
        interacao := InteracaoAlavanca{
            NomeJogador:    nomeJogador,
            NomeAlavanca:   nomeAlavanca,
            SequenceNumber: alavancaSeqNum,
        }
        var resposta string
        err := clienteRPC.Call("ServidorJogo.AtivarAlavanca", interacao, &resposta)
        if err != nil {
            jogo.StatusMsg = "Erro ao ativar alavanca"
        } else {
            jogo.StatusMsg = resposta
        }
    } else {
        jogo.StatusMsg = "Não há alavanca ao lado!"
    }
}