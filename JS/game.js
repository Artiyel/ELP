"use strict";

const Controller = require("./controller"); // on va chercher controller.js

class Game {
  constructor() {
    this.players = [];
    this.deck = [];
    this.discard = [];
    this.controller = new Controller(this);
    this.currentPlayerIndex = 0;
    this.log = [];
    this.gameOver = false;
    this.prediction=[]
  }

  init(playerNames) {
    this.players = playerNames.map(name => ({
      name,
      cards: [],
      status: "active",
      score: 0,
      secondChance: false
    }));

    this.initDeck();
    this.shuffleDeck();
  }

  initDeck() {
    this.deck = [];

    for (let i = 0; i <= 12; i++) {
      let count = 12 - i + 1;
      for (let j = 0; j < count; j++) this.deck.push(i);
    }

    const specials = ["freeze", "flip3", "secondChance", "+", "x"];
    specials.forEach(card => {
      for (let j = 0; j < 2; j++) this.deck.push(card);
    });
  }

  shuffle(array) {
    for (let i = array.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
  }

  shuffleDeck() {
    this.deck = this.shuffle(this.deck);
  }

  drawCard() {
    if (this.deck.length === 0) {
      this.reshuffleDiscard();
    }

    const card = this.deck.pop();
    return card;
  }

  predict() {
    let cartes_qui_restent = [1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
    let nb_carte_en_jeu = 94
    let liste_proba=[0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    for (var i = 0; i < this.deck.length; i++){
      for (var j = 0; j < 13; j++){
        if (j === this.deck[i]) {
          cartes_qui_restent[j]-=1
        }
      }
      nb_carte_en_jeu -= 1
    }
    for (var j = 0; j < 13; j++){
        liste_proba[j]= cartes_qui_restent[j]/nb_carte_en_jeu
    }
    let max_proba = list[0];
    let max_index=0
    for (var i = 1; i < list.length; i++) {
      if (list[i] > max_proba) {
        max_proba = list[i];
        max_index=i
      }
    }
    return max_index
  }

  reshuffleDiscard() {
    this.logAction("Pioche vide → mélange de la défausse");
    this.deck = this.shuffle([...this.discard]);
    this.discard = [];
  }

  logAction(action) {
    console.log(action);
    this.log.push(action);
  }

  isFlip7(player) {
    const unique = new Set(player.cards.filter(c => typeof c === "number"));
    return unique.size >= 7;
  }

  endOfRound() {
  const active = this.players.filter(p => p.status === "active");
  return active.length === 1 || this.players.some(p => this.isFlip7(p));
}

  endOfRoundCleanup() {
  this.players.forEach(player => {
    // Envoyer toutes les cartes dans la défausse
    this.discard.push(...player.cards);
    // Reset main et cartes action
    player.cards = [];
    // Reset status pour prochaine manche
    player.status = "active";
    player.secondChance = false;
  });
}
  

 calculateScores() {
  this.players.forEach(player => {
    if (player.status === "busted") return;

    let sum = 0;
    let mult = 1;
    let bonus = 0;

    player.cards.forEach(c => {
      if (typeof c === "number") sum += c;
      else if (c === "+") bonus += 5;
      else if (c === "x") mult *= 2;
    });

    if (this.isFlip7(player)) bonus += 15;

    player.score += sum * mult + bonus;
  });

  this.endOfRoundCleanup();
}


  resetPlayersForNextRound() {
    this.players.forEach(p => {
      p.cards = [];
      p.status = "active";
      p.secondChance = false;
    });
  }

  checkForWinner() {
    const winner = this.players.find(p => p.score >= 200);
    if (winner) {
      console.log(`${winner.name} a atteint 200 points et remporte la partie !`);
      this.gameOver = true;
      return true;
    }
    return false;
  }
}

module.exports = Game;
