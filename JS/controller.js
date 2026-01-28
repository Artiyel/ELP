// controller.js
"use strict";

class Controller {
  constructor(game) {
    this.game = game;
  }

  hit(player) {
    if (player.status !== "active") return;

    const card = this.game.drawCard();
    this.game.logAction(`${player.name} drew ${card}`);

    // ----- CARTES NORMALES -----
    if (typeof card === "number") {
      this.handleNumberCard(player, card);
      return card;
    }

    // ----- BONUS -----
    if (card === "+" || card === "x") {
      player.cards.push(card);
      this.game.logAction(`${player.name} keeps bonus ${card}`);
      return card;
    }

    // ----- ACTIONS -----
    if (card === "secondChance") {
      this.applySecondChance(player);
      this.game.discard.push(card);
      return card;
    }

    if (card === "freeze") {
      this.applyFreeze(player);
      this.game.discard.push(card);
      return; card
    }

    if (card === "flip3") {
      this.handleFlip3(player);
      this.game.discard.push(card);
      return card;
    }
  }

  // ======================
  // CARTES NORMALES
  // ======================
  handleNumberCard(player, card) {
    if (player.cards.includes(card)) {
      if (player.secondChance) {
        player.secondChance = false;
        player.cards.push(card);
        this.game.logAction(`${player.name} used Second Chance!`);
      } else {
        player.status = "busted";
        this.game.logAction(`${player.name} drew duplicate ${card} and is busted!`);
      }
    } else {
      player.cards.push(card);
      this.game.logAction(
        `${player.name} new cards: [${player.cards.join(", ")}]`
      );
    }

    if (this.game.isFlip7(player)) {
      this.game.logAction(`${player.name} completed FLIP 7!`);
    }
  }

  // ======================
  // ACTIONS SIMPLES
  // ======================
  applySecondChance(player) {
    if (!player.secondChance) {
        player.secondChance = true;
        this.game.logAction(`${player.name} gains Second Chance`);
        return;
    }

    this.game.logAction(
        `${player.name} already has Second Chance and must give it`
    );

    const target = this.chooseTargetManually(player, "SecondChance", "secondChance");
    if (!target) return;

    target.secondChance = true;
    this.game.logAction(
        `${player.name} gives Second Chance to ${target.name}`
    );
    }



  applyFreeze(player) {
    const target = this.chooseTargetManually(player, "Freeze", "freeze");
    if (!target) return;

    target.status = "busted";
    this.game.logAction(`${player.name} freezes ${target.name}`);
}

  // ======================
  // FLIP 3
  // ======================
  handleFlip3(player) {
    const target = this.chooseTargetManually(player, "Flip3", "flip3");
    if (!target) return; 

    this.game.logAction(`${target.name} activates Flip3!`);

    let pendingActions = [];

    for (let i = 0; i < 3; i++) {
      const card = this.game.drawCard();
      this.game.logAction(`${target.name} drew ${card}`);

      if (typeof card === "number") {
        if (target.cards.includes(card)) {
          if (target.secondChance) {
            target.secondChance = false;
            target.cards.push(card);
            this.game.logAction(`${target.name} used Second Chance!`);
          } else {
            target.status = "busted";
            this.game.logAction(`${target.name} busted during Flip3!`);
            break;
          }
        } else {
          target.cards.push(card);
        }

        if (this.game.isFlip7(target)) {
          this.game.logAction(`${target.name} completed FLIP 7!`);
          break;
        }

      } else if (card === "+" || card === "x") {
        target.cards.push(card);

      } else {
        // action → stockée
        pendingActions.push(card);
        this.game.discard.push(card);
      }
    }

    // actions données même si busted
    for (const action of pendingActions) {
        if (action === "freeze") {
            const target_2 = this.chooseTargetManually(target, "Freeze", "freeze");
            if (target_2) {
            target_2.status = "busted";
            this.game.logAction(
                `${target.name} gives Freeze to ${target_2.name}`
            );
            }
        }

        if (action === "secondChance") {
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
    const readlineSync = require("readline-sync");

    // On exclut le joueur lui-même
    const targets = this.game.players.filter(
        p => p.status === "active" && p !== fromPlayer
    );

    if (targets.length === 0) {
        // Aucun autre joueur actif
        if (actionName !== "SecondChance") {
            this.game.logAction(
                `${fromPlayer.name} keeps ${actionName} because no other players are active`
            );
        } else {
            this.game.logAction(
                `${fromPlayer.name} discards ${actionName} because no other players are active`
            );
            this.game.discard.push(card);
        }
        return null;
    }

    // Affichage des cibles possibles
    targets.forEach((p, i) => {
        console.log(`[${i}] ${p.name}`);
    });

    let index;
    do {
        index = readlineSync.questionInt("Choose a player to give your card to: ");
    } while (index < 0 || index >= targets.length);

    return targets[index];
}



  stay(player) {
    this.game.logAction(`${player.name} stays.`);
  }
}

module.exports = Controller;

