"use strict"; // active le mode strict : JS =  plus sûr et moins permissif

const readlineSync = require("readline-sync"); // pour lire les entrées clavier
const Game = require("./game"); //on va chercher game.js

async function main() {
  // fonction qui fait tourner le jeu 

  console.log("Welcome to Flip7 Solo!");

  let logtxt ='';
  logtxt = "date, player, card, action, score \n"
  let logline = ''

  const numPlayers = parseInt(readlineSync.question("Number of players: "));
  const players = [];
  
  for (let i = 0; i < numPlayers; i++) {
    const name = readlineSync.question(`Enter name for player ${i + 1}: `); // pour chaque joueur on récupère son nom
    players.push(name); // on l'ajoute à notre liste de joueurs
  }

  const game = new Game();
  game.init(players); // on initialise les joueurs, leurs cartes...

  while (true) {
    for (let i = 0; i < game.players.length; i++) {
      logline = ', ';
      logline += i; // numéro du joueur
      logline += ", "
      const player = game.players[i]; // on prend chaque joueur dans l'ordre
      if (player.status !== "active") continue;

      console.log(`\n ${player.name}'s turn (score: ${player.score})`);
      const action = readlineSync.keyIn(`Press [Space] to HIT, [f] to STAY, [q] to QUIT : `, {limit: " fq"}); // on fait jouer le joueur 
      if (action === " ") { // si le joueur HIT
        logline += game.controller.hit(player); // on récupère la carte tirée
        logline += ', hit, ';
      } else if (action === "f") { // si le joueur STAY
        logline += "none, stay, ";
        game.controller.stay(player); 
      }
      else if (action === "q") { // Si on décide d'arrêter 
        logline += "game closed";
        logtxt += Date().concat(logline,"\n")
        console.log("Bye !");
        console.log(logtxt);
        process.exit(0); // 0 = succès
        }
      logline += player.score; // le score du joueur
      logtxt += Date().concat(logline,"\n")
      console.log(logline);
      if (game.endOfRound()) break;
    }

    if (game.endOfRound()) { // fin d'une manche
      console.log("\n End of round!");
      game.calculateScores(); // on calcule les scores de chaque joueur
      game.players.forEach(p => console.log(`${p.name}: ${p.score} points`)); // on affiche leurs scores
      if (game.checkForWinner()) { // si un joueur a atteint 200 points 
        console.log("Fin de la partie !"); 
        console.log(logtxt);
        process.exit(0); // quitte Node.js proprement
    }

      // Réinitialiser pour nouvelle manche
      game.players.forEach(p => {
        p.status = "active";
        p.cards = [];
        p.secondChance = false;
      });
    }
  }
}

main();
