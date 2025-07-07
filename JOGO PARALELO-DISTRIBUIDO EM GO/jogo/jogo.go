package main

import (
	"bufio"
	"os"
)

type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool 
}

type Jogo struct {
	Mapa           [][]Elemento 
	PosX, PosY     int          
	UltimoVisitado Elemento     
	StatusMsg      string       
	OutrosJogadores []Jogador 
}


var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
	Alavanca   = Elemento{'A', CorVerde, CorPadrao, true}
)

func jogoNovo() Jogo {
	return Jogo{UltimoVisitado: Vazio}
}


func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			 case Alavanca.simbolo:
        		e = Alavanca
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y

				//A posição onde estava o personagem será substituída pelo elemento Vazio.
				//Assim, o personagem não fica "fixo" no mapa, 
				//mas sua posição é controlada separadamente pelas variáveis PosX e PosY.
				e = Vazio
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {

	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	
    for _, j := range jogo.OutrosJogadores {
        if j.X == x && j.Y == y && j.Nome != nomeJogador {
            return false
        }
    }

	
	return true
}

func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	
	elemento := jogo.Mapa[y][x] 

	jogo.Mapa[y][x] = jogo.UltimoVisitado   
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] 
	jogo.Mapa[ny][nx] = elemento            
}
