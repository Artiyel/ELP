// controller.js
"use strict";

class Controller {
  constructor(game) {
    this.game = game; // on récupère le jeu en coursw
  }

  hit(player) {
    if (player.status !== "active") return;

    const card = this.game.drawCard(); // le joueur pioche une carte s'il n'est pas actif
    this.game.logAction(`${player.name} drew ${card}`);


    if (typeof card === "number") {
      this.handleNumberCard(player, card);
      return card;
    }

    if (card === "+" || card === "x") {
      player.cards.push(card);
      this.game.logAction(`${player.name} keeps bonus ${card}`);
      return card;
    }


    if (card === "secondChance") {
      this.applySecondChance(player);
      this.game.discard.push(card);
      return card;
    }

    if (card === "freeze") {
      this.applyFreeze(player); // si c'est une freeze
      this.game.discard.push(card); // la carte est mise dans la défausse
      return card;
    }

    if (card === "flip3") {
      this.handleFlip3(player); // si c'est une fli3
      this.game.discard.push(card); // la carte est mise dans la défausse
      return card;
    }
  }

// Pour les cartes normales
  handleNumberCard(player, card) {
    if (player.cards.includes(card)) { // si le joueur a déjà la carte
      if (player.secondChance) { // si le joueur a une carte second chance
        player.secondChance = false; // il ne l'a plus
        player.cards.push(card); // il garde sa carte nombre
        this.game.logAction(`${player.name} used Second Chance!`);
      } else {
        player.status = "busted"; // sinon il perd
        this.game.logAction(`${player.name} drew duplicate ${card} and is busted!`);
      }
    } else {
      player.cards.push(card); // sinon il garde sa carte nombre
      this.game.logAction(
        `${player.name} new cards: [${player.cards.join(", ")}]`
      );
    }

    if (this.game.isFlip7(player)) {
      this.game.logAction(`${player.name} completed FLIP 7!`); // si il a un flip7
    }
  }

  applySecondChance(player) {
    if (!player.secondChance) { // si le joueur n'a pas de secondChance
        player.secondChance = true; // il la garde
        this.game.logAction(`${player.name} gains Second Chance`);
        return;
    }

    this.game.logAction(
        `${player.name} already has Second Chance and must give it` // si il a déjà une second chance
    );

    const target = this.chooseTargetManually(player, "SecondChance", "secondChance"); // il doit donner la carte à quelqu'un d'autre
    if (!target) return; 
    // si y'a un joueur choisi :
    target.secondChance = true;
    this.game.logAction(
        `${player.name} gives Second Chance to ${target.name}`
    );
    }



  applyFreeze(player) {
    const target = this.chooseTargetManually(player, "Freeze", "freeze"); // le joueur peut choisir de donner sa carte freeze
    if (!target) return;

    target.status = "busted"; // le joueur perd
    this.game.logAction(`${player.name} freezes ${target.name}`);
}

  handleFlip3(player) {
    const target = this.chooseTargetManually(player, "Flip3", "flip3"); // le joueur peut choisir de donner sa carte flip3
    if (!target) return; 

    this.game.logAction(`${target.name} activates Flip3!`); // si un joueur est choisi

    let pendingActions = []; // actions piochées à accomplir après avoir pioché

    for (let i = 0; i < 3; i++) { // le joueur pioche 3 fois
      const card = this.game.drawCard();
      this.game.logAction(`${target.name} drew ${card}`);

      if (typeof card === "number") { // si la carte est un nombre
        this.handleNumberCard(player, card)

      } else if (card === "+" || card === "x") { // si le joueur a une carte action
        target.cards.push(card);

      } else {
        // action = stockée
        pendingActions.push(card);
        this.game.discard.push(card);
      }
    }

    // on parcourt les actions à réaliser
    for (const action of pendingActions) {
        if (action === "freeze") { // si c'est un freeze 
            const target_2 = this.chooseTargetManually(target, "Freeze", "freeze");
            if (target_2) {
            target_2.status = "busted";
            this.game.logAction(
                `${target.name} gives Freeze to ${target_2.name}`
            );
            }
        }

        if (action === "secondChance") { // si c'est un second chance
            const target_2 = this.chooseTargetManually(target, "SecondChance", "secondChance");
            if (target_2) {
            target_2.secondChance = true;
            this.game.logAction(
                `${target.name} gives Second Chance to ${target_2.name}`
            );
            }
        }
}

  }


  chooseTargetManually(fromPlayer, actionName, card) {
    const readlineSync = require("readline-sync"); // pour lire les entrées

    if (actionName !== "SecondChance"){ // si la carte n'est pas une second chance
       let keep = readlineSync.question("Do you want to keep the card ? YES : [y], NO : [n] "); 
       if (keep === "y"){ // si le joueur veut jouer la carte
        return fromPlayer
      }
    }
    
    // On exclut le joueur lui-même
    const targets = this.game.players.filter(
        p => p.status === "active" && p !== fromPlayer // on filtre les autres joueurs, on garde les actifs et on enlève le joueur
    );

    if (targets.length === 0) {  
        // Aucun autre joueur actif
        if (actionName !== "SecondChance") { // si c'est pas une second chance
            this.game.logAction(
                `${fromPlayer.name} keeps ${actionName} because no other players are active`
            );
        } else { // si c'est une second chance
            this.game.logAction(
                `${fromPlayer.name} discards ${actionName} because no other players are active`
            );
            this.game.discard.push(card); // on met la carte dans la défausse
        }
        return null;
    }

    // Affichage des cibles possibles
    targets.forEach((p, i) => {
        console.log(`[${i}] ${p.name}`);
    });

    let index; // contient le numéro choisi par le joueur (numéro du joueur cible dans la liste des joueurs actifs)
    do {
        index = readlineSync.questionInt("Choose a player to give your card to: ");
    } while (index < 0 || index >= targets.length); // Tant que l'index saisi n'est pas valide, on répète le prompt

    return targets[index]; // Retourne le joueur correspondant à l'index choisi
}



  stay(player) {
    this.game.logAction(`${player.name} stays.`); // si le joueur passe
  }
}

module.exports = Controller;

