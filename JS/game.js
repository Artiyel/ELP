"use strict";

const Controller = require("./controller"); // on va chercher controller.js

class Game { 
  constructor() { // on initialise la pioche, la défausse, les joueurs...
    this.players = []; // les joueurs
    this.deck = []; // la pioche
    this.discard = []; // la défausse
    this.controller = new Controller(this); // le controlleur
    this.currentPlayerIndex = 0; // index du joueur qui joue actuellement
    this.log = []; // historique de la partie 
    this.gameOver = false; // indique si la partie est finie
    this.prediction=[] // prédictions IA
  }

  init(playerNames) {
    this.players = playerNames.map(name => ({ // pour chaque joueur dans playerNames on initialise des paramètres
      name, // nom
      cards: [], // ses cartes
      status: "active", // son statut
      score: 0, // son score
      secondChance: false // si il a une second chance
    }));

    this.initDeck(); // on initialise la pioche
    this.shuffleDeck(); // on la mélange
  }

  initDeck() {
    this.deck = [];

    for (let i = 0; i <= 12; i++) { // on crée les cartes nombres
      let count = 12 - i + 1;
      for (let j = 0; j < count; j++) this.deck.push(i);
    }

    const specials = ["freeze", "flip3", "secondChance", "+", "x"]; // on crée les cartes spéciales
    specials.forEach(card => {
      for (let j = 0; j < 2; j++) this.deck.push(card);
    });
  }

  shuffle(array) { // mélanger
    for (let i = array.length - 1; i > 0; i--) { // pour chaque carte d'un array
      const j = Math.floor(Math.random() * (i + 1));
      [array[i], array[j]] = [array[j], array[i]]; // on échange la carte avec une carte d'un indice aléatoire
    }
    return array;
  }

  shuffleDeck() {
    this.deck = this.shuffle(this.deck); // on mélange le deck
  }

  drawCard() { // on pioche une carte
    if (this.deck.length === 0) {
      this.reshuffleDiscard(); // si la pioche est vide on mélange la défausse
    }

    const card = this.deck.pop(); // on enlève cette carte du deck
    return card;
  }

  predict() {
    let cartes_qui_restent = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    let nb_carte_en_deck = this.deck.length
    let liste_proba = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    
    for (var i = 0; i < nb_carte_en_deck; i++){
      for (var j = 0; j < 13; j++){
        if (j === this.deck[i]) { //si une carte dans le deck est une carte a numéro:
          cartes_qui_restent[j]+=1  //on compte le nombre de cartes a numéro présentes encore dans le deck
        }
      }
    }
    for (var j = 0; j < 13; j++){
        liste_proba[j]= cartes_qui_restent[j]/nb_carte_en_deck
    }
    let max_proba = list[0];
    let max_index=0
    for (var i = 1; i < list.length; i++) {
      if (list[i] > max_proba) {  
        max_proba = list[i];
        max_index=i
      }
    }
    return max_index // Renvoie l'index de la carte la plus probable de sortir (entre 0 et 12)
  }

  is_beneficial(player) {
    const main_joueur = new Set(player.cards.filter(c => typeof c === "number"));
    carte_qui_sort_probablement = this.predict()
    if (main_joueur.has(carte_qui_sort_probablement)) { //si la carte qui sort est dans la main du joueur
      console.log("The so-called 'AI' recommends that you 'HIT'")
    }
    else {
      console.log("The so-called 'AI' recommends that you 'STAY'")
    }
  }

  reshuffleDiscard() { 
    this.logAction("Pioche vide = mélange de la défausse");
    this.deck = this.shuffle([...this.discard]); // La pioche devient la défausse mélangée
    this.discard = []; // la défausse est vide
  }

  logAction(action) { // pour afficher des infos dans le terminal et les stocker dans log
    console.log(action);
    this.log.push(action);
  }

  isFlip7(player) {
    const unique = new Set(player.cards.filter(c => typeof c === "number")); // on filtre les numéros et on garde que les cartes uniques avec set
    return unique.size >= 7; // on regarde si les cartes nombres du joueur sont 7 ou plus
  }

  endOfRound() { // fin d'une manche
  const active = this.players.filter(p => p.status === "active"); // on garde que les joueurs actifs
  return active.length === 1 || this.players.some(p => this.isFlip7(p)); // on retourne si un seul joueur est actif ou si un joueur a fait un flip7
}

  endOfRoundCleanup() {
  this.players.forEach(player => { // pour chaque joueur
    // Envoyer toutes ses cartes dans la défausse
    this.discard.push(...player.cards);
    // On reset ses cartes
    player.cards = [];
    // On reset son status pour la prochaine manche et on enlève si il a une second chance
    player.status = "active";
    player.secondChance = false;
  });
}
  

 calculateScores() { // on calcule le score
  this.players.forEach(player => {
    if (player.status === "busted") return; // si le joueur a perdu il a rien d'ajouté

    let sum = 0;
    let mult = 1;
    let bonus = 0;

    player.cards.forEach(c => {
      if (typeof c === "number") sum += c; // on ajoute les cartes nombres
      else if (c === "+") bonus += 5; // on ajoute un éventuel bonus
      else if (c === "x") mult *= 2; // on multiplie éventuellement en cas de bonus
    });

    if (this.isFlip7(player)) bonus += 15; // si le joueur a fait flip7 

    player.score += sum * mult + bonus; // calcul du score
  });

  this.endOfRoundCleanup();
}


  checkForWinner() { // on regarde si un joueur a atteint au moins 200 points 
    const winner = this.players.find(p => p.score >= 200);
    if (winner) {
      console.log(`${winner.name} a atteint 200 points et remporte la partie !`);
      this.gameOver = true; // le jeu est alors fini
      return true;
    }
    return false;
  }
}

module.exports = Game;
