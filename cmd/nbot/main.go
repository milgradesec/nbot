package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/internal/config"
)

var data = []string{
	"Saca pecho",
	"Si te lo llevas me voy afk",
	"ff 15",
	"Esta over",
	"Percutimos percutimos",
	"Diveamos eh diveamos",
	"Go go go",
	"Estas trolleando",
	"No me llevo ninguna kill",
	"Una mas y me voy afk",
	"Era para baitear el zhonyas",
	"De verdad, de verdad jungler difference",
	"Hemos perdido",
	"Suspiro*",
	"Estas enfermo",
	"Pero puede dejar de pasar bolas el culogordo?",
	"Que tal con tu novia?",
	"Me has encontrado ya novia?",
	"Que pesado",
	"Cabezon",
	"Pideme perdon",
	"PIDE PERDON",
	"Estas bien de la cabeza?",
	"Estas mal de la cabeza",
	"Tengo un papel que demuestra mi 150 de iq",
	"Llevo mas goles que todos vosotros juntos",
	"Xd",
	"XDDDDDDDDDDDDDDDDD",
	"Calvo terrorista",
	"Espera que me pongo las gafas",
	"Va 0/8 y me mata, impresionante",
	"Estoy smurfing",
	"Estoy 1v9",
	"VERDADES",
	"Puedes dejar de hacer el ridiculo?",
	"Pues voy a martillearme la polla",
	"Venid que puedo contra todos",
	"Top daño de la partida",
	"Es worth ya no tiene flash",
	"Play around the caja",
	"Deja de jugar flex como un bot",
	"Problemas", //nolint
	"Eres la persona más egoísta que conozco",
	"Yo me ducho todos los días, bueno hoy no",
	"Caduco",
	"Trilero",
	"Doble moonstone, doble moonstone",
	"Estoy oneshoteando al Nautilus",
	"De verdad me solea una Leona support?",
	"Claro genios",
	"Estamos en la winners queue",
	"Menudo carreo te estoy metiendo",
}

func main() {
	token, found := config.GetToken()
	if !found {
		log.Fatal("error: Bot token not found")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	rand.Seed(time.Now().Unix())
	session.AddHandler(handler)

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	fmt.Println("Bot is running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}

func handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!frases" {
		var msg string
		for _, frase := range data {
			msg += frase
			msg += "\n"
		}
		s.ChannelMessageSend(m.ChannelID, msg) //nolint
		return
	}

	if strings.Contains(m.Content, "nbot") {
		s.ChannelMessageSend(m.ChannelID, data[rand.Intn(len(data))]) //nolint
	}
}
