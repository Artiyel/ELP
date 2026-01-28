# Rewrite des Règles
## Étapes d'une partie :
1. Au tout début on indique le nombre de joueur, avec un nom associé a chaque
1. le jeu se déroule en plusieur manches, et chaque
1. premier tour :
chaque joueur prend une carte, si cette carte est une action elle est executée instantanément
1. Pour chaque tour, pour chaque joueur : choisir de "hit" or "stay", 
    * si "hit" donner une carte au joueur, si elle est unique, le joueur rejoue au tour suivant, sinon, il perd instantanément
    * si "stay" le joueur passe son tour sans prendre de carte
1. la manche se finit si un joueur a 7 cartes distinctes dans son jeu (flip7) ou si tout les autres joueurs sonts éliminés
2. les joueurs calculent leurs points a la fin de chaque manches (cf Score)
3. le jeu se continue avec d'autres manches tant que personne n'a 200 points

## Cartes
* Cartes normales : 12 douzes, 11 onzes, 10 dix, ...., 1 un et 1 zéro
* Cartes spéciales : les cartes spéciales sont dans le deck avec les autres, elle peuvent etre appliqué a n'importe quel joueur possédant des cartes, y compris sois même.
    * Cartes action :
        * Frezze : instant élimination de la manche.
        * flip 3 : pioche 3 cartes, si tu pioche des actions elles sont éxécutée (données) a la fin, si tu fini un flip 7 tu interrompt le flip 3, si tu perd et que tu a pioché une action, tu la donne quand même.
        * Second chance : joker, permet de ne pas perdre une seule fois, si un deuxième est pioché, il doit etre donné.
    * Cartes bonus :
        * Cartes "+" : score additionné a celui des cartes a la fin de la manche
        * Carte "x" : score des cartes multiplié

## Score
Le Score total est la somme des Cartes, multiplié par un éventuel multiplicateur, additionné aux éventuels points bonus, auquel s'ajoutent 15 en cas de victoire par flip7  
Points Joueur += (sommes des cartes) * (multiplicateurs) + (points bonus) + (flip7)

# TO DO
- [ ] rajouter le fichier de log (Matin)
- [ ] rajouter "IA" de comptage de carte 
  - [ ] --> renvoie la carte la plus probable de sortir 
    - [x] regarder toute les cartes dans le jeu
    - [x] compter chaque cartes encore présentes
    - [x] calculer les chances
    - [ ] renvoyer la chance la plus élevée
  - [ ] --> regarde si ça bénéficie le joueur  (si il bust ou non)
    - [ ] acceder a la main du joueur
    - [ ] comparer carte prédite avec cartes de la main
    - [ ] afficher dans la console l'avis de l'ia
- [ ] Stats ? (pourcentage de type de carte tirée, nb de flip7 sur la partie, personne qui a le plus été éliminé, durée moyenne d'une manche/tour )
