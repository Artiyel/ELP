// index.js
"use strict";

const readlineSync = require("readline-sync");
const Game = require("./game");

async function main() {
  console.log("Welcome to Flip7 Solo!");

  const numPlayers = parseInt(readlineSync.question("Number of players: "));
  const players = [];
  for (let i = 0; i < numPlayers; i++) {
    const name = readlineSync.question(`Enter name for player ${i + 1}: `);
    players.push(name);
  }

  const game = new Game();
  game.init(players);

  while (true) {
    for (let i = 0; i < game.players.length; i++) {
      const player = game.players[i];
      if (player.status !== "active") continue;

      console.log(`\n ${player.name}'s turn (score: ${player.score})`);
      const action = readlineSync.keyIn(`Press [Space] to HIT, [f] to STAY, [q] to QUIT : `, {limit: " fq"});
      if (action === " ") {
        game.controller.hit(player);
      } else if (action === "f") {
        game.controller.stay(player);
      }
      else if (action === "q") {
        console.log("Bye !");
        process.exit(0); // 0 = succès
        }

      if (game.endOfRound()) break;
    }

    if (game.endOfRound()) {
      console.log("\n End of round!");
      game.calculateScores();
      game.players.forEach(p => console.log(`${p.name}: ${p.score} points`));
      if (game.checkForWinner()) {
        console.log("Fin de la partie !");
        process.exit(0); // quitte Node.js proprement
    }

      // Réinitialiser pour nouvelle manche
      game.players.forEach(p => {
        p.status = "active";
        p.cards = [];
        p.secondChance = false;
      });
      game.initDeck();
      game.shuffleDeck();
    }
  }
}

main();
